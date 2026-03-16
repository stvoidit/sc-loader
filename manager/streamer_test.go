package manager_test

import (
	"sc-loader/manager"
	"testing"
)

func TestParseUsername(t *testing.T) {
	host := "examle.com"
	tests := []struct {
		name string
		arg  string
		host string
		want manager.Streamer
	}{
		{
			name: "TestParseUsername.username",
			arg:  "superUser1",
			host: host,
			want: manager.Streamer{
				Username: "superUser1",
				SreamURL: "https://examle.com/superUser1",
			},
		},
		{
			name: "TestParseUsername.url",
			arg:  "https://examle.com/superUser2",
			host: host,
			want: manager.Streamer{
				Username: "superUser2",
				SreamURL: "https://examle.com/superUser2",
			},
		},
		{
			name: "TestParseUsername.url2",
			arg:  "https://examle.com/superUser3/video",
			host: host,
			want: manager.Streamer{
				Username: "superUser3",
				SreamURL: "https://examle.com/superUser3",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, gotErr := manager.ParseUsername(tt.arg, tt.host)
			if gotErr != nil {
				t.Errorf("ParseUsername() failed: %v", gotErr)
				return
			}
			if got.Username != tt.want.Username || got.SreamURL != tt.want.SreamURL {
				t.Errorf("ParseUsername() = %+v, want %+v", got, tt.want)
			}
		})
	}
}
