package client

import (
	"context"
	"crypto/sha256"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"sc-loader/models"
	"sc-loader/utils"
)

// TestMain глушит slog-вывод тестируемого кода, чтобы не засорять лог теста.
func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

// Реальная пара ключ/шифртекст DRM mouflon (см. tests/main.js).
const (
	testKey       = "EQueeGh2kaewa3ch"
	testEncrypted = "gqaq55h8zoL+Lnlg9nmPX7"
	testDecrypted = "CblB7Ti4PVvMt24g"
	testRoomID    = 123
)

func testPSCH() utils.PSCH { return sha256.Sum256([]byte(testKey)) }

// mediaPart собирает ссылку на part-сегмент, где зашифрованное имя стоит
// третьим с конца элементом пути (именно так его ищет utils.DecodePart).
func mediaPart(host, encrypted string) string {
	return fmt.Sprintf("https://%s/hls/%d/master/%d_auto/init_a1b2_part_%s_7_9.m4s",
		host, testRoomID, testRoomID, encrypted)
}

// newTestClient поднимает TLS-сервер и возвращает ClientAPI, все запросы
// которого уходят на этот сервер независимо от host в ссылке.
func newTestClient(t *testing.T, handler http.Handler) (*ClientAPI, *utils.Config) {
	t.Helper()
	srv := httptest.NewTLSServer(handler)
	t.Cleanup(srv.Close)

	u, err := url.Parse(srv.URL)
	if err != nil {
		t.Fatal(err)
	}
	cnf := &utils.Config{
		Host:   u.Host,
		Folder: t.TempDir(),
		Keys:   map[string]string{testKey: testKey},
	}
	cnf.CacheKeys = make(utils.KeysDRM, len(cnf.Keys))
	cnf.CacheKeys.LoadKeys(cnf.Keys)

	c := NewClient(cnf)
	c._c.Transport = &http.Transport{
		DialTLSContext: func(ctx context.Context, network, _ string) (net.Conn, error) {
			var d net.Dialer
			raw, err := d.DialContext(ctx, network, u.Host)
			if err != nil {
				return nil, err
			}
			conn := tls.Client(raw, &tls.Config{InsecureSkipVerify: true})
			if err := conn.HandshakeContext(ctx); err != nil {
				utils.DeferCloseReader(raw)
				return nil, err
			}
			return conn, nil
		},
	}
	return c, cnf
}

// jsonHandler отдаёт фиксированное тело по указанному пути.
func jsonHandler(path, body string) *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc(path, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		io.WriteString(w, body) //nolint:errcheck
	})
	return mux
}

func TestNewClient(t *testing.T) {
	cnf := &utils.Config{Host: "example.com"}
	c := NewClient(cnf)

	if c.defaultHeaders.Get("Accept") != "application/json" {
		t.Errorf("Accept = %q, want application/json", c.defaultHeaders.Get("Accept"))
	}
	if c.defaultHeaders.Get("User-Agent") != ua {
		t.Errorf("User-Agent = %q, want %q", c.defaultHeaders.Get("User-Agent"), ua)
	}
	if c.htmlHeaders.Get("Origin") != "https://example.com" {
		t.Errorf("Origin = %q", c.htmlHeaders.Get("Origin"))
	}
	if c.htmlHeaders.Get("Referer") != "https://example.com/" {
		t.Errorf("Referer = %q", c.htmlHeaders.Get("Referer"))
	}
	if c._c == nil || c._c.Timeout != 10*time.Second {
		t.Errorf("http.Client = %+v, want timeout 10s", c._c)
	}
	if !c.IsNewURL("x") {
		t.Error("IsNewURL() = false для нового множества")
	}
}

func TestMakeRequest(t *testing.T) {
	c := NewClient(&utils.Config{Host: "example.com"})
	ctx := context.Background()

	req, err := c.makeRequest(ctx, "https://example.com/api")
	if err != nil {
		t.Fatalf("makeRequest() error = %v", err)
	}
	if req.Method != http.MethodGet {
		t.Errorf("Method = %q, want GET", req.Method)
	}
	if req.Header.Get("User-Agent") != ua {
		t.Errorf("User-Agent = %q, want %q", req.Header.Get("User-Agent"), ua)
	}
	if req.Context() != ctx {
		t.Error("контекст не проброшен в запрос")
	}

	// заголовки запроса - копия дефолтных, изменение не должно их затрагивать
	req.Header.Set("Accept", "text/plain")
	if c.defaultHeaders.Get("Accept") != "application/json" {
		t.Error("makeRequest() изменил defaultHeaders")
	}

	if _, err := c.makeRequest(ctx, "https://example.com/%zz"); err == nil {
		t.Error("makeRequest() error = nil, want error")
	}
}

