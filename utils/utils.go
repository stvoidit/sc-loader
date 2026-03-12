package utils

import (
	"bytes"
	"encoding/json/v2"
	"fmt"
	"io"
	"maps"
	"math/rand/v2"
	"os"
	"slices"
)

func SaveFile(filename string, rc io.ReadCloser) io.ReadCloser {
	f, err := os.Create(filename)
	if err != nil {
		panic(err)
	}
	defer DeferCloseReader(f)
	defer DeferCloseReader(rc)
	b, err := io.ReadAll(rc)
	if err != nil {
		panic(err)
	}
	if _, err := f.Write(b); err != nil {
		panic(err)
	}
	return io.NopCloser(bytes.NewReader(b))
}

func MakeUniq() string {
	var pack = []rune("0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz")
	var buf = make([]rune, 16)
	for i := range buf {
		buf[i] = pack[rand.IntN(len(pack))]
	}
	return string(buf)
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

func FormatFileSize[T ~int | int64](size T) string {
	switch {
	case size < KB:
		return fmt.Sprintf("%d B", size)
	case size < MB:
		return fmt.Sprintf("%.1f KB", float64(size)/KB)
	case size < GB:
		return fmt.Sprintf("%.1f MB", float64(size)/MB)
	case size < TB:
		return fmt.Sprintf("%.1f GB", float64(size)/GB)
	default:
		return fmt.Sprintf("%.1f TB", float64(size)/TB)
	}
}

type Probe struct{}

var probe = Probe{}

type Set[T comparable] map[T]Probe

func (s Set[T]) Has(value T) (ok bool) {
	_, ok = s[value]
	return
}
func (s Set[T]) Set(value T) (ok bool) {
	s[value] = probe
	return
}

func (s Set[T]) String() string {
	var keys = slices.Collect(maps.Keys(s))
	out, _ := json.Marshal(keys)
	return string(out)
}
