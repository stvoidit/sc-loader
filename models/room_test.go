package models

import (
	"encoding/json/v2"
	"testing"
	"time"
)

func TestRoomResponseJSON(t *testing.T) {
	const body = `{
		"user": {
			"user": {
				"id": 11111111,
				"isBlocked": false,
				"isCam": true,
				"isLive": true,
				"isModel": true,
				"login": "alex",
				"name": "Alex",
				"status": "public",
				"username": "Alex"
			}
		}
	}`

	var rr RoomResponse
	if err := json.Unmarshal([]byte(body), &rr); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if got := rr.GetRoomId(); got != 11111111 {
		t.Errorf("GetRoomId() = %d, want %d", got, 11111111)
	}
	if got := rr.GetStatus(); got != "public" {
		t.Errorf("GetStatus() = %q, want %q", got, "public")
	}
	if !rr.IsOnline() {
		t.Error("IsOnline() = false, want true")
	}
}

func TestRoomResponseIsOnline(t *testing.T) {
	tests := []struct {
		name   string
		status string
		isLive bool
		want   bool
	}{
		{"public and live", "public", true, true},
		{"PUBLIC and live", "PUBLIC", true, true},
		{"public but not live", "public", false, false},
		{"private but live", "private", true, false},
		{"empty", "", false, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var rr RoomResponse
			rr.User.User.Status = tt.status
			rr.User.User.IsLive = tt.isLive
			if got := rr.IsOnline(); got != tt.want {
				t.Errorf("IsOnline() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBroadcastJSON(t *testing.T) {
	const body = `{
		"item": {
			"appType": "web",
			"isApprovedModel": true,
			"isBadStream": false,
			"isBlocked": false,
			"isDeleted": false,
			"isLive": true,
			"isMobile": false,
			"modelId": 55555,
			"settings": {
				"fps": 30,
				"height": 720,
				"width": 1280,
				"mediaTransport": "hls",
				"video": {"bitrate": 2500, "codec": "h264", "keyInterval": 2}
			},
			"snapshotTimestamp": 1773386791,
			"status": "public",
			"statusChangedAt": "2026-03-13T10:26:31.846+03:00",
			"streamName": "stream_55555",
			"username": "alex",
			"webRTCAppKey": "key"
		}
	}`

	var b Broadcast
	if err := json.Unmarshal([]byte(body), &b); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if !b.IsOnline() {
		t.Error("IsOnline() = false, want true")
	}
	if got := b.GetStatus(); got != "public" {
		t.Errorf("GetStatus() = %q, want %q", got, "public")
	}
	// проверка, что вложенные поля и время действительно распарсились
	if b.Item.ModelID != 55555 || b.Item.Username != "alex" {
		t.Errorf("Item = %+v, want modelId=55555 username=alex", b.Item)
	}
	if b.Item.Settings.Video.Bitrate != 2500 || b.Item.Settings.Fps != 30 {
		t.Errorf("Settings = %+v, want bitrate=2500 fps=30", b.Item.Settings)
	}
	want := time.Date(2026, 3, 13, 10, 26, 31, 846000000, time.FixedZone("", 3*3600))
	if !b.Item.StatusChangedAt.Equal(want) {
		t.Errorf("StatusChangedAt = %v, want %v", b.Item.StatusChangedAt, want)
	}
}

func TestBroadcastOffline(t *testing.T) {
	var b Broadcast
	if err := json.Unmarshal(
		[]byte(`{"item":{"isLive":false,"status":"offline"}}`),
		&b,
	); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if b.IsOnline() {
		t.Error("IsOnline() = true, want false")
	}
	if got := b.GetStatus(); got != "offline" {
		t.Errorf("GetStatus() = %q, want %q", got, "offline")
	}
}

func TestBroadcastEmptyJSON(t *testing.T) {
	// API может вернуть {} или тело другой структуры - геттеры не должны паниковать
	var b Broadcast
	if err := json.Unmarshal([]byte(`{}`), &b); err != nil {
		t.Fatalf("Unmarshal() error = %v", err)
	}
	if b.IsOnline() || b.GetStatus() != "" {
		t.Errorf("Broadcast{} = (%v, %q), want (false, \"\")", b.IsOnline(), b.GetStatus())
	}
}