func TestDoGetRequestError(t *testing.T) {
	c := NewClient(&utils.Config{Host: "example.com"})
	c._c.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("dial failed")
	})
	req, err := c.makeRequest(context.Background(), "https://example.com/api")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := c.doGetRequest(req); err == nil {
		t.Error("doGetRequest() error = nil, want error")
	}
}

type roundTripperFunc func(*http.Request) (*http.Response, error)

func (f roundTripperFunc) RoundTrip(r *http.Request) (*http.Response, error) { return f(r) }

func TestGetRoomID(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		var gotPath string
		var gotHeader http.Header
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/api/front/users/user-ids/alex",
			func(w http.ResponseWriter, r *http.Request) {
				gotPath, gotHeader = r.URL.Path, r.Header.Clone()
				w.Header().Set("Content-Type", "application/json")
				io.WriteString(w, `{"id":11111111}`) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		streamer := &models.Streamer{Username: "alex"}
		id, online, err := c.GetRoomID(ctx, streamer)
		if err != nil {
			t.Fatalf("GetRoomID() error = %v", err)
		}
		if id != 11111111 || !online {
			t.Errorf("GetRoomID() = (%d, %v), want (11111111, true)", id, online)
		}
		if streamer.RoomID != 11111111 {
			t.Errorf("streamer.RoomID = %d, want 11111111", streamer.RoomID)
		}
		if gotPath != "/api/front/users/user-ids/alex" {
			t.Errorf("path = %q", gotPath)
		}
		if gotHeader.Get("User-Agent") != ua {
			t.Errorf("User-Agent = %q, want %q", gotHeader.Get("User-Agent"), ua)
		}
		if gotHeader.Get("Accept") != "application/json" {
			t.Errorf("Accept = %q, want application/json", gotHeader.Get("Accept"))
		}
	})

	t.Run("gzip response", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/api/front/users/user-ids/alex",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Encoding", "gzip")
				w.Write(gzipBytes(t, `{"id":777}`)) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		id, _, err := c.GetRoomID(ctx, &models.Streamer{Username: "alex"})
		if err != nil {
			t.Fatalf("GetRoomID() error = %v", err)
		}
		if id != 777 {
			t.Errorf("GetRoomID() id = %d, want 777", id)
		}
	})

	t.Run("no id in response", func(t *testing.T) {
		mux := jsonHandler("/api/front/users/user-ids/alex", `{}`)
		c, _ := newTestClient(t, mux)

		id, online, err := c.GetRoomID(ctx, &models.Streamer{Username: "alex"})
		if err == nil {
			t.Fatal("GetRoomID() error = nil, want error")
		}
		if id != -1 || online {
			t.Errorf("GetRoomID() = (%d, %v), want (-1, false)", id, online)
		}
	})

	t.Run("unexpected json shape", func(t *testing.T) {
		mux := jsonHandler("/api/front/users/user-ids/alex", `{"id":"not-a-number"}`)
		c, _ := newTestClient(t, mux)

		if _, _, err := c.GetRoomID(ctx, &models.Streamer{Username: "alex"}); err == nil {
			t.Error("GetRoomID() error = nil, want unmarshal error")
		}
	})

	t.Run("broken json", func(t *testing.T) {
		mux := jsonHandler("/api/front/users/user-ids/alex", `<html>not found</html>`)
		c, _ := newTestClient(t, mux)

		if _, _, err := c.GetRoomID(ctx, &models.Streamer{Username: "alex"}); err == nil {
			t.Error("GetRoomID() error = nil, want unmarshal error")
		}
	})

	t.Run("http 500", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/api/front/users/user-ids/alex",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusInternalServerError)
				io.WriteString(w, "internal error") //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		id, online, err := c.GetRoomID(ctx, &models.Streamer{Username: "alex"})
		if err == nil {
			t.Fatal("GetRoomID() error = nil для http 500")
		}
		if id != -1 || online {
			t.Errorf("GetRoomID() = (%d, %v), want (-1, false)", id, online)
		}
	})

	t.Run("canceled context", func(t *testing.T) {
		mux := jsonHandler("/api/front/users/user-ids/alex", `{"id":1}`)
		c, _ := newTestClient(t, mux)

		canceled, cancel := context.WithCancel(ctx)
		cancel()
		_, _, err := c.GetRoomID(canceled, &models.Streamer{Username: "alex"})
		if err == nil || !utils.IsCancel(err) {
			t.Errorf("GetRoomID() error = %v, want canceled", err)
		}
	})
}

