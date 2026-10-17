package client

import (
	"bytes"
	"compress/gzip"
	"compress/zlib"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strings"
	"testing"
)

func gzipBytes(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := gzip.NewWriter(&buf)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func deflateBytes(t *testing.T, s string) []byte {
	t.Helper()
	var buf bytes.Buffer
	w := zlib.NewWriter(&buf)
	if _, err := w.Write([]byte(s)); err != nil {
		t.Fatal(err)
	}
	if err := w.Close(); err != nil {
		t.Fatal(err)
	}
	return buf.Bytes()
}

func newResponse(body []byte, encoding string) *http.Response {
	h := http.Header{}
	if encoding != "" {
		h.Set("Content-Encoding", encoding)
	}
	return &http.Response{
		StatusCode: http.StatusOK,
		Header:     h,
		Body:       io.NopCloser(bytes.NewReader(body)),
	}
}

func TestUnzipResponse(t *testing.T) {
	const plain = "hello world"

	t.Run("no encoding", func(t *testing.T) {
		res := newResponse([]byte(plain), "")
		if err := unzipResponse(res); err != nil {
			t.Fatalf("unzipResponse() error = %v", err)
		}
		got, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != plain {
			t.Errorf("body = %q, want %q", got, plain)
		}
		if res.Uncompressed {
			t.Error("Uncompressed = true, want false")
		}
	})

	t.Run("gzip", func(t *testing.T) {
		res := newResponse(gzipBytes(t, plain), "gzip")
		if err := unzipResponse(res); err != nil {
			t.Fatalf("unzipResponse() error = %v", err)
		}
		defer res.Body.Close()
		got, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != plain {
			t.Errorf("body = %q, want %q", got, plain)
		}
		if !res.Uncompressed {
			t.Error("Uncompressed = false, want true")
		}
		if res.Header.Get("Content-Encoding") != "" {
			t.Error("Content-Encoding не удалён")
		}
	})

	t.Run("deflate", func(t *testing.T) {
		res := newResponse(deflateBytes(t, plain), "deflate")
		if err := unzipResponse(res); err != nil {
			t.Fatalf("unzipResponse() error = %v", err)
		}
		defer res.Body.Close()
		got, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != plain {
			t.Errorf("body = %q, want %q", got, plain)
		}
		if !res.Uncompressed {
			t.Error("Uncompressed = false, want true")
		}
	})

	t.Run("unknown encoding", func(t *testing.T) {
		res := newResponse([]byte(plain), "br")
		if err := unzipResponse(res); !errors.Is(err, ErrUnknownCompressionMethod) {
			t.Errorf("unzipResponse() error = %v, want %v", err, ErrUnknownCompressionMethod)
		}
	})

	t.Run("broken gzip", func(t *testing.T) {
		res := newResponse([]byte("definitely not gzip"), "gzip")
		if err := unzipResponse(res); err == nil {
			t.Error("unzipResponse() error = nil, want error")
		}
	})

	t.Run("already uncompressed", func(t *testing.T) {
		res := newResponse(gzipBytes(t, plain), "gzip")
		res.Uncompressed = true
		if err := unzipResponse(res); err != nil {
			t.Fatalf("unzipResponse() error = %v", err)
		}
		got, err := io.ReadAll(res.Body)
		if err != nil {
			t.Fatal(err)
		}
		if bytes.Equal(got, []byte(plain)) {
			t.Error("тело было распаковано повторно")
		}
		if res.Header.Get("Content-Encoding") != "gzip" {
			t.Error("Content-Encoding удалён, хотя распаковки не было")
		}
	})
}

func TestMakeURLs(t *testing.T) {
	tests := []struct {
		name string
		got  string
		want string
	}{
		{
			name: "config",
			got:  makeConfigURL("example.com"),
			want: "https://example.com/api/front/v3/config/initial-dynamic?requestPath=/",
		},
		{
			name: "broadcast",
			got:  makeBroadcastURL("example.com", 123),
			want: "https://example.com/api/front/v2/broadcasts/123",
		},
		{
			name: "user id",
			got:  makeUserIDURL("example.com", "alex"),
			want: "https://example.com/api/front/users/user-ids/alex",
		},
		{
			name: "master playlist",
			got:  makeMasterPlaylistURL("doppiocdn.com", 123),
			want: "https://edge-hls.doppiocdn.com/hls/123/master/123_auto.m3u8?playlistType=lowLatency",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.got != tt.want {
				t.Errorf("got %q, want %q", tt.got, tt.want)
			}
		})
	}
}

