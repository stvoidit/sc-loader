package utils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
)

type ReplaceOption func(a *slog.Attr)

// func uppercaseLevel(a *slog.Attr) {
// 	if a.Key != slog.LevelKey {
// 		return
// 	}
// 	if level, ok := a.Value.Any().(slog.Level); ok {
// 		a.Value = slog.StringValue(level.String())
// 	}
// }

func shortSourceFile(a *slog.Attr) {
	fmt.Printf("%+v\n", a)
	if a.Key != slog.SourceKey {
		return
	}
	src, ok := a.Value.Any().(*slog.Source)
	if !ok || src == nil {
		return
	}
	src.File = filepath.Base(src.File)
	src.Function = src.Function
}

func NewSlogHandlerOptions(logLevel slog.Level, opts ...ReplaceOption) slog.HandlerOptions {
	replaceFn := func(_ []string, a slog.Attr) slog.Attr {
		for _, opt := range opts {
			opt(&a)
		}
		return a
	}
	return slog.HandlerOptions{
		AddSource:   true,
		Level:       logLevel,
		ReplaceAttr: replaceFn,
	}
}

func NewSloglogger(level slog.Level, outs ...io.Writer) *slog.Logger {
	var w io.Writer
	if lenOutes := len(outs); lenOutes == 1 {
		w = outs[0]
	} else if lenOutes > 1 {
		w = io.MultiWriter(outs...)
	} else {
		w = os.Stdout
	}
	options := NewSlogHandlerOptions(level, shortSourceFile)
	handler := slog.NewTextHandler(w, &options)
	logger := slog.New(handler)
	return logger
}

func SlogParseEnvLevel(varName string) (logLevel slog.Level) {
	str, ok := os.LookupEnv(varName)
	if !ok {
		return slog.LevelInfo
	}
	switch strings.ToLower(str) {
	case "debug", "dbg", "d", "-4":
		logLevel = slog.LevelDebug
	case "info", "i", "0":
		logLevel = slog.LevelInfo
	case "warn", "warning", "w", "4":
		logLevel = slog.LevelWarn
	case "err", "error", "e", "8":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}
	return logLevel
}