func TestGetRoomStatus(t *testing.T) {
	ctx := context.Background()

	t.Run("online", func(t *testing.T) {
		var gotPath string
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/api/front/v2/broadcasts/123",
			func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				io.WriteString(
					w,
					`{"item":{"isLive":true,"status":"public","modelId":5}}`,
				) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		status, online, err := c.GetRoomStatus(ctx, &models.Streamer{RoomID: testRoomID})
		if err != nil {
			t.Fatalf("GetRoomStatus() error = %v", err)
		}
		if status != models.StatusPublic || !online {
			t.Errorf("GetRoomStatus() = (%s, %v), want (public, true)", status, online)
		}
		if gotPath != "/api/front/v2/broadcasts/123" {
			t.Errorf("path = %q", gotPath)
		}
	})

	t.Run("offline", func(t *testing.T) {
		mux := jsonHandler(
			"/api/front/v2/broadcasts/123",
			`{"item":{"isLive":false,"status":"offline"}}`,
		)
		c, _ := newTestClient(t, mux)

		status, online, err := c.GetRoomStatus(ctx, &models.Streamer{RoomID: testRoomID})
		if err != nil {
			t.Fatalf("GetRoomStatus() error = %v", err)
		}
		if status != "offline" || online {
			t.Errorf("GetRoomStatus() = (%s, %v), want (offline, false)", status, online)
		}
	})

	t.Run("broken json", func(t *testing.T) {
		mux := jsonHandler("/api/front/v2/broadcasts/123", `{"item":`)
		c, _ := newTestClient(t, mux)

		status, online, err := c.GetRoomStatus(ctx, &models.Streamer{RoomID: testRoomID})
		if err == nil {
			t.Fatal("GetRoomStatus() error = nil, want error")
		}
		// при ошибке разбора online остаётся true, чтобы запись не прерывалась
		if status != "" || !online {
			t.Errorf("GetRoomStatus() = (%s, %v), want (\"\", true)", status, online)
		}
	})

	t.Run("http 404 with json body", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/api/front/v2/broadcasts/123",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `{"error":"not found"}`) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		// код ответа не проверяется, поэтому тело разбирается как Broadcast
		status, online, err := c.GetRoomStatus(ctx, &models.Streamer{RoomID: testRoomID})
		if err != nil {
			t.Fatalf("GetRoomStatus() error = %v", err)
		}
		if status != "" || online {
			t.Errorf("GetRoomStatus() = (%s, %v), want (\"\", false)", status, online)
		}
	})

	t.Run("http 404 with html body", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/api/front/v2/broadcasts/123",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusNotFound)
				io.WriteString(w, `<html>not found</html>`) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		if _, _, err := c.GetRoomStatus(ctx, &models.Streamer{RoomID: testRoomID}); err == nil {
			t.Error("GetRoomStatus() error = nil, want unmarshal error")
		}
	})
}

const initialDynamicJSON = `{
	"initialDynamic": {
		"csrfToken": "csrf",
		"guestId": 42,
		"players": {
			"isHLSLLEnabled": true,
			"cdnConfig": [
				{"domain": "doppiocdn.com", "features": {"ll": true}},
				{"domain": "doppiocdn2.com", "features": {"ll": true}}
			]
		}
	}
}`

