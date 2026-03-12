package utils

import (
	"context"
	"encoding/json/jsontext"
	"encoding/json/v2"
	"errors"
	"io"
	"log/slog"
	"os"
	"strings"
)

func FileStat(f *os.File) (string, error) {
	fi, err := f.Stat()
	if err != nil {
		return "", err
	}
	size := fi.Size()
	return FormatFileSize(size), nil
}

func IsCancel(err error) bool {
	const interroptErr = "interrupt signal received"
	return errors.Is(err, context.Canceled) || strings.Contains(err.Error(), interroptErr)
}

func DeferCloseReader(rc io.ReadCloser) {
	if rc == nil {
		return
	}
	if err := rc.Close(); err != nil {
		slog.Error("DeferCloseReader", slog.String("error", err.Error()))
	}
}

func IsEmpty(s string) bool {
	return strings.EqualFold(s, "")
}

func PPrint(w io.Writer, i any) {
	if err := json.MarshalWrite(w, i, jsontext.WithIndent("  ")); err != nil {
		slog.Error("PPrint", slog.String("error", err.Error()))
	}
}
