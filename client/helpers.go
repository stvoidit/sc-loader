package client

import (
	"bytes"
	"compress/gzip"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sc-loader/utils"
)

var ErrUnknownCompressionMethod = errors.New("unknown compression method")

// unzipResponse - распаковка сжатого ответа
func unzipResponse(response *http.Response) (err error) {
	if response.Uncompressed {
		return nil
	}
	switch response.Header.Get("Content-Encoding") {
	case "":
		return nil
	case "gzip":
		b, err := io.ReadAll(response.Body)
		if err != nil {
			return err
		}
		if err := response.Body.Close(); err != nil {
			return err
		}
		gz, err := gzip.NewReader(bytes.NewReader(b))
		if err != nil {
			return err
		}
		response.Body = io.NopCloser(gz)
		response.Header.Del("Content-Encoding")
		response.Uncompressed = true
		return nil
	default:
		return ErrUnknownCompressionMethod
	}
}

const (
	scheme           = "https"
	endpointConfig   = `/api/front/v3/config/initial-dynamic`
	endpointRoom     = `/api/front/v2/models/username/%s/cam`
	endpointPlayList = `/hls/%d/master/%d_auto.m3u8`
)

const (
	playlistTypeStandard   = "standard"
	playlistTypeLowLatency = "lowLatency"
)

var (
	lowLatencyQuery = (url.Values{
		"playlistType": []string{playlistTypeLowLatency},
	}).Encode()
	roomQuery = (url.Values{
		"triggerRequest": []string{"loadCam"},
		"uniq":           []string{utils.MakeUniq()},
	}).Encode()
)

func makeConfigURL(host string) string {
	link := url.URL{
		Scheme: scheme,
		Host:   host,
		Path:   endpointConfig,
	}
	return link.String()
}

func makeRoomURL(host, username string) string {
	// uniq := utils.MakeUniq()
	link := (url.URL{
		Scheme:   scheme,
		Host:     host,
		Path:     fmt.Sprintf(endpointRoom, username),
		RawQuery: roomQuery,
	})
	return link.String()
}

func makeMasterPlaylistURL(cdnDomain string, roomId int) string {
	link := url.URL{
		Scheme:   scheme,
		Host:     "edge-hls." + cdnDomain,
		Path:     fmt.Sprintf(endpointPlayList, roomId, roomId),
		RawQuery: lowLatencyQuery,
	}
	return link.String()
}

func makeMasterPlaylist(link, pkey string) (string, error) {
	URL, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	q := URL.Query()
	q.Set("pkey", pkey)
	q.Set("psch", "v2")
	q.Set("playlistType", playlistTypeLowLatency)
	URL.RawQuery = q.Encode()
	return URL.String(), nil
}

func makeVideoPlaylist(videoList string, more url.Values) (string, error) {
	link, err := url.Parse(videoList)
	if err != nil {
		return videoList, err
	}
	query := link.Query()
	for k, v := range more {
		query.Set(k, v[0])
	}
	link.RawQuery = query.Encode()
	return link.String(), nil
}