func TestInitClientConfig(t *testing.T) {
	ctx := context.Background()

	t.Run("ok", func(t *testing.T) {
		var configHeader http.Header
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/" {
				return
			}
			w.Header().Set("Content-Type", "text/html")
			io.WriteString(
				w,
				"<html>\n<script>\nrelease_version: 2026.10.11\n</script>\n</html>",
			) //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v3/config/initial-dynamic",
			func(w http.ResponseWriter, r *http.Request) {
				configHeader = r.Header.Clone()
				io.WriteString(w, initialDynamicJSON) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		id, err := c.InitClientConfig(ctx)
		if err != nil {
			t.Fatalf("InitClientConfig() error = %v", err)
		}
		if id.CSRFToken != "csrf" || id.GuestID != 42 {
			t.Errorf("InitialDynamic = %+v", id)
		}
		if got := c.defaultHeaders.Get("front-version"); got != "2026.10.11" {
			t.Errorf("front-version = %q, want 2026.10.11", got)
		}
		if configHeader.Get("front-version") != "2026.10.11" {
			t.Errorf("запрос конфига ушёл без front-version: %v", configHeader)
		}
		if got := id.GetDomainsCDN(); len(got) != 1 || got[0] != "doppiocdn.com" {
			t.Errorf("GetDomainsCDN() = %v, want [doppiocdn.com]", got)
		}
	})

	t.Run("Init returns cdn domain", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `release_version: 1.0`) //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v3/config/initial-dynamic",
			func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, initialDynamicJSON) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		domain, err := c.Init(ctx)
		if err != nil {
			t.Fatalf("Init() error = %v", err)
		}
		if domain != "doppiocdn.com" {
			t.Errorf("Init() = %q, want doppiocdn.com", domain)
		}
	})

	t.Run("no release version in html", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `<html>nothing here</html>`) //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v3/config/initial-dynamic",
			func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, initialDynamicJSON) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		if _, err := c.InitClientConfig(ctx); err != nil {
			t.Fatalf("InitClientConfig() error = %v", err)
		}
		if got := c.defaultHeaders.Get("front-version"); got != "" {
			t.Errorf("front-version = %q, want empty", got)
		}
	})

	t.Run("x-api-version header wins", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `release_version: 1.0`) //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v3/config/initial-dynamic",
			func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("X-Api-Version", "v9")
				io.WriteString(w, initialDynamicJSON) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		if _, err := c.InitClientConfig(ctx); err != nil {
			t.Fatalf("InitClientConfig() error = %v", err)
		}
		if got := c.defaultHeaders.Get("front-version"); got != "v9" {
			t.Errorf("front-version = %q, want v9", got)
		}
	})

	t.Run("config endpoint broken json", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, `release_version: 1.0`) //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v3/config/initial-dynamic",
			func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusServiceUnavailable)
				io.WriteString(w, `service unavailable`) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)

		if _, err := c.InitClientConfig(ctx); err == nil {
			t.Error("InitClientConfig() error = nil, want error")
		}
	})

	t.Run("html request failed", func(t *testing.T) {
		c := NewClient(&utils.Config{Host: "example.com"})
		c._c.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("no route to host")
		})
		if _, err := c.InitClientConfig(ctx); err == nil {
			t.Error("InitClientConfig() error = nil, want error")
		}
	})
}

func masterPlaylistBody(variant, pkey string) string {
	return strings.Join([]string{
		"#EXTM3U",
		"#EXT-X-MOUFLON:PSCH:v2:" + pkey,
		"#EXT-X-STREAM-INF:BANDWIDTH=2500000,RESOLUTION=1280x720",
		variant,
	}, "\n")
}

