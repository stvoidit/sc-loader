package models

import (
	"strings"
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
