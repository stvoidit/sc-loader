package utils

import (
	"crypto/sha256"
	"encoding/json/v2"
	"errors"
	"log/slog"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const defaultFilename = "config.json"

func findConfigFile(filename string) (string, error) {
	configDir, err := os.UserConfigDir()
	if err != nil {
		return filename, err
	}
	pwd, err := os.Getwd()
	if err != nil {
		return filename, err
	}
	slog.Debug("findConfigFile", "configDir", configDir)
	var search = [3]string{
		filepath.Join(configDir, "sc-loader", filename),
		filepath.Join(pwd, filename),
		filename,
	}
	slog.Debug("findConfigFile", slog.Any("search", search))
	for _, p := range search {
		result, err := filepath.Glob(p)
		if err != nil {
			return filename, err
		}
		if len(result) > 0 {
			return result[0], err
		}
	}
	return filename, nil
}

func ReadConfig(filename string) (cnf Config, err error) {
	if filename == "" {
		filename = defaultFilename
	}
	slog.Info("ReadConfig", slog.String("filename", filename))
	b, err := os.ReadFile(filename)
	if err != nil {
		return cnf, err
	}
	if err := json.Unmarshal(b, &cnf, json.StringifyNumbers(true)); err != nil {
		return cnf, err
	}
	if err := cnf.checkFolder(); err != nil {
		return cnf, err
	}
	return cnf, err
}

func LoadConfig() (cnf Config, err error) {
	filename, err := findConfigFile(defaultFilename)
	if err != nil {
		return cnf, err
	}
	return ReadConfig(filename)
}

type PSCH [32]byte
type KeysDRM map[string]PSCH

func (keys KeysDRM) LoadKeys(m map[string]string) {
	if len(m) == 0 {
		return
	}
	for k, v := range m {
		keys[k] = sha256.Sum256([]byte(v))
	}
}

var ErrKeyNotFound = errors.New("decode key not found")

func (keys KeysDRM) GetKey(key string) (PSCH, error) {
	value, ok := keys[key]
	if !ok {
		return value, ErrKeyNotFound
	}
	return value, nil
}

type Config struct {
	Debug     bool              `json:"debug"`
	Host      string            `json:"host"`
	Folder    string            `json:"folder"`
	Keys      map[string]string `json:"keys"`
	CacheKeys KeysDRM           `json:"-"`
}

func (cnf *Config) checkFolder() (err error) {
	cnf.Folder, err = filepath.Abs(cnf.Folder)
	if err != nil {
		return err
	}
	if err := os.MkdirAll(cnf.Folder, 0700); err != nil && !os.IsExist(err) {
		return err
	}
	cnf.cacheKeys()
	return
}

func (cnf *Config) cacheKeys() {
	cnf.CacheKeys = make(KeysDRM, len(cnf.Keys))
	cnf.CacheKeys.LoadKeys(cnf.Keys)
}

func (cnf Config) String() string {
	var sb strings.Builder
	PPrint(&sb, cnf)
	return sb.String()
}

func (cnf Config) CreateVideoFile(username string) (file *os.File, err error) {
	const format = "20060102-150405"
	timeSuffix := time.Now().Format(format)
	filename := username + "_" + timeSuffix + "_tmp.mp4"
	fullPath := filepath.Join(cnf.Folder, filename)
	return os.OpenFile(fullPath, os.O_WRONLY|os.O_CREATE|os.O_APPEND|os.O_SYNC, 0644)
}

func (cnf Config) GetUserURL(username string) string {
	link := url.URL{
		Scheme: "https",
		Host:   cnf.Host,
		Path:   username,
	}
	return link.String()
}