func TestGetPlaylistVariants(t *testing.T) {
	ctx := context.Background()
	const variant = "https://edge-hls.doppiocdn.com/hls/123/master/123_auto.m3u8"

	newSrv := func(t *testing.T, status int, body string) *ClientAPI {
		mux := http.NewServeMux()
		mux.HandleFunc(
			"/hls/123/master/123_auto.m3u8",
			func(w http.ResponseWriter, r *http.Request) {
				if status != http.StatusOK {
					w.WriteHeader(status)
				}
				io.WriteString(w, body) //nolint:errcheck
			},
		)
		c, _ := newTestClient(t, mux)
		return c
	}

	t.Run("ok", func(t *testing.T) {
		c := newSrv(t, http.StatusOK, masterPlaylistBody(variant, testKey))

		playlist, psch, err := c.GetPlaylistVariants(ctx, "doppiocdn.com", &models.Streamer{RoomID: testRoomID})
		if err != nil {
			t.Fatalf("GetPlaylistVariants() error = %v", err)
		}
		if psch != testPSCH() {
			t.Errorf("psch = %x, want %x", psch, testPSCH())
		}
		u, err := url.Parse(playlist)
		if err != nil {
			t.Fatal(err)
		}
		if u.Path != "/hls/123/master/123_auto.m3u8" {
			t.Errorf("playlist path = %q", u.Path)
		}
		q := u.Query()
		if q.Get("pkey") != testKey || q.Get("psch") != "v2" ||
			q.Get("playlistType") != playlistTypeLowLatency {
			t.Errorf("playlist query = %v", q)
		}
	})

	t.Run("http 404", func(t *testing.T) {
		c := newSrv(t, http.StatusNotFound, masterPlaylistBody(variant, testKey))

		_, _, err := c.GetPlaylistVariants(ctx, "doppiocdn.com", &models.Streamer{RoomID: testRoomID})
		if err == nil {
			t.Fatal("GetPlaylistVariants() error = nil для http 404")
		}
		if err.Error() != "404" {
			t.Errorf("error = %q, want \"404\"", err.Error())
		}
	})

	t.Run("pkey not found", func(t *testing.T) {
		c := newSrv(t, http.StatusOK, "#EXTM3U\n"+variant)

		if _, _, err := c.GetPlaylistVariants(
			ctx,
			"doppiocdn.com",
			&models.Streamer{RoomID: testRoomID},
		); !errors.Is(
			err,
			ErrPKeyNotFound,
		) {
			t.Errorf("error = %v, want %v", err, ErrPKeyNotFound)
		}
	})

	t.Run("unknown pkey", func(t *testing.T) {
		c := newSrv(t, http.StatusOK, masterPlaylistBody(variant, "unknown-key"))

		if _, _, err := c.GetPlaylistVariants(
			ctx,
			"doppiocdn.com",
			&models.Streamer{RoomID: testRoomID},
		); !errors.Is(
			err,
			utils.ErrKeyNotFound,
		) {
			t.Errorf("error = %v, want %v", err, utils.ErrKeyNotFound)
		}
	})
}

func mediaPlaylistBody(host string) string {
	return strings.Join([]string{
		"#EXTM3U",
		"#EXT-X-TARGETDURATION:2",
		`#EXT-X-MAP:URI="init.mp4"`,
		`#EXT-X-RENDITION-REPORT:URI="../720p/live.m3u8",LAST-MSN=100500,LAST-PART=3`,
		"#EXTINF:2.000,",
		"#EXT-X-MOUFLON:URI:" + mediaPart(host, testEncrypted),
	}, "\n")
}

func TestGetPlaylistVideo(t *testing.T) {
	ctx := context.Background()

	t.Run("decodes parts and remembers them", func(t *testing.T) {
		var gotQuery url.Values
		var c *ClientAPI
		var cnf *utils.Config
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			gotQuery = r.URL.Query()
			io.WriteString(w, mediaPlaylistBody(c.cnf.Host)) //nolint:errcheck
		})
		c, cnf = newTestClient(t, mux)

		plist := fmt.Sprintf(
			"https://%s/hls/%d/master/media.m3u8?pkey=%s",
			cnf.Host,
			testRoomID,
			testKey,
		)
		more := url.Values{"_HLS_msn": {"100500"}}

		vids, hlsQuery, err := c.GetPlaylistVideo(ctx, plist, testPSCH(), more)
		if err != nil {
			t.Fatalf("GetPlaylistVideo() error = %v", err)
		}
		want := []string{"init.mp4", mediaPart(cnf.Host, testDecrypted)}
		if len(vids) != len(want) {
			t.Fatalf("vids = %v, want %v", vids, want)
		}
		for i := range want {
			if vids[i] != want[i] {
				t.Errorf("vids[%d] = %q, want %q", i, vids[i], want[i])
			}
		}
		if hlsQuery.Get("_HLS_msn") != "100500" || hlsQuery.Get("_HLS_part") != "3" {
			t.Errorf("hlsQuery = %v", hlsQuery)
		}
		if gotQuery.Get("_HLS_msn") != "100500" || gotQuery.Get("pkey") != testKey {
			t.Errorf("запрос ушёл с query = %v", gotQuery)
		}

		// после пометки сегмента как записанного он отфильтровывается
		c.SetCheckURL("init.mp4")
		vids, _, err = c.GetPlaylistVideo(ctx, plist, testPSCH(), nil)
		if err != nil {
			t.Fatalf("GetPlaylistVideo() error = %v", err)
		}
		if len(vids) != 1 || vids[0] != want[1] {
			t.Errorf("vids = %v, want [%s]", vids, want[1])
		}
	})

	t.Run("http 500 with empty body", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		})
		c, _ := newTestClient(t, mux)

		vids, hlsQuery, err := c.GetPlaylistVideo(
			ctx,
			"https://x/hls/123/master/media.m3u8",
			testPSCH(),
			nil,
		)
		if err != nil {
			t.Fatalf("GetPlaylistVideo() error = %v", err)
		}
		if len(vids) != 0 || hlsQuery != nil {
			t.Errorf("vids = %v, hlsQuery = %v, want пусто", vids, hlsQuery)
		}
	})

	t.Run("broken part in playlist", func(t *testing.T) {
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "#EXT-X-MOUFLON:URI:"+mediaPart("127.0.0.1", "!!!")) //nolint:errcheck
		})
		c, _ := newTestClient(t, mux)

		if _, _, err := c.GetPlaylistVideo(
			ctx,
			"https://x/hls/123/master/media.m3u8",
			testPSCH(),
			nil,
		); err == nil {
			t.Error("GetPlaylistVideo() error = nil, want decode error")
		}
	})

	t.Run("request error", func(t *testing.T) {
		c := NewClient(&utils.Config{Host: "example.com"})
		c._c.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("connection reset by peer")
		})
		if _, _, err := c.GetPlaylistVideo(
			ctx,
			"https://x/media.m3u8",
			testPSCH(),
			nil,
		); err == nil {
			t.Error("GetPlaylistVideo() error = nil, want error")
		}
	})

	t.Run("invalid playlist url", func(t *testing.T) {
		c, _ := newTestClient(t, http.NotFoundHandler())
		if _, _, err := c.GetPlaylistVideo(ctx, "https://x/%zz", testPSCH(), nil); err == nil {
			t.Error("GetPlaylistVideo() error = nil, want error")
		}
	})
}

