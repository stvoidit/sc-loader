package utils

import (
	"bytes"
	"log/slog"
	"os"
	"strings"
	"testing"
)

func TestSlogParseEnvLevel(t *testing.T) {
	const envName = "TEST_LOG_LEVEL"

	tests := []struct {
		value string
		want  slog.Level
	}{
		{"debug", slog.LevelDebug},
		{"DBG", slog.LevelDebug},
		{"d", slog.LevelDebug},
		{"-4", slog.LevelDebug},
		{"info", slog.LevelInfo},
		{"i", slog.LevelInfo},
		{"0", slog.LevelInfo},
		{"warn", slog.LevelWarn},
		{"Warning", slog.LevelWarn},
		{"w", slog.LevelWarn},
		{"4", slog.LevelWarn},
		{"err", slog.LevelError},
		{"ERROR", slog.LevelError},
		{"e", slog.LevelError},
		{"8", slog.LevelError},
		{"verbose", slog.LevelInfo}, // неизвестное значение -> info
		{"", slog.LevelInfo},
	}
	for _, tt := range tests {
		t.Run(tt.value, func(t *testing.T) {
			t.Setenv(envName, tt.value)
			if got := SlogParseEnvLevel(envName); got != tt.want {
				t.Errorf("SlogParseEnvLevel(%q) = %v, want %v", tt.value, got, tt.want)
			}
		})
	}

	t.Run("env not set", func(t *testing.T) {
		old, hadOld := os.LookupEnv(envName)
		if err := os.Unsetenv(envName); err != nil {
			t.Fatal(err)
		}
		defer func() {
			if hadOld {
				os.Setenv(envName, old) //nolint:errcheck
			}
		}()
		if got := SlogParseEnvLevel(envName); got != slog.LevelInfo {
			t.Errorf("SlogParseEnvLevel() = %v, want %v", got, slog.LevelInfo)
		}
	})
}

func TestNewSloglogger(t *testing.T) {
	t.Run("single writer", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewSloglogger(slog.LevelInfo, &buf)
		logger.Info("hello", slog.String("user", "alex"))

		out := buf.String()
		for _, want := range []string{"level=INFO", "msg=hello", "user=alex"} {
			if !strings.Contains(out, want) {
				t.Errorf("лог = %q, want substring %q", out, want)
			}
		}
		// shortSourceFile обрезает путь до имени файла
		if !strings.Contains(out, "source=logging_test.go:") {
			t.Errorf("лог = %q, want source=logging_test.go:NN", out)
		}
		if strings.Contains(out, "sc-loader/utils/logging_test.go") {
			t.Errorf("лог = %q, путь к файлу не обрезан", out)
		}
	})

	t.Run("multi writer", func(t *testing.T) {
		var first, second bytes.Buffer
		logger := NewSloglogger(slog.LevelInfo, &first, &second)
		logger.Info("both")
		if first.Len() == 0 || second.Len() == 0 {
			t.Errorf("запись в несколько writer не сработала: %d/%d", first.Len(), second.Len())
		}
	})

	t.Run("level filter", func(t *testing.T) {
		var buf bytes.Buffer
		logger := NewSloglogger(slog.LevelError, &buf)
		logger.Info("skipped")
		logger.Debug("skipped")
		if buf.Len() != 0 {
			t.Errorf("буфер = %q, want empty", buf.String())
		}
		logger.Error("shown")
		if !strings.Contains(buf.String(), "msg=shown") {
			t.Errorf("буфер = %q, want msg=shown", buf.String())
		}
	})

	t.Run("default stdout", func(t *testing.T) {
		f, err := os.CreateTemp(t.TempDir(), "stdout")
		if err != nil {
			t.Fatal(err)
		}
		defer DeferCloseReader(f)

		old := os.Stdout
		os.Stdout = f
		logger := NewSloglogger(slog.LevelInfo)
		logger.Info("to stdout")
		os.Stdout = old

		got, err := os.ReadFile(f.Name())
		if err != nil {
			t.Fatal(err)
		}
		if !strings.Contains(string(got), "msg=\"to stdout\"") {
			t.Errorf("stdout = %q, want msg=\"to stdout\"", got)
		}
	})
}

func TestShortSourceFile(t *testing.T) {
	t.Run("source attr", func(t *testing.T) {
		src := &slog.Source{File: "/home/user/project/main.go", Line: 42, Function: "main.main"}
		a := slog.Any(slog.SourceKey, src)
		shortSourceFile(&a)
		if src.File != "main.go" {
			t.Errorf("Source.File = %q, want %q", src.File, "main.go")
		}
		if src.Line != 42 {
			t.Errorf("Source.Line = %d, want 42", src.Line)
		}
	})

	t.Run("other attr is untouched", func(t *testing.T) {
		a := slog.String(slog.MessageKey, "hello")
		shortSourceFile(&a)
		if a.Key != slog.MessageKey || a.Value.String() != "hello" {
			t.Errorf("attr = %v, want без изменений", a)
		}
	})

	t.Run("source key with wrong type", func(t *testing.T) {
		a := slog.String(slog.SourceKey, "not a source")
		shortSourceFile(&a) // не должен паниковать
		if a.Value.String() != "not a source" {
			t.Errorf("attr = %v, want без изменений", a)
		}
	})
}

func TestNewSlogHandlerOptions(t *testing.T) {
	called := false
	opt := ReplaceOption(func(a *slog.Attr) {
		if a.Key == slog.MessageKey {
			called = true
			a.Value = slog.StringValue(strings.ToUpper(a.Value.String()))
		}
	})

	opts := NewSlogHandlerOptions(slog.LevelWarn, opt)
	if !opts.AddSource {
		t.Error("AddSource = false, want true")
	}
	if opts.Level != slog.LevelWarn {
		t.Errorf("Level = %v, want %v", opts.Level, slog.LevelWarn)
	}

	var buf bytes.Buffer
	slog.New(slog.NewTextHandler(&buf, &opts)).Warn("quiet")
	if !called {
		t.Error("ReplaceOption не вызван")
	}
	if !strings.Contains(buf.String(), "msg=QUIET") {
		t.Errorf("буфер = %q, want msg=QUIET", buf.String())
	}
}
