package utils

import (
	"io"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync"
	"testing"
)

func TestFormatFileSize(t *testing.T) {
	tests := []struct {
		size int64
		want string
	}{
		{0, "0 B"},
		{-5, "-5 B"},
		{512, "512 B"},
		{KB - 1, "1023 B"},
		{KB, "1.0 KB"},
		{KB + KB/2, "1.5 KB"},
		{MB - 1, "1024.0 KB"},
		{MB, "1.0 MB"},
		{MB + MB/4, "1.2 MB"},
		{GB, "1.0 GB"},
		{3*GB + GB/2, "3.5 GB"},
		{TB, "1.0 TB"},
		{2*TB + TB/2, "2.5 TB"},
	}
	for _, tt := range tests {
		if got := FormatFileSize(tt.size); got != tt.want {
			t.Errorf("FormatFileSize(%d) = %q, want %q", tt.size, got, tt.want)
		}
	}
	//generic-вариант для обычного int
	if got := FormatFileSize(2048); got != "2.0 KB" {
		t.Errorf("FormatFileSize(int) = %q, want %q", got, "2.0 KB")
	}
}

func TestSaveFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "dump.json")
	const content = `{"id":1}`

	rc := SaveFile(path, io.NopCloser(strings.NewReader(content)))
	defer DeferCloseReader(rc)

	got, err := io.ReadAll(rc)
	if err != nil {
		t.Fatalf("ReadAll() error = %v", err)
	}
	if string(got) != content {
		t.Errorf("SaveFile() reader = %q, want %q", got, content)
	}

	onDisk, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("os.ReadFile() error = %v", err)
	}
	if string(onDisk) != content {
		t.Errorf("файл = %q, want %q", onDisk, content)
	}
}

func TestSaveFilePanicsOnError(t *testing.T) {
	defer func() {
		if r := recover(); r == nil {
			t.Error("SaveFile() не запаниковал при ошибке создания файла")
		}
	}()
	// каталог не существует -> os.Create вернёт ошибку
	SaveFile(filepath.Join(t.TempDir(), "no-such-dir", "x.json"),
		io.NopCloser(strings.NewReader("x")))
}

func TestSet(t *testing.T) {
	s := NewSet[string]()

	if s.Has("a") {
		t.Error("Has(a) = true для нового множества")
	}
	s.Set("a")
	s.Set("b")
	s.Set("a") // повтор не должен ничего ломать
	if !s.Has("a") || !s.Has("b") {
		t.Errorf("Has() = (%v, %v), want (true, true)", s.Has("a"), s.Has("b"))
	}

	values := s.Values()
	slices.Sort(values)
	if !slices.Equal(values, []string{"a", "b"}) {
		t.Errorf("Values() = %v, want [a b]", values)
	}

	t.Run("FilterNew", func(t *testing.T) {
		got := s.FilterNew([]string{"a", "c", "b", "d"})
		slices.Sort(got)
		if !slices.Equal(got, []string{"c", "d"}) {
			t.Errorf("FilterNew() = %v, want [c d]", got)
		}
		// FilterNew только фильтрует, но не добавляет значения в множество
		if s.Has("c") || s.Has("d") {
			t.Error("FilterNew() изменил множество")
		}
		if got := s.FilterNew(nil); len(got) != 0 {
			t.Errorf("FilterNew(nil) = %v, want empty", got)
		}
		if got := s.FilterNew([]string{"a"}); len(got) != 0 {
			t.Errorf("FilterNew([a]) = %v, want empty", got)
		}
	})

	t.Run("String", func(t *testing.T) {
		empty := NewSet[int]()
		if got := empty.String(); got != "[]" {
			t.Errorf("String() = %q, want %q", got, "[]")
		}
		one := NewSet[string]()
		one.Set("x")
		if got := one.String(); got != `["x"]` {
			t.Errorf("String() = %q, want %q", got, `["x"]`)
		}
	})
}

func TestSetConcurrent(t *testing.T) {
	s := NewSet[int]()
	var wg sync.WaitGroup
	for i := range 50 {
		wg.Go(func() {
			s.Set(i)
			s.Has(i)
			s.FilterNew([]int{i, i + 1})
			_ = s.Values()
			_ = s.String()
		})
	}
	wg.Wait()
	if len(s.Values()) != 50 {
		t.Errorf("Values() len = %d, want 50", len(s.Values()))
	}
}
