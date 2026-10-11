package utils

import (
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"net/url"
	"strings"
	"testing"
)

// Реальная пара ключ/шифртекст DRM mouflon (проверена в tests/main.js).
const (
	drmKey       = "EQueeGh2kaewa3ch"
	drmEncrypted = "gqaq55h8zoL+Lnlg9nmPX7"
	drmDecrypted = "CblB7Ti4PVvMt24g"
)

func testPSCH(key string) PSCH { return sha256.Sum256([]byte(key)) }

// encodePart - операция, обратная Decode: нужна для синтетических проверок.
func encodePart(plain string, psch PSCH) string {
	src := []byte(plain)
	xored := make([]byte, len(src))
	for i := range src {
		xored[i] = src[i] ^ psch[i%len(psch)]
	}
	return reverseString(base64.RawStdEncoding.EncodeToString(xored))
}

func TestDecode(t *testing.T) {
	psch := testPSCH(drmKey)
	tests := []struct {
		name    string
		in      string
		want    string
		wantErr bool
	}{
		{name: "real vector", in: drmEncrypted, want: drmDecrypted},
		{name: "empty", in: "", want: ""},
		{name: "invalid base64", in: "!!!not base64!!!", want: "!!!not base64!!!", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Decode(tt.in, psch)
			if (err != nil) != tt.wantErr {
				t.Fatalf("Decode() error = %v, wantErr %v", err, tt.wantErr)
			}
			if got != tt.want {
				t.Errorf("Decode() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeRoundTrip(t *testing.T) {
	psch := testPSCH(drmKey)
	plainList := []string{
		"a",
		drmDecrypted,
		strings.Repeat("long-part-name-", 5), // длиннее psch, проверяет i%len(psch)
		"init_12345.m4s",
	}
	for _, plain := range plainList {
		t.Run(plain[:min(len(plain), 12)], func(t *testing.T) {
			got, err := Decode(encodePart(plain, psch), psch)
			if err != nil {
				t.Fatalf("Decode() error = %v", err)
			}
			if got != plain {
				t.Errorf("Decode() = %q, want %q", got, plain)
			}
		})
	}
}

func TestDecodeWrongKey(t *testing.T) {
	got, err := Decode(drmEncrypted, testPSCH("another-key"))
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}
	if got == drmDecrypted {
		t.Errorf("Decode() с чужим ключом вернул верный результат %q", got)
	}
}

func TestReverseString(t *testing.T) {
	tests := []struct{ in, want string }{
		{"", ""},
		{"a", "a"},
		{"abc", "cba"},
		{"aba", "aba"},
		{"привет", "тевирп"},
	}
	for _, tt := range tests {
		if got := reverseString(tt.in); got != tt.want {
			t.Errorf("reverseString(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestParseRendition(t *testing.T) {
	tests := []struct {
		name     string
		in       string
		wantMSN  string
		wantPart string
		wantLen  int
	}{
		{
			name:     "both params",
			in:       `"../720p/live.m3u8",LAST-MSN=1234,LAST-PART=2`,
			wantMSN:  "1234",
			wantPart: "2",
			wantLen:  2,
		},
		{
			name:    "only msn",
			in:      `LAST-MSN=99`,
			wantMSN: "99",
			wantLen: 1,
		},
		{
			name:    "no params",
			in:      `"../720p/live.m3u8"`,
			wantLen: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := parseRendition(tt.in)
			if len(got) != tt.wantLen {
				t.Errorf("parseRendition() len = %d, want %d (%v)", len(got), tt.wantLen, got)
			}
			if got.Get("_HLS_msn") != tt.wantMSN {
				t.Errorf("_HLS_msn = %q, want %q", got.Get("_HLS_msn"), tt.wantMSN)
			}
			if got.Get("_HLS_part") != tt.wantPart {
				t.Errorf("_HLS_part = %q, want %q", got.Get("_HLS_part"), tt.wantPart)
			}
		})
	}
}

func TestParseRaw(t *testing.T) {
	const variant = "https://edge-hls.example.com/hls/123/master/123_auto.m3u8"
	tests := []struct {
		name        string
		body        string
		wantVariant string
		wantPKey    string
	}{
		{
			name: "full playlist",
			body: strings.Join([]string{
				"#EXTM3U",
				"#EXT-X-MOUFLON:PSCH:v2:" + drmKey,
				"#EXT-X-STREAM-INF:BANDWIDTH=1000",
				variant,
			}, "\n"),
			wantVariant: variant,
			wantPKey:    drmKey,
		},
		{
			name:        "no pkey",
			body:        "#EXTM3U\n" + variant,
			wantVariant: variant,
		},
		{
			name:     "no variant",
			body:     "#EXTM3U\n#EXT-X-MOUFLON:PSCH:v2:" + drmKey,
			wantPKey: drmKey,
		},
		{
			name: "first values win",
			body: strings.Join([]string{
				"#EXT-X-MOUFLON:PSCH:v2:first",
				"#EXT-X-MOUFLON:PSCH:v2:second",
				"https://example.com/one.m3u8",
				"https://example.com/two.m3u8",
			}, "\n"),
			wantVariant: "https://example.com/one.m3u8",
			wantPKey:    "first",
		},
		{name: "empty", body: ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			variant, pkey := ParseRaw(strings.NewReader(tt.body))
			if variant != tt.wantVariant || pkey != tt.wantPKey {
				t.Errorf("ParseRaw() = (%q, %q), want (%q, %q)",
					variant, pkey, tt.wantVariant, tt.wantPKey)
			}
		})
	}
}

// partURL собирает ссылку на part-сегмент, где зашифрованное имя стоит
// третьим с конца элементом пути (ровно так его ищет DecodePart).
func partURL(encrypted string) string {
	return "https://edge-hls.example.com/hls/123/master/123_auto/init_a1b2_part_" +
		encrypted + "_7_9.m4s"
}

func TestDecodePart(t *testing.T) {
	psch := testPSCH(drmKey)
	tests := []struct {
		name     string
		link     string
		want     string
		wantErr  error
		wantAnyE bool
	}{
		{
			name: "part link decoded",
			link: partURL(drmEncrypted),
			want: partURL(drmDecrypted),
		},
		{
			name:    "not a part link",
			link:    "https://edge-hls.example.com/hls/123/master/123_auto.m3u8",
			want:    "https://edge-hls.example.com/hls/123/master/123_auto.m3u8",
			wantErr: ErrCanSkip,
		},
		{
			name:    "invalid url",
			link:    "https://example.com/%zz",
			wantErr: url.EscapeError("%zz"),
		},
		{
			name:     "invalid base64 in part",
			link:     partURL("!!!"),
			want:     partURL("!!!"),
			wantAnyE: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := DecodePart(tt.link, psch)
			switch {
			case tt.wantAnyE:
				if err == nil {
					t.Fatal("DecodePart() error = nil, want error")
				}
			case tt.wantErr != nil:
				if !errors.Is(err, tt.wantErr) {
					t.Fatalf("DecodePart() error = %v, want %v", err, tt.wantErr)
				}
			case err != nil:
				t.Fatalf("DecodePart() unexpected error = %v", err)
			}
			if tt.want != "" && got != tt.want {
				t.Errorf("DecodePart() = %q, want %q", got, tt.want)
			}
		})
	}
}

func TestDecodeRaw(t *testing.T) {
	psch := testPSCH(drmKey)
	t.Run("full playlist", func(t *testing.T) {
		body := strings.Join([]string{
			"#EXTM3U",
			"#EXT-X-TARGETDURATION:2",
			`#EXT-X-MAP:URI="init.mp4"`,
			"#EXT-X-MOUFLON:PSCH:v2:" + drmKey,
			`#EXT-X-RENDITION-REPORT:URI="../720p/live.m3u8",LAST-MSN=1234,LAST-PART=2`,
			"#EXTINF:2.000,",
			"#EXT-X-MOUFLON:URI:" + partURL(drmEncrypted),
			"#EXT-X-MOUFLON:URI:https://edge-hls.example.com/hls/123/skip-me.m4s",
		}, "\n")

		vids, hlsQuery, err := DecodeRaw(strings.NewReader(body), psch)
		if err != nil {
			t.Fatalf("DecodeRaw() error = %v", err)
		}
		want := []string{"init.mp4", partURL(drmDecrypted)}
		if len(vids) != len(want) {
			t.Fatalf("DecodeRaw() vids = %v, want %v", vids, want)
		}
		for i := range want {
			if vids[i] != want[i] {
				t.Errorf("vids[%d] = %q, want %q", i, vids[i], want[i])
			}
		}
		if hlsQuery.Get("_HLS_msn") != "1234" || hlsQuery.Get("_HLS_part") != "2" {
			t.Errorf("hlsQuery = %v, want _HLS_msn=1234&_HLS_part=2", hlsQuery)
		}
	})

	t.Run("empty body", func(t *testing.T) {
		vids, hlsQuery, err := DecodeRaw(strings.NewReader(""), psch)
		if err != nil {
			t.Fatalf("DecodeRaw() error = %v", err)
		}
		if vids == nil || len(vids) != 0 {
			t.Errorf("vids = %v, want empty slice", vids)
		}
		if hlsQuery != nil {
			t.Errorf("hlsQuery = %v, want nil", hlsQuery)
		}
	})

	t.Run("decode error", func(t *testing.T) {
		body := "#EXT-X-MOUFLON:URI:" + partURL("!!!")
		if _, _, err := DecodeRaw(strings.NewReader(body), psch); err == nil {
			t.Fatal("DecodeRaw() error = nil, want base64 error")
		}
	})

	t.Run("map uri without quotes", func(t *testing.T) {
		vids, _, err := DecodeRaw(strings.NewReader("#EXT-X-MAP:URI=init.mp4"), psch)
		if err != nil {
			t.Fatalf("DecodeRaw() error = %v", err)
		}
		if len(vids) != 1 || vids[0] != "init.mp4" {
			t.Errorf("vids = %v, want [init.mp4]", vids)
		}
	})
}
