package models

import (
	"strings"
	"time"
)

type RoomResponse struct {
	User struct {
		User struct {
			ID        int    `json:"id"`
			IsBlocked bool   `json:"isBlocked"`
			IsCam     bool   `json:"isCam"`
			IsLive    bool   `json:"isLive"`
			IsModel   bool   `json:"isModel"`
			Login     string `json:"login"`
			Name      string `json:"name"`
			Status    string `json:"status"`
			Username  string `json:"username"`
		} `json:"user"`
	} `json:"user"`
}

func (rr RoomResponse) GetRoomId() int {
	return rr.User.User.ID
}

func (rr RoomResponse) IsOnline() bool {
	return strings.EqualFold(rr.User.User.Status, "public") && rr.User.User.IsLive
}

func (rr RoomResponse) GetStatus() string {
	return rr.User.User.Status
}

type BroadcastStatus string

func (bs BroadcastStatus) String() string {
	return string(bs)
}

const (
	StatusPublic    = "public"
	StatusPrivate   = "private"
	StatusGroupShow = "groupShow"
	StatusOff       = "off"
)

type Broadcast struct {
	Item struct {
		AppType         string `json:"appType"`
		IsApprovedModel bool   `json:"isApprovedModel"`
		IsBadStream     bool   `json:"isBadStream"`
		IsBlocked       bool   `json:"isBlocked"`
		IsDeleted       bool   `json:"isDeleted"`
		IsLive          bool   `json:"isLive"`
		IsMicInactive   bool   `json:"isMicInactive"`
		IsMobile        bool   `json:"isMobile"`
		IsNonNude       bool   `json:"isNonNude"`
		ModelID         int    `json:"modelId"`
		Settings        struct {
			AudioQuality struct {
				Network string `json:"network"`
				Stream  string `json:"stream"`
			} `json:"audioQuality"`
			Fps            int    `json:"fps"`
			Height         int    `json:"height"`
			IsBadStream    bool   `json:"isBadStream"`
			IsMicInactive  bool   `json:"isMicInactive"`
			IsSourcePreset bool   `json:"isSourcePreset"`
			MediaTransport string `json:"mediaTransport"`
			Presets        []any  `json:"presets"`
			Video          struct {
				BFrames     bool   `json:"bFrames"`
				Bitrate     int    `json:"bitrate"`
				Codec       string `json:"codec"`
				KeyInterval int    `json:"keyInterval"`
			} `json:"video"`
			VideoQuality struct {
				Network string `json:"network"`
				Stream  string `json:"stream"`
			} `json:"videoQuality"`
			Width int `json:"width"`
		} `json:"settings"`
		SnapshotTimestamp int             `json:"snapshotTimestamp"`
		Status            BroadcastStatus `json:"status"`
		StatusChangedAt   time.Time       `json:"statusChangedAt"`
		StreamName        string          `json:"streamName"`
		Username          string          `json:"username"`
		WebRTCAppKey      string          `json:"webRTCAppKey"`
	} `json:"item"`
}

func (b Broadcast) IsOnline() bool {
	return b.Item.Status == StatusPublic
}

func (b Broadcast) GetStatus() BroadcastStatus {
	return b.Item.Status
}
