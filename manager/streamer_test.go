package manager_test

import (
	"log/slog"
	"regexp"
	"testing"

	"sc-loader/manager"
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

func TestParseUsernameEdgeCases(t *testing.T) {
	const host = "examle.com"
	tests := []struct {
		name     string
		arg      string
		wantUser string
		wantURL  string
		wantErr  bool
	}{
		{
			name:     "trailing slash",
			arg:      "https://examle.com/superUser/",
			wantUser: "superUser",
			wantURL:  "https://examle.com/superUser",
		},
		{
			name:     "query is kept in url",
			arg:      "https://examle.com/superUser?utm=1",
			wantUser: "superUser",
			wantURL:  "https://examle.com/superUser?utm=1",
		},
		{
			name:     "empty username",
			arg:      "",
			wantUser: "",
			wantURL:  "https://examle.com",
		},
		{
			// текущее поведение: filepath.Base("") возвращает "."
			name:     "url without path",
			arg:      "https://examle.com",
			wantUser: ".",
			wantURL:  "https://examle.com",
		},
		{
			name:    "invalid url",
			arg:     "https://examle.com/%zz",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := manager.ParseUsername(tt.arg, host)
			if tt.wantErr {
				if err == nil {
					t.Fatalf("ParseUsername(%q) error = nil, want error", tt.arg)
				}
				return
			}
			if err != nil {
				t.Fatalf("ParseUsername(%q) error = %v", tt.arg, err)
			}
			if got.Username != tt.wantUser || got.SreamURL != tt.wantURL {
				t.Errorf("ParseUsername(%q) = %+v, want {%s %s}", tt.arg, got, tt.wantUser, tt.wantURL)
			}
		})
	}
}

func TestMakeFilename(t *testing.T) {
	s := manager.Streamer{Username: "alex", SreamURL: "https://examle.com/alex"}
	re := regexp.MustCompile(`^alex_\d{8}-\d{6}_tmp\.mp4$`)
	if got := s.MakeFilename(); !re.MatchString(got) {
		t.Errorf("MakeFilename() = %q, want alex_YYYYMMDD-HHMMSS_tmp.mp4", got)
	}
}

func TestStreamerLogValue(t *testing.T) {
	s := manager.Streamer{Username: "alex", SreamURL: "https://examle.com/alex"}
	v := s.LogValue()
	if v.Kind() != slog.KindGroup {
		t.Fatalf("LogValue().Kind() = %v, want %v", v.Kind(), slog.KindGroup)
	}
	attrs := map[string]string{}
	for _, a := range v.Group() {
		attrs[a.Key] = a.Value.String()
	}
	if attrs["username"] != "alex" || attrs["url"] != "https://examle.com/alex" {
		t.Errorf("LogValue() = %v", attrs)
	}
}