func TestMakeRoomURL(t *testing.T) {
	got := makeRoomURL("example.com", "alex")
	if !strings.HasPrefix(got, "https://example.com/api/front/v2/models/username/alex/cam?") {
		t.Fatalf("makeRoomURL() = %q", got)
	}
	q, err := url.ParseQuery(strings.SplitN(got, "?", 2)[1])
	if err != nil {
		t.Fatal(err)
	}
	if q.Get("triggerRequest") != "loadCam" {
		t.Errorf("triggerRequest = %q, want %q", q.Get("triggerRequest"), "loadCam")
	}
	if len(q.Get("uniq")) != 16 {
		t.Errorf("uniq = %q, want 16 символов", q.Get("uniq"))
	}
}

func TestMakeMasterPlaylist(t *testing.T) {
	got, err := makeMasterPlaylist(
		"https://edge-hls.example.com/hls/123/master/123_auto.m3u8",
		"PKEY",
	)
	if err != nil {
		t.Fatalf("makeMasterPlaylist() error = %v", err)
	}
	u, err := url.Parse(got)
	if err != nil {
		t.Fatal(err)
	}
	q := u.Query()
	if q.Get("pkey") != "PKEY" || q.Get("psch") != "v2" ||
		q.Get("playlistType") != playlistTypeLowLatency {
		t.Errorf("query = %v, want pkey/psch/playlistType", q)
	}

	t.Run("keeps existing query", func(t *testing.T) {
		got, err := makeMasterPlaylist(
			"https://example.com/a.m3u8?playlistType=lowLatency&x=1",
			"K",
		)
		if err != nil {
			t.Fatal(err)
		}
		q, err := url.ParseQuery(strings.SplitN(got, "?", 2)[1])
		if err != nil {
			t.Fatal(err)
		}
		if q.Get("x") != "1" || q.Get("pkey") != "K" {
			t.Errorf("query = %v, want x=1&pkey=K", q)
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		if _, err := makeMasterPlaylist("https://example.com/%zz", "K"); err == nil {
			t.Error("makeMasterPlaylist() error = nil, want error")
		}
	})
}

func TestMakeVideoPlaylist(t *testing.T) {
	base := "https://edge-hls.example.com/hls/123/master/media.m3u8?pkey=K"
	more := url.Values{"_HLS_msn": {"5"}, "_HLS_part": {"2"}}

	got, err := makeVideoPlaylist(base, more)
	if err != nil {
		t.Fatalf("makeVideoPlaylist() error = %v", err)
	}
	q, err := url.ParseQuery(strings.SplitN(got, "?", 2)[1])
	if err != nil {
		t.Fatal(err)
	}
	if q.Get("pkey") != "K" || q.Get("_HLS_msn") != "5" || q.Get("_HLS_part") != "2" {
		t.Errorf("query = %v, want pkey=K&_HLS_msn=5&_HLS_part=2", q)
	}

	t.Run("more overwrites base", func(t *testing.T) {
		got, err := makeVideoPlaylist(base, url.Values{"pkey": {"OTHER"}})
		if err != nil {
			t.Fatal(err)
		}
		q, err := url.ParseQuery(strings.SplitN(got, "?", 2)[1])
		if err != nil {
			t.Fatal(err)
		}
		if q.Get("pkey") != "OTHER" {
			t.Errorf("pkey = %q, want OTHER", q.Get("pkey"))
		}
	})

	t.Run("nil more", func(t *testing.T) {
		got, err := makeVideoPlaylist(base, nil)
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(got, "pkey=K") {
			t.Errorf("makeVideoPlaylist() = %q, want pkey=K", got)
		}
	})

	t.Run("invalid url", func(t *testing.T) {
		got, err := makeVideoPlaylist("https://example.com/%zz", nil)
		if err == nil {
			t.Fatal("makeVideoPlaylist() error = nil, want error")
		}
		if got != "https://example.com/%zz" {
			t.Errorf("makeVideoPlaylist() = %q, want исходную ссылку", got)
		}
	})
}