func TestIsNewURLAndSetCheckURL(t *testing.T) {
	c := NewClient(&utils.Config{Host: "example.com"})
	if !c.IsNewURL("a") {
		t.Error("IsNewURL(a) = false, want true")
	}
	c.SetCheckURL("a")
	if c.IsNewURL("a") {
		t.Error("IsNewURL(a) = true после SetCheckURL")
	}
	if got := c.filterParts([]string{"a", "b"}); len(got) != 1 || got[0] != "b" {
		t.Errorf("filterParts() = %v, want [b]", got)
	}
}

func TestDownloadBuf(t *testing.T) {
	ctx := context.Background()

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/part1":
			io.WriteString(w, "PART-1-DATA") //nolint:errcheck
		case "/missing":
			w.WriteHeader(http.StatusNotFound)
			io.WriteString(w, "not found") //nolint:errcheck
		default:
			w.WriteHeader(http.StatusInternalServerError)
		}
	}))
	t.Cleanup(srv.Close)
	c := NewClient(&utils.Config{Host: "example.com"})

	t.Run("ok", func(t *testing.T) {
		link := srv.URL + "/part1"
		buf, err := c.DownloadBuf(ctx, link)
		if err != nil {
			t.Fatalf("DownloadBuf() error = %v", err)
		}
		if string(buf) != "PART-1-DATA" {
			t.Errorf("DownloadBuf() = %q", buf)
		}
		if c.IsNewURL(link) {
			t.Error("ссылка не помечена как загруженная")
		}
	})

	t.Run("http 404 body is returned as is", func(t *testing.T) {
		buf, err := c.DownloadBuf(ctx, srv.URL+"/missing")
		if err != nil {
			t.Fatalf("DownloadBuf() error = %v", err)
		}
		if string(buf) != "not found" {
			t.Errorf("DownloadBuf() = %q, want \"not found\"", buf)
		}
	})

	t.Run("request error", func(t *testing.T) {
		canceled, cancel := context.WithCancel(ctx)
		cancel()
		link := srv.URL + "/part2"
		if _, err := c.DownloadBuf(canceled, link); err == nil {
			t.Fatal("DownloadBuf() error = nil, want error")
		}
		// SetCheckURL вызывается через defer даже при ошибке
		if c.IsNewURL(link) {
			t.Error("ссылка не помечена после ошибки")
		}
	})

	t.Run("truncated body", func(t *testing.T) {
		// сервер обещает 100 байт и обрывает соединение
		broken := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Content-Length", "100")
			w.WriteHeader(http.StatusOK)
			io.WriteString(w, "short") //nolint:errcheck
			hj, ok := w.(http.Hijacker)
			if !ok {
				return
			}
			conn, _, err := hj.Hijack()
			if err != nil {
				return
			}
			conn.Close() //nolint:errcheck
		}))
		defer broken.Close()

		if _, err := c.DownloadBuf(ctx, broken.URL+"/part"); err == nil {
			t.Error("DownloadBuf() error = nil для оборванного ответа")
		}
	})
}

