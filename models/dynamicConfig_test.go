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

// TestInitialDynamicResponseJSON проверяет разбор ответа
// /api/front/v3/config/initial-dynamic: неизвестные поля игнорируются.
func TestInitialDynamicResponseJSON(t *testing.T) {
	const body = `{
		"initialDynamic": {
			"csrfToken": "csrf-token",
			"guestId": 4242,
			"players": {
				"playoutDelayHint": 0.5,
				"isHLSLLEnabled": true,
				"cdnConfig": [
					{"domain": "doppiocdn.com", "features": {"ll": true}},
					{"domain": "edge-hls.doppiocdn.com", "features": {"ll": true}},
					{"domain": "doppiocdn2.com", "features": {"ll": false}}
				]
			},
			"unknownField": "ignored"
		}
	}`

	var resp InitialDynamicResponse
	if err := json.Unmarshal([]byte(body), &resp); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	id := resp.InitialDynamic
	if len(id.Players.CDNConfig) != 3 {
		t.Fatalf("CDNConfig = %+v, want 3 элемента", id.Players.CDNConfig)
	}
	if got := id.GetDomainsCDN(); !slices.Equal(
		got,
		[]string{"doppiocdn.com", "edge-hls.doppiocdn.com"},
	) {
		t.Errorf("GetDomainsCDN() = %v, want [doppiocdn.com edge-hls.doppiocdn.com]", got)
	}
}
