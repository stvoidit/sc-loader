package manager

import (
	"log/slog"
	"net/url"
	"path/filepath"
	"strings"
)

type Streamer struct {
	Username string
	SreamURL string
}

func (s Streamer) LogValue() slog.Value {
	return slog.GroupValue(
		slog.String("username", s.Username),
		slog.String("url", s.SreamURL))
}

func ParseUsername(arg, host string) (Streamer, error) {
	var s Streamer
	if strings.HasPrefix(arg, "https") {
		URL, err := url.Parse(arg)
		if err != nil {
			return s, err
		}
		parts := strings.Split(URL.Path, "/")
		if len(parts) > 2 {
			URL.Path = "/" + parts[1]
		}
		s.SreamURL = URL.String()
		s.Username = filepath.Base(URL.Path)
	} else {
		s.Username = arg
		s.SreamURL = (&url.URL{Scheme: "https", Host: host, Path: arg}).String()
	}
	return s, nil
}