func TestDownloadWrite(t *testing.T) {
	ctx := context.Background()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "0123456789") //nolint:errcheck
	}))
	t.Cleanup(srv.Close)

	c := NewClient(&utils.Config{Host: "example.com"})
	path := filepath.Join(t.TempDir(), "out.mp4")
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0o600)
	if err != nil {
		t.Fatal(err)
	}
	defer utils.DeferCloseReader(f)

	link := srv.URL + "/part"
	n, err := c.DownloadWrite(ctx, f, link)
	if err != nil {
		t.Fatalf("DownloadWrite() error = %v", err)
	}
	if n != 10 {
		t.Errorf("DownloadWrite() n = %d, want 10", n)
	}

	// повторная запись той же ссылки пропускается
	n, err = c.DownloadWrite(ctx, f, link)
	if err != nil {
		t.Fatalf("DownloadWrite() error = %v", err)
	}
	if n != 0 {
		t.Errorf("DownloadWrite() n = %d, want 0", n)
	}

	got, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if string(got) != "0123456789" {
		t.Errorf("файл = %q, want %q", got, "0123456789")
	}

	// ошибка загрузки пробрасывается наружу
	if _, err := c.DownloadWrite(ctx, f, "https://x/%zz"); err == nil {
		t.Error("DownloadWrite() error = nil, want error")
	}
}

func TestStartPlaylistLoop(t *testing.T) {
	t.Run("parts are pushed to channel", func(t *testing.T) {
		var c *ClientAPI
		var cnf *utils.Config
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, mediaPlaylistBody(c.cnf.Host)) //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v2/broadcasts/123",
			func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `{"item":{"isLive":true,"status":"public"}}`) //nolint:errcheck
			},
		)
		c, cnf = newTestClient(t, mux)

		ctx, cancel := context.WithCancel(context.Background())
		defer cancel()

		plist := fmt.Sprintf(
			"https://%s/hls/%d/master/media.m3u8?pkey=%s",
			cnf.Host,
			testRoomID,
			testKey,
		)
		ch := c.StartPlaylistLoop(ctx, &models.Streamer{RoomID: testRoomID, Username: "alex"}, plist, testPSCH())

		want := mediaPart(cnf.Host, testDecrypted)
		deadline := time.After(15 * time.Second)
		select {
		case got, ok := <-ch:
			if !ok {
				t.Fatal("канал закрыт до получения сегмента")
			}
			if got != "init.mp4" {
				t.Errorf("первый сегмент = %q, want init.mp4", got)
			}
		case <-deadline:
			t.Fatal("таймаут ожидания сегмента")
		}

		select {
		case got, ok := <-ch:
			if !ok {
				t.Fatal("канал закрыт до получения part-сегмента")
			}
			if got != want {
				t.Errorf("сегмент = %q, want %q", got, want)
			}
		case <-deadline:
			t.Fatal("таймаут ожидания part-сегмента")
		}

		cancel()
		waitClosed(t, ch)
	})

	t.Run("offline room closes channel", func(t *testing.T) {
		var c *ClientAPI
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			// плейлист без новых сегментов
			io.WriteString(w, "#EXTM3U\n#EXTINF:2.0,\n") //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v2/broadcasts/123",
			func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `{"item":{"isLive":false,"status":"offline"}}`) //nolint:errcheck
			},
		)
		c, cnf := newTestClient(t, mux)

		ctx := t.Context()

		plist := fmt.Sprintf(
			"https://%s/hls/%d/master/media.m3u8?pkey=%s",
			cnf.Host,
			testRoomID,
			testKey,
		)
		ch := c.StartPlaylistLoop(ctx, &models.Streamer{RoomID: testRoomID, Username: "alex"}, plist, testPSCH())

		select {
		case link, ok := <-ch:
			if ok {
				t.Errorf("получен неожиданный сегмент %q", link)
			}
		case <-time.After(15 * time.Second):
			t.Fatal("канал не закрылся после offline-статуса")
		}
		waitClosed(t, ch)
	})

	t.Run("status request error closes channel", func(t *testing.T) {
		var c *ClientAPI
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "#EXTM3U\n#EXTINF:2.0,\n") //nolint:errcheck
		})
		mux.HandleFunc(
			"/api/front/v2/broadcasts/123",
			func(w http.ResponseWriter, r *http.Request) {
				io.WriteString(w, `<html>backend error</html>`) //nolint:errcheck
			},
		)
		c, cnf := newTestClient(t, mux)

		ctx := t.Context()

		plist := fmt.Sprintf(
			"https://%s/hls/%d/master/media.m3u8?pkey=%s",
			cnf.Host,
			testRoomID,
			testKey,
		)
		ch := c.StartPlaylistLoop(ctx, &models.Streamer{RoomID: testRoomID, Username: "alex"}, plist, testPSCH())
		waitClosed(t, ch)
	})

	t.Run("playlist error closes channel", func(t *testing.T) {
		var c *ClientAPI
		mux := http.NewServeMux()
		mux.HandleFunc("/hls/123/master/media.m3u8", func(w http.ResponseWriter, r *http.Request) {
			io.WriteString(w, "#EXT-X-MOUFLON:URI:"+mediaPart("127.0.0.1", "!!!")) //nolint:errcheck
		})
		c, cnf := newTestClient(t, mux)

		ctx := t.Context()

		plist := fmt.Sprintf(
			"https://%s/hls/%d/master/media.m3u8?pkey=%s",
			cnf.Host,
			testRoomID,
			testKey,
		)
		ch := c.StartPlaylistLoop(ctx, &models.Streamer{RoomID: testRoomID, Username: "alex"}, plist, testPSCH())
		waitClosed(t, ch)
	})
}

