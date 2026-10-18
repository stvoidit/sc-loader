package utils

import (
	"crypto/sha256"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestKeysDRMLoadKeys(t *testing.T) {
	keys := make(KeysDRM)
	keys.LoadKeys(map[string]string{drmKey: drmKey, "empty-value": ""})

	if len(keys) != 2 {
		t.Fatalf("LoadKeys() len = %d, want 2", len(keys))
	}
	if want := sha256.Sum256([]byte(drmKey)); keys[drmKey] != want {
		t.Errorf("LoadKeys() ключ = %x, want sha256 = %x", keys[drmKey], want)
	}
	if want := sha256.Sum256(nil); keys["empty-value"] != want {
		t.Errorf("LoadKeys() для пустого значения = %x, want %x", keys["empty-value"], want)
	}
}

func TestKeysDRMLoadKeysEmpty(t *testing.T) {
	keys := make(KeysDRM)
	keys.LoadKeys(nil)
	keys.LoadKeys(map[string]string{})
	if len(keys) != 0 {
		t.Errorf("LoadKeys(nil) len = %d, want 0", len(keys))
	}
}

func TestKeysDRMGetKey(t *testing.T) {
	keys := KeysDRM{drmKey: testPSCH(drmKey)}

	got, err := keys.GetKey(drmKey)
	if err != nil {
		t.Fatalf("GetKey() error = %v", err)
	}
	if got != testPSCH(drmKey) {
		t.Errorf("GetKey() = %x, want %x", got, testPSCH(drmKey))
	}

	if _, err := keys.GetKey("unknown"); !errors.Is(err, ErrKeyNotFound) {
		t.Errorf("GetKey(unknown) error = %v, want %v", err, ErrKeyNotFound)
	}
}

const testConfigJSON = `{
  "debug": true,
  "host": "example.com",
  "folder": %q,
  "keys": {%q: "secret"}
}`

func writeTestConfig(t *testing.T, dir, folder string) string {
	t.Helper()
	path := filepath.Join(dir, "config.json")
	body := fmt.Sprintf(testConfigJSON, folder, drmKey)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func TestReadConfig(t *testing.T) {
	dir := t.TempDir()
	folder := filepath.Join(dir, "downloads", "nested")
	path := writeTestConfig(t, dir, folder)

	cnf, err := ReadConfig(path)
	if err != nil {
		t.Fatalf("ReadConfig() error = %v", err)
	}
	if !cnf.Debug {
		t.Error("Debug = false, want true")
	}
	if cnf.Host != "example.com" {
		t.Errorf("Host = %q, want %q", cnf.Host, "example.com")
	}
	if cnf.Folder != folder {
		t.Errorf("Folder = %q, want %q", cnf.Folder, folder)
	}
	if _, err := os.Stat(folder); err != nil {
		t.Errorf("checkFolder() не создал каталог: %v", err)
	}
	if got, err := cnf.CacheKeys.GetKey(drmKey); err != nil || got != testPSCH("secret") {
		t.Errorf("CacheKeys.GetKey() = (%x, %v), want sha256(secret)", got, err)
	}
}

func TestReadConfigErrors(t *testing.T) {
	dir := t.TempDir()

	t.Run("missing file", func(t *testing.T) {
		if _, err := ReadConfig(filepath.Join(dir, "nope.json")); err == nil {
			t.Error("ReadConfig() error = nil, want error")
		}
	})

	t.Run("broken json", func(t *testing.T) {
		path := filepath.Join(dir, "broken.json")
		if err := os.WriteFile(path, []byte("{not json"), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadConfig(path); err == nil {
			t.Error("ReadConfig() error = nil, want error")
		}
	})

	t.Run("default filename from cwd", func(t *testing.T) {
		t.Chdir(dir)
		folder := filepath.Join(dir, "dl")
		writeTestConfig(t, dir, folder)
		cnf, err := ReadConfig("")
		if err != nil {
			t.Fatalf("ReadConfig(\"\") error = %v", err)
		}
		if cnf.Host != "example.com" {
			t.Errorf("Host = %q, want %q", cnf.Host, "example.com")
		}
	})

	t.Run("folder is not a dir", func(t *testing.T) {
		path := filepath.Join(dir, "bad-folder.json")
		body := `{"host":"example.com","folder":"` + filepath.Join(dir, "config.json") + `"}`
		if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
			t.Fatal(err)
		}
		if _, err := ReadConfig(path); err == nil {
			t.Error("ReadConfig() error = nil, want error для файла вместо каталога")
		}
	})
}

func TestConfigMakeFilePath(t *testing.T) {
	cnf := Config{Folder: "/data/video"}
	if got := cnf.MakeFilePath("a.mp4"); got != "/data/video/a.mp4" {
		t.Errorf("MakeFilePath() = %q, want %q", got, "/data/video/a.mp4")
	}
}

func TestConfigString(t *testing.T) {
	cnf := Config{
		Debug:  true,
		Host:   "example.com",
		Folder: "/tmp",
		Keys:   map[string]string{"k": "v"},
	}
	got := cnf.String()
	for _, want := range []string{`"host": "example.com"`, `"debug": true`, `"k": "v"`} {
		if !strings.Contains(got, want) {
			t.Errorf("String() = %s, want substring %q", got, want)
		}
	}
}

func TestConfigLogValue(t *testing.T) {
	cnf := Config{
		Debug:  true,
		Host:   "example.com",
		Folder: "/tmp",
		Keys:   map[string]string{"k": "v"},
	}
	v := cnf.LogValue()
	if v.Kind() != slog.KindGroup {
		t.Fatalf("LogValue().Kind() = %v, want %v", v.Kind(), slog.KindGroup)
	}
	found := map[string]bool{}
	for _, a := range v.Group() {
		found[a.Key] = true
		if a.Key == "host" && a.Value.String() != "example.com" {
			t.Errorf("host = %q, want %q", a.Value.String(), "example.com")
		}
	}
	for _, key := range []string{"debug", "host", "folder", "keys"} {
		if !found[key] {
			t.Errorf("LogValue() не содержит ключ %q", key)
		}
	}
}

func TestFindConfigFile(t *testing.T) {
	t.Run("user config dir", func(t *testing.T) {
		xdg := t.TempDir()
		t.Setenv("XDG_CONFIG_HOME", xdg)
		t.Chdir(t.TempDir())

		dir := filepath.Join(xdg, "sc-loader")
		if err := os.MkdirAll(dir, 0o700); err != nil {
			t.Fatal(err)
		}
		want := filepath.Join(dir, defaultFilename)
		if err := os.WriteFile(want, []byte("{}"), 0o600); err != nil {
			t.Fatal(err)
		}

		got, err := findConfigFile(defaultFilename)
		if err != nil {
			t.Fatalf("findConfigFile() error = %v", err)
		}
		if got != want {
			t.Errorf("findConfigFile() = %q, want %q", got, want)
		}
	})

	t.Run("current dir", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		cwd := t.TempDir()
		t.Chdir(cwd)
		if err := os.WriteFile(
			filepath.Join(cwd, defaultFilename),
			[]byte("{}"),
			0o600,
		); err != nil {
			t.Fatal(err)
		}

		got, err := findConfigFile(defaultFilename)
		if err != nil {
			t.Fatalf("findConfigFile() error = %v", err)
		}
		if got != filepath.Join(cwd, defaultFilename) {
			t.Errorf("findConfigFile() = %q, want %q", got, filepath.Join(cwd, defaultFilename))
		}
	})

	t.Run("not found", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		t.Chdir(t.TempDir())

		got, err := findConfigFile(defaultFilename)
		if err != nil {
			t.Fatalf("findConfigFile() error = %v", err)
		}
		if got != defaultFilename {
			t.Errorf("findConfigFile() = %q, want %q", got, defaultFilename)
		}
	})

	t.Run("user config dir unavailable", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "")

		got, err := findConfigFile(defaultFilename)
		if err == nil {
			t.Fatal("findConfigFile() error = nil, want error")
		}
		if got != defaultFilename {
			t.Errorf("findConfigFile() = %q, want %q", got, defaultFilename)
		}
	})
}

