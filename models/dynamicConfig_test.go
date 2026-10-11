package models

import (
	"encoding/json/v2"
	"slices"
	"testing"
)

func cdn(domain string) CDNConfig {
	var c CDNConfig
	c.Domain = domain
	return c
}

func TestGetDomainsCDN(t *testing.T) {
	tests := []struct {
		name    string
		domains []string
		want    []string
	}{
		{
			name:    "filters domains with digits",
			domains: []string{"doppiocdn.com", "edge-hls2.doppiocdn.com", "strip-cdn.com"},
			want:    []string{"doppiocdn.com", "strip-cdn.com"},
		},
		{name: "all filtered", domains: []string{"a1.com", "b2.com"}, want: []string{}},
		{name: "empty", domains: nil, want: []string{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var id InitialDynamic
			for _, d := range tt.domains {
				id.Players.CDNConfig = append(id.Players.CDNConfig, cdn(d))
			}
			got := id.GetDomainsCDN()
			if !slices.Equal(got, tt.want) {
				t.Errorf("GetDomainsCDN() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetDomainCDN(t *testing.T) {
	t.Run("no domains", func(t *testing.T) {
		var id InitialDynamic
		got, err := id.GetDomainCDN()
		if err == nil {
			t.Fatalf("GetDomainCDN() error = nil, want error")
		}
		if got != "" {
			t.Errorf("GetDomainCDN() = %q, want empty", got)
		}
	})

	t.Run("single domain", func(t *testing.T) {
		var id InitialDynamic
		id.Players.CDNConfig = []CDNConfig{cdn("doppiocdn.com")}
		got, err := id.GetDomainCDN()
		if err != nil {
			t.Fatalf("GetDomainCDN() error = %v", err)
		}
		if got != "doppiocdn.com" {
			t.Errorf("GetDomainCDN() = %q, want %q", got, "doppiocdn.com")
		}
	})

	t.Run("random from list", func(t *testing.T) {
		var id InitialDynamic
		want := []string{"one.com", "two.com", "three.com"}
		for _, d := range want {
			id.Players.CDNConfig = append(id.Players.CDNConfig, cdn(d))
		}
		for range 50 {
			got, err := id.GetDomainCDN()
			if err != nil {
				t.Fatalf("GetDomainCDN() error = %v", err)
			}
			if !slices.Contains(want, got) {
				t.Fatalf("GetDomainCDN() = %q, want one of %v", got, want)
			}
		}
	})
}

func TestPlayersGetCDNConfigs(t *testing.T) {
	p := Players{CDNConfig: []CDNConfig{cdn("a.com"), cdn("b.com")}}
	if got := p.GetCDNConfigs(); len(got) != 2 || got[0].Domain != "a.com" {
		t.Errorf("GetCDNConfigs() = %v, want [a.com b.com]", got)
	}
}

// TestInitialDynamicResponseJSON проверяет разбор реального ответа
// /api/front/v3/config/initial-dynamic (только используемые поля).
func TestInitialDynamicResponseJSON(t *testing.T) {
	const body = `{
		"initialDynamic": {
			"csrfToken": "csrf-token",
			"csrfTimestamp": "2026-03-13T10:26:31.846Z",
			"guestId": 4242,
			"sessionHash": "hash",
			"tabId": "tab",
			"websocketUrl": "wss://example.com/ws",
			"websocket": {"url": "wss://example.com/ws", "token": "tok"},
			"players": {
				"playoutDelayHint": 0.5,
				"isHLSLLEnabled": true,
				"llCountriesBlacklist": ["RU"],
				"cdnConfig": [
					{"domain": "doppiocdn.com", "features": {"ll": true}},
					{"domain": "edge-hls.doppiocdn.com", "features": {"ll": true}},
					{"domain": "doppiocdn2.com", "features": {"ll": false}}
				],
				"broadcastResolutionList": {
					"720p": {"bitrate": 2500, "width": 1280, "height": 720, "isEnabled": true}
				}
			},
			"unknownField": "ignored"
		}
	}`

	var resp InitialDynamicResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	id := resp.InitialDynamic
	if id.CSRFToken != "csrf-token" || id.GuestID != 4242 {
		t.Errorf("InitialDynamic = %+v, want csrfToken/guestId", id)
	}
	if id.Websocket.Token != "tok" {
		t.Errorf("Websocket = %+v, want token=tok", id.Websocket)
	}
	if !id.Players.IsHLSLLEnabled || id.Players.PlayoutDelayHint != 0.5 {
		t.Errorf("Players = %+v, want isHLSLLEnabled/playoutDelayHint", id.Players)
	}
	if got := id.GetDomainsCDN(); !slices.Equal(
		got,
		[]string{"doppiocdn.com", "edge-hls.doppiocdn.com"},
	) {
		t.Errorf("GetDomainsCDN() = %v, want [doppiocdn.com edge-hls.doppiocdn.com]", got)
	}
	if res := id.Players.BroadcastResolutionList["720p"]; res.Width != 1280 || !res.IsEnabled {
		t.Errorf("BroadcastResolutionList = %+v, want width=1280", res)
	}
}
