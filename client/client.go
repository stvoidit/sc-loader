package client

import (
	"context"
	"encoding/json/v2"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"sc-loader/models"
	"sc-loader/utils"
	"time"
)

type Probe struct{}

type ClientAPI struct {
	_c             *http.Client
	cnf            utils.Config
	defaultHeaders http.Header
	writedParts    utils.Set[string]
}

func NewClient(cnf utils.Config) *ClientAPI {
	jar, _ := cookiejar.New(nil)
	transport := &http.Transport{
		Proxy:               http.ProxyFromEnvironment,
		MaxIdleConnsPerHost: 10,
		MaxIdleConns:        10,
		IdleConnTimeout:     time.Millisecond * 200,
	}
	_c := &http.Client{
		Timeout:   10 * time.Second,
		Jar:       jar,
		Transport: transport,
	}
	return &ClientAPI{
		_c:          _c,
		cnf:         cnf,
		writedParts: make(utils.Set[string]),
		defaultHeaders: http.Header{
			"Accept":          []string{"*/*"},
			`Accept-Language`: []string{`en-US`},
			`Accept-Encoding`: []string{`gzip`},
			`Sec-Fetch-Mode`:  []string{"cors"},
			`Sec-Fetch-Site`:  []string{`same-origin`},
			`Sec-GPC`:         []string{"1"},
			"Origin":          []string{fmt.Sprintf(`https://%s`, cnf.Host)},
			`Referer`:         []string{fmt.Sprintf(`https://%s/`, cnf.Host)},
			`User-Agent`:      []string{`Mozilla/5.0 (X11; Linux x86_64; rv:147.0) Gecko/20100101 Firefox/147.0`},
		},
	}
}

func (c *ClientAPI) filterParts(links []string) []string {
	if links == nil {
		return nil
	}
	var newLinks = make([]string, 0, len(links))
	for i := range links {
		if !c.writedParts.Has(links[i]) {
			newLinks = append(newLinks, links[i])
		}
	}
	return newLinks
}

func (c *ClientAPI) makeRequest(ctx context.Context, link string) (req *http.Request, err error) {
	req, err = http.NewRequestWithContext(ctx, http.MethodGet, link, nil)
	if err != nil {
		return req, err
	}
	req.Header = c.defaultHeaders.Clone()
	return req, err
}

func (c *ClientAPI) doGetRequest(req *http.Request) (res *http.Response, err error) {
	res, err = c._c.Do(req)
	if err != nil {
		return nil, err
	}
	slog.Debug("doGetRequest",
		slog.String("status", res.Status),
		slog.String("proto", res.Proto),
		slog.String("url", req.URL.String()),
		// slog.Any("headers", res.Header),
	)
	err = unzipResponse(res)
	return
}

func (c *ClientAPI) Init(ctx context.Context) (cdnDomain string, err error) {
	id, err := c.InitClientConfig(ctx)
	return id.GetDomainCDN(), err
}

