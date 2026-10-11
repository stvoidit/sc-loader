package client

import (
	"context"
	"errors"
	"net/http"
	"testing"
	"time"

	"sc-loader/utils"
)

func TestSrvCDNString(t *testing.T) {
	cdn := SrvCDN{Domain: "doppiocdn.com"}
	if got, want := cdn.String(), "https://edge-hls.doppiocdn.com"; got != want {
		t.Errorf("String() = %q, want %q", got, want)
	}
}

func TestPickBestHTTPS(t *testing.T) {
	ctx := context.Background()

	t.Run("no domains", func(t *testing.T) {
		c := NewClient(&utils.Config{Host: "example.com"})
		if got := c.PickBestHTTPS(ctx, nil); got != "" {
			t.Errorf("PickBestHTTPS(nil) = %q, want empty", got)
		}
	})

	t.Run("transport error", func(t *testing.T) {
		c := NewClient(&utils.Config{Host: "example.com"})
		c._c.Transport = roundTripperFunc(func(*http.Request) (*http.Response, error) {
			return nil, errors.New("dial tcp: no route to host")
		})
		if got := c.PickBestHTTPS(ctx, []string{"a.example", "b.example"}); got != "" {
			t.Errorf("PickBestHTTPS() = %q, want empty", got)
		}
	})

	t.Run("all unhealthy", func(t *testing.T) {
		c, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusInternalServerError)
		}))
		if got := c.PickBestHTTPS(ctx, []string{"a.example"}); got != "" {
			t.Errorf("PickBestHTTPS() = %q, want empty для http 500", got)
		}
	})

	t.Run("picks fastest healthy", func(t *testing.T) {
		c, _ := newTestClient(t, http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.Host == "edge-hls.slow.example" {
				time.Sleep(300 * time.Millisecond)
			}
			w.WriteHeader(http.StatusNotFound) // 404 <= 404, значит хост жив
		}))

		got := c.PickBestHTTPS(ctx, []string{"slow.example", "fast.example"})
		if got != "fast.example" {
			t.Errorf("PickBestHTTPS() = %q, want fast.example", got)
		}
	})
}
