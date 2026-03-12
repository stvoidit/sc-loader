package client

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"net/url"
	"sync"
	"time"
)

type SrvCDN struct {
	Domain  string
	TTFB    time.Duration
	Healthy bool
	Err     error
}

func (cdn SrvCDN) String() string {
	return (&url.URL{
		Scheme: "https",
		Host:   "edge-hls." + cdn.Domain,
	}).String()
}

func (c *ClientAPI) probeHTTPS(ctx context.Context, cdn *SrvCDN, timeout time.Duration) {
	reqCtx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()
	req, _ := http.NewRequestWithContext(reqCtx, http.MethodConnect, cdn.String(), nil)
	start := time.Now()
	resp, err := c._c.Do(req)
	ttfb := time.Since(start)

	if err != nil {
		cdn.TTFB += timeout
		cdn.Err = err
		return
	}
	if err := resp.Body.Close(); err != nil {
		cdn.TTFB += timeout
		cdn.Err = err
		return
	}
	cdn.TTFB = ttfb
	cdn.Healthy = resp.StatusCode <= 404
	if !cdn.Healthy {
		cdn.TTFB += timeout
		cdn.Err = fmt.Errorf("http %d", resp.StatusCode)
	}
}

func (c ClientAPI) PickBestHTTPS(ctx context.Context, domains []string) string {
	const cycles = 3
	timeout := 2 * time.Second
	var wg sync.WaitGroup
	out := make([]SrvCDN, len(domains))
	for i := range domains {
		out[i] = SrvCDN{Domain: domains[i]}
		wg.Go(func() {
			for range cycles {
				c.probeHTTPS(ctx, &out[i], timeout)
			}
			slog.Debug("probeHTTPS",
				slog.String("Domain", out[i].Domain),
				slog.String("URL", out[i].String()),
				slog.Duration("Duration", out[i].TTFB/cycles),
				slog.Int64("Microseconds", out[i].TTFB.Microseconds()/cycles),
				slog.Any("error", out[i].Err))
		})
	}
	wg.Wait()

	var best SrvCDN
	best.TTFB = timeout
	for _, r := range out {
		if !r.Healthy || r.TTFB == 0 {
			continue
		}
		if r.TTFB < best.TTFB {
			best = r
		}
	}
	if best.Domain == "" {
		for _, r := range out {
			if r.TTFB > 0 && r.TTFB < best.TTFB {
				best = r
			}
		}
	}
	slog.Info("PickBestHTTPS", slog.Any("best", best))
	return best.Domain
}
