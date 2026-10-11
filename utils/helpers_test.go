package utils

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestFileStat(t *testing.T) {
	path := filepath.Join(t.TempDir(), "file.bin")
	if err := os.WriteFile(path, bytes.Repeat([]byte{1}, 2048), 0o600); err != nil {
		t.Fatal(err)
	}

	f, err := os.Open(path)
	if err != nil {
		t.Fatal(err)
	}
	size, err := FileStat(f)
	if err != nil {
		t.Fatalf("FileStat() error = %v", err)
	}
	if size != "2.0 KB" {
		t.Errorf("FileStat() = %q, want %q", size, "2.0 KB")
	}

	if err := f.Close(); err != nil {
		t.Fatal(err)
	}
	if _, err := FileStat(f); err == nil {
		t.Error("FileStat() для закрытого файла error = nil, want error")
	}
}

func TestIsCancel(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{"canceled", context.Canceled, true},
		{"wrapped canceled", fmt.Errorf("download: %w", context.Canceled), true},
		{"interrupt signal", errors.New("interrupt signal received"), true},
		{
			"wrapped interrupt",
			fmt.Errorf("main: %w", errors.New("interrupt signal received")),
			true,
		},
		{"deadline", context.DeadlineExceeded, false},
		{"other", errors.New("connection refused"), false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := IsCancel(tt.err); got != tt.want {
				t.Errorf("IsCancel(%v) = %v, want %v", tt.err, got, tt.want)
			}
		})
	}
}

type trackingCloser struct {
	closed bool
	err    error
}

func (c *trackingCloser) Read(p []byte) (int, error) { return 0, io.EOF }

func (c *trackingCloser) Close() error {
	c.closed = true
	return c.err
}

func TestDeferClose(t *testing.T) {
	DeferClose(nil) // не должен паниковать

	ok := &trackingCloser{}
	DeferClose(ok)
	if !ok.closed {
		t.Error("DeferClose() не вызвал Close()")
	}

	// ошибка Close() логируется и не пробрасывается дальше
	bad := &trackingCloser{err: errors.New("close failed")}
	DeferClose(bad)
	if !bad.closed {
		t.Error("DeferClose() не вызвал Close() для ошибочного closer")
	}
}

func TestDeferCloseReader(t *testing.T) {
	DeferCloseReader(nil) // не должен паниковать

	rc := &trackingCloser{}
	DeferCloseReader(rc)
	if !rc.closed {
		t.Error("DeferCloseReader() не вызвал Close()")
	}
}

func TestPPrint(t *testing.T) {
	var buf bytes.Buffer
	PPrint(&buf, struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}{"alex", 3})

	want := "{\n  \"name\": \"alex\",\n  \"age\": 3\n}"
	if got := strings.TrimSpace(buf.String()); got != want {
		t.Errorf("PPrint() =\n%s\nwant:\n%s", got, want)
	}
}

func TestPPrintUnsupported(t *testing.T) {
	// ошибка сериализации логируется, паники быть не должно
	var buf bytes.Buffer
	PPrint(&buf, make(chan int))
}
