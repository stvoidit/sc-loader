package models

import (
	"math/rand/v2"
	"regexp"
	"time"
)

type InitialDynamicResponse struct {
	InitialDynamic InitialDynamic `json:"initialDynamic"`
}

type InitialDynamic struct {
	CSRFToken           string    `json:"csrfToken"`
	CSRFTimestamp       time.Time `json:"csrfTimestamp"`
	CSRFNotifyTimestamp time.Time `json:"csrfNotifyTimestamp"`
	JwtToken            *string   `json:"jwtToken"`
	GuestID             int64     `json:"guestId"`
	SessionHash         string    `json:"sessionHash"`
	UserHash            string    `json:"userHash"`
	IPHash              int64     `json:"ipHash"`
	CometAuth           string    `json:"cometAuth"`
	TabID               string    `json:"tabId"`
	WebsocketURL        string    `json:"websocketUrl"`
	Websocket           Websocket `json:"websocket"`
	Players             Players   `json:"players"`
	// WebsocketApps        Websocket `json:"websocketApps"`
	// CAPTCHA              CAPTCHA   `json:"captcha"`
	// Turnstile            Turnstile `json:"turnstile"`
	// Flags                []string  `json:"flags"`
	// User                 any       `json:"user"`
	// IsPublic             bool      `json:"isPublic"`
	// PreferredTag         string    `json:"preferredTag"`
	// DetectedPreferredTag string    `json:"detectedPreferredTag"`

}

type CAPTCHA struct {
	IsEnabled  bool      `json:"isEnabled"`
	IsRequired bool      `json:"isRequired"`
	Type       string    `json:"type"`
	Config     Turnstile `json:"config"`
}

type Turnstile struct {
	SiteKey string `json:"siteKey"`
}

type Players struct {
	// IsVP8BroadcastAvailable                  bool                               `json:"isVP8BroadcastAvailable"`
	// IsNewPresetsEnabled                      bool                               `json:"isNewPresetsEnabled"`
	// AutoResolutionPlayers                    []string                           `json:"autoResolutionPlayers"`
	// WebRTCMeasurementLockTimeout             int64                              `json:"webRTCMeasurementLockTimeout"`
	// WebrtcEdgeWSURL                          string                             `json:"webrtcEdgeWSUrl"`
	// WebrtcOriginWSURL                        string                             `json:"webrtcOriginWSUrl"`
	// WebrtcPixelateURL                        string                             `json:"webrtcPixelateUrl"`
	PlayoutDelayHint                         float64                            `json:"playoutDelayHint"`
	HLSLowLatencyPlayerTypeDisabledPlatforms []string                           `json:"hlsLowLatencyPlayerTypeDisabledPlatforms"`
	HLSLowLatencyDisabledBrowsersMap         HLSLowLatencyDisabledBrowsersMap   `json:"hlsLowLatencyDisabledBrowsersMap"`
	WebRTCMeasurementFactors                 WebRTCMeasurementFactors           `json:"webRTCMeasurementFactors"`
	WebRTCMeasurementConstants               WebRTCMeasurementConstants         `json:"webRTCMeasurementConstants"`
	CDNConfig                                []CDNConfig                        `json:"cdnConfig"`
	LlCountriesBlacklist                     []string                           `json:"llCountriesBlacklist"`
	BroadcastResolutionList                  map[string]BroadcastResolutionList `json:"broadcastResolutionList"`
	IsHLSLLEnabled                           bool                               `json:"isHLSLLEnabled"`
}

func (p Players) GetCDNConfigs() []CDNConfig {
	return p.CDNConfig
}

type BroadcastResolutionList struct {
	Bitrate      int64         `json:"bitrate"`
	Width        int64         `json:"width"`
	Height       int64         `json:"height"`
	IsEnabled    bool          `json:"isEnabled"`
	BitrateByFPS *BitrateByFPS `json:"bitrateByFPS,omitempty"`
}

type BitrateByFPS struct {
	The50   int64 `json:"50"`
	Default int64 `json:"default"`
}

type CDNConfig struct {
	Domain   string `json:"domain"`
	Features struct {
		Ll bool `json:"ll"`
	} `json:"features"`
}

type HLSLowLatencyDisabledBrowsersMap struct {
	Firefox string `json:"firefox"`
}

type WebRTCMeasurementConstants struct {
	ABRAverageUpFactor   float64 `json:"abrAverageUpFactor"`
	ABRAverageDownFactor float64 `json:"abrAverageDownFactor"`
	ABRVectorUpFactor    float64 `json:"abrVectorUpFactor"`
	ABRVectorDownFactor  int64   `json:"abrVectorDownFactor"`
}

type WebRTCMeasurementFactors struct {
	Rtt         Bitrate `json:"rtt"`
	PacketsLost Bitrate `json:"packetsLost"`
	Bitrate     Bitrate `json:"bitrate"`
}

type Bitrate struct {
	MeasurementFault float64 `json:"measurementFault"`
	Multiplier       float64 `json:"multiplier"`
	Ratio            float64 `json:"ratio"`
}

type Websocket struct {
	URL   string `json:"url"`
	Token string `json:"token"`
}

func (id InitialDynamic) GetDomainsCDN() []string {
	var rgx = regexp.MustCompile(`\d+`)
	var cleaned = make([]string, 0, len(id.Players.CDNConfig))
	for i := range id.Players.CDNConfig {
		if !rgx.MatchString(id.Players.CDNConfig[i].Domain) {
			cleaned = append(cleaned, id.Players.CDNConfig[i].Domain)
		}
	}
	return cleaned
}

func (id InitialDynamic) GetDomainCDN() (domain string) {
	var cleaned = id.GetDomainsCDN()
	return cleaned[rand.IntN(len(cleaned))]

}