func (c *ClientAPI) InitClientConfig(ctx context.Context) (id models.InitialDynamic, err error) {
	var urlConfig = makeConfigURL(c.cnf.Host)
	slog.Debug("Init", slog.String("userRoomUrl", urlConfig))
	req, err := c.makeRequest(ctx, urlConfig)
	if err != nil {
		return id, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.doGetRequest(req)
	if err != nil {
		return id, err
	}
	defer utils.DeferCloseReader(res.Body)
	if xapiversion := res.Header.Get("X-Api-Version"); !utils.IsEmpty(xapiversion) {
		c.defaultHeaders.Set("front-version", xapiversion)
	}
	var idr models.InitialDynamicResponse
	if err := json.UnmarshalRead(res.Body, &idr); err != nil {
		return id, err
	}
	return idr.InitialDynamic, err
}

func (c *ClientAPI) GetRoomID(ctx context.Context, username string) (roomID int, online bool, err error) {
	var userRoomUrl = makeRoomURL(c.cnf.Host, username)
	slog.Debug("GetRoomID", slog.String("userRoomUrl", userRoomUrl))
	req, err := c.makeRequest(ctx, userRoomUrl)
	if err != nil {
		return -1, false, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.doGetRequest(req)
	if err != nil {
		return -1, false, err
	}
	defer utils.DeferCloseReader(res.Body)
	// res.Body = utils.SaveFile("room.json", res.Body)
	var rr models.RoomResponse
	if err := json.UnmarshalRead(res.Body, &rr); err != nil {
		return -1, false, err
	}
	return rr.GetRoomId(), rr.IsOnline(), err
}

func (c *ClientAPI) GetRoomStatus(ctx context.Context, username string) (status string, online bool, err error) {
	var userRoomUrl = makeRoomURL(c.cnf.Host, username)
	slog.Debug("GetRoomID", slog.String("userRoomUrl", userRoomUrl))
	req, err := c.makeRequest(ctx, userRoomUrl)
	if err != nil {
		return "", false, err
	}
	req.Header.Set("Accept", "application/json")
	res, err := c.doGetRequest(req)
	if err != nil {
		return "", false, err
	}
	defer utils.DeferCloseReader(res.Body)
	// res.Body = utils.SaveFile("room.json", res.Body)
	var rr models.RoomResponse
	if err := json.UnmarshalRead(res.Body, &rr); err != nil {
		return "", false, err
	}
	return rr.GetStatus(), rr.IsOnline(), err
}

var ErrPKeyNotFound = errors.New("pkey not found")

func (c *ClientAPI) GetPlaylistVariants(
	ctx context.Context,
	cdnDomain string,
	roomId int,
) (playlist string, psch utils.PSCH, err error) {
	var playlistUrl = makeMasterPlaylistURL(cdnDomain, roomId)
	slog.Debug("GetPlaylistVariants", slog.String("playlistUrl", playlistUrl))
	req, err := c.makeRequest(ctx, playlistUrl)
	if err != nil {
		return playlist, psch, err
	}
	res, err := c.doGetRequest(req)
	if err != nil {
		return playlist, psch, err
	}
	variant, pkey := utils.ParseRaw(res.Body)
	if utils.IsEmpty(pkey) {
		return playlist, psch, ErrPKeyNotFound
	}
	psch, err = c.cnf.CacheKeys.GetKey(pkey)
	if err != nil {
		return playlist, psch, err
	}
	playlist, err = makeMasterPlaylist(variant, pkey)
	return playlist, psch, err
}

func (c *ClientAPI) GetPlaylistVideo(
	ctx context.Context,
	videoList string,
	psch utils.PSCH,
	more url.Values,
) ([]string, url.Values, error) {
	link, err := makeVideoPlaylist(videoList, more)
	if err != nil {
		return nil, nil, err
	}
	req, err := c.makeRequest(ctx, link)
	if err != nil {
		return nil, nil, err
	}
	res, err := c.doGetRequest(req)
	if err != nil {
		return nil, nil, err
	}
	defer utils.DeferCloseReader(res.Body)
	links, hlsQuery, err := utils.DecodeRaw(res.Body, psch)
	if err != nil {
		return nil, nil, err
	}
	return c.filterParts(links), hlsQuery, nil
}

func (c *ClientAPI) Download(ctx context.Context, f *os.File, url string) error {
	if c.writedParts.Has(url) {
		return nil
	}
	if err := c.download(ctx, url, f); err != nil {
		return err
	}
	c.writedParts.Set(url)
	return nil
}

func (c *ClientAPI) download(ctx context.Context, link string, f *os.File) error {
	req, err := c.makeRequest(ctx, link)
	if err != nil {
		return err
	}
	res, err := c.doGetRequest(req)
	if err != nil {
		return err
	}
	defer utils.DeferCloseReader(res.Body)
	b, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}
	if err := res.Body.Close(); err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		return err
	}
	return nil
}

func (c *ClientAPI) startPlaylistLoop(
	ctx context.Context,
	ch chan<- string,
	plist string,
	pkey utils.PSCH,
) {
	defer close(ch)
	var currentTry = 0
	const maxRetry = 30
	const timeout = (time.Second * 1) + (time.Millisecond * 431)
	ticker := time.NewTicker(timeout)
	defer ticker.Stop()
	var more url.Values

loop:
	for {
		select {
		case <-ctx.Done():
			slog.Debug("context done")
			break loop
		case <-ticker.C:
			vids, hlsQuery, err := c.GetPlaylistVideo(ctx, plist, pkey, more)
			if len(hlsQuery) > 0 {
				more = hlsQuery
			}
			if err != nil {
				if !utils.IsCancel(err) {
					slog.Error("GetPlaylistVideo", "error", err.Error())
				}
				break loop
			}
			for i := range vids {
				ch <- vids[i]
			}
			if len(vids) == 0 {
				currentTry++
				if currentTry >= maxRetry {
					slog.Error("GetPlaylistVideo", "error", "may be offline, max retries")
					break loop
				}
			} else {
				currentTry = 0
			}
		}
	}
}

func (c *ClientAPI) StartPlaylistLoop(ctx context.Context, plist string, pkey utils.PSCH) <-chan string {
	var ch = make(chan string, 100)
	go c.startPlaylistLoop(ctx, ch, plist, pkey)
	return ch
}