func TestLoadConfig(t *testing.T) {
	xdg := t.TempDir()
	t.Setenv("XDG_CONFIG_HOME", xdg)
	t.Chdir(t.TempDir())

	dir := filepath.Join(xdg, "sc-loader")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	folder := filepath.Join(xdg, "downloads")
	writeTestConfig(t, dir, folder)

	cnf, err := LoadConfig()
	if err != nil {
		t.Fatalf("LoadConfig() error = %v", err)
	}
	if cnf.Host != "example.com" {
		t.Errorf("Host = %q, want %q", cnf.Host, "example.com")
	}
	if cnf.Folder != folder {
		t.Errorf("Folder = %q, want %q", cnf.Folder, folder)
	}
	if _, err := os.Stat(folder); err != nil {
		t.Errorf("каталог для загрузок не создан: %v", err)
	}

	t.Run("missing config", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", t.TempDir())
		t.Chdir(t.TempDir())
		if _, err := LoadConfig(); err == nil {
			t.Error("LoadConfig() error = nil, want error")
		}
	})

	t.Run("user config dir unavailable", func(t *testing.T) {
		t.Setenv("XDG_CONFIG_HOME", "")
		t.Setenv("HOME", "")
		if _, err := LoadConfig(); err == nil {
			t.Error("LoadConfig() error = nil, want error")
		}
	})
}
