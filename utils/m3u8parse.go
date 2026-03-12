package utils

import (
	"bufio"
	"encoding/base64"
	"errors"
	"io"
	"net/url"
	"strings"
)

// const _mouflon_filename = "media.mp4"
const _mouflon_file_attr = "#EXT-X-MOUFLON:URI:"
const _map_uri = "#EXT-X-MAP:URI="
const _mouflon_psch = "#EXT-X-MOUFLON:PSCH:v2:"
const INDEPENDENT = "INDEPENDENT"
const _rendition = "#EXT-X-RENDITION-REPORT:URI="

func parseRendition(s string) (hlsQuery url.Values) {
	var params = [...][2]string{
		[...]string{"LAST-MSN=", "_HLS_msn"},
		[...]string{"LAST-PART=", "_HLS_part"},
	}
	hlsQuery = make(url.Values, len(params))
	var parts = strings.Split(s, ",")
	for _, param := range params {
		for _, part := range parts {
			if after, ok := strings.CutPrefix(part, param[0]); ok {
				hlsQuery.Set(param[1], after)
			}
		}
	}
	return
}

func Decode(encryptedPart string, psch [32]byte) (string, error) {
	decodedB64, err := base64.RawStdEncoding.DecodeString(reverseString(encryptedPart))
	if err != nil {
		return encryptedPart, err
	}
	var out = make([]byte, len(decodedB64))
	for i := range decodedB64 {
		out[i] = decodedB64[i] ^ psch[i%len(psch)]
	}
	return string(out), nil
}

func reverseString(s string) string {
	r := []rune(s)
	for i, j := 0, len(r)-1; i < j; i, j = i+1, j-1 {
		r[i], r[j] = r[j], r[i]
	}
	return string(r)
}

func ParseRaw(r io.Reader) (variant, pkey string) {
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if afterCut, ok := strings.CutPrefix(line, _mouflon_psch); IsEmpty(pkey) && ok {
			pkey = afterCut
		}
		if IsEmpty(variant) && strings.HasPrefix(line, "https") {
			variant = line
		}
	}
	// slog.Debug("ParseRaw", slog.String("variant", variant), slog.String("pkey", pkey))
	return variant, pkey
}

var ErrCanSkip = errors.New("can skip this part")

func DecodePart(link string, psch [32]byte) (string, error) {
	URL, err := url.Parse(link)
	if err != nil {
		return "", err
	}
	// slog.Debug("decodePart", slog.String("URL.Path", URL.Path))
	parts := strings.Split(URL.Path, "_")
	var encryptedPart string
	if strings.EqualFold(encryptedPart, "init") {
		encryptedPart = parts[len(parts)-1]
	} else if strings.Contains(link, "_part") {
		encryptedPart = parts[len(parts)-3]
	} else {
		// encryptedPart = parts[len(parts)-2]
		return link, ErrCanSkip
	}
	// slog.Debug("decodePart", "encryptedPart", encryptedPart)
	encryptedB64, err := Decode(encryptedPart, psch)
	if err != nil {
		return link, err
	}
	// slog.Debug("decodePart", "encryptedB64", encryptedB64)
	URL.Path = strings.Replace(URL.Path, encryptedPart, encryptedB64, 1)
	return URL.String(), nil
}

func DecodeRaw(r io.Reader, psch PSCH) (vids []string, hlsQuery url.Values, err error) {
	vids = make([]string, 0, 4)
	s := bufio.NewScanner(r)
	for s.Scan() {
		line := s.Text()
		if after, ok := strings.CutPrefix(line, _map_uri); ok {
			vids = append(vids, strings.Trim(after, `"`))
		} else if after, ok := strings.CutPrefix(line, _rendition); ok {
			hlsQuery = parseRendition(after)
		} else if after, ok := strings.CutPrefix(line, _mouflon_file_attr); ok {
			playlist, err := DecodePart(after, psch)
			if err != nil {
				if errors.Is(err, ErrCanSkip) {
					continue
				}
				return vids, hlsQuery, err
			}
			vids = append(vids, playlist)
		}
	}
	return vids, hlsQuery, err
}