// waitClosed ждёт закрытия канала, но не дольше 15 секунд.
func waitClosed(t *testing.T, ch <-chan string) {
	t.Helper()
	deadline := time.After(15 * time.Second)
	for {
		select {
		case _, ok := <-ch:
			if !ok {
				return
			}
		case <-deadline:
			t.Fatal("канал не закрылся")
		}
	}
}

// TestInvalidHost проверяет, что ошибка создания запроса пробрасывается наружу.
func TestInvalidHost(t *testing.T) {
	ctx := context.Background()
	c := NewClient(&utils.Config{Host: "invalid host"})

	if _, _, err := c.GetRoomID(ctx, &models.Streamer{Username: "alex"}); err == nil {
		t.Error("GetRoomID() error = nil, want error")
	}
	if _, _, err := c.GetRoomStatus(ctx, &models.Streamer{RoomID: testRoomID}); err == nil {
		t.Error("GetRoomStatus() error = nil, want error")
	}
	if _, _, err := c.GetPlaylistVariants(ctx, "invalid cdn", &models.Streamer{RoomID: testRoomID}); err == nil {
		t.Error("GetPlaylistVariants() error = nil, want error")
	}
	if _, err := c.InitClientConfig(ctx); err == nil {
		t.Error("InitClientConfig() error = nil, want error")
	}
	if _, err := c.Init(ctx); err == nil {
		t.Error("Init() error = nil, want error")
	}
}

func TestInitError(t *testing.T) {
	c := NewClient(&utils.Config{Host: "example.com"})
	c._c.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
		return nil, errors.New("no route to host")
	})

	domain, err := c.Init(context.Background())
	if err == nil {
		t.Fatal("Init() error = nil, want error")
	}
	if domain != "" {
		t.Errorf("Init() = %q, want empty", domain)
	}
}

func TestGetFrontendVersionBrokenEncoding(t *testing.T) {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Encoding", "br")  // неподдерживаемое сжатие
		io.WriteString(w, "release_version: 1.0") //nolint:errcheck
	})
	c, _ := newTestClient(t, mux)

	if _, err := c.InitClientConfig(
		context.Background(),
	); !errors.Is(
		err,
		ErrUnknownCompressionMethod,
	) {
		t.Errorf("InitClientConfig() error = %v, want %v", err, ErrUnknownCompressionMethod)
	}
}

func TestDownloadWriteFileError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "data") //nolint:errcheck
	}))
	t.Cleanup(srv.Close)

	c := NewClient(&utils.Config{Host: "example.com"})
	path := filepath.Join(t.TempDir(), "closed.mp4")
	f, err := os.Create(path)
	if err != nil {
		t.Fatal(err)
	}
	if err := f.Close(); err != nil {
		t.Fatal(err)
	}

	if _, err := c.DownloadWrite(context.Background(), f, srv.URL+"/part"); err == nil {
		t.Error("DownloadWrite() error = nil для закрытого файла")
	}
}
