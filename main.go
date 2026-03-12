package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sc-loader/client"
	"sc-loader/utils"
	"strings"
	"time"
)

var logLevel = slog.LevelInfo

func init() {
	const logLevelVar = "LOG_LEVEL"
	logLevel = utils.SlogParseEnvLevel(logLevelVar)
	slog.SetLogLoggerLevel(logLevel)
}

func finalRemux(filename string) {
	prog, err := exec.LookPath("ffmpeg")
	if err != nil {
		slog.Error("finalRemux.LookPath", slog.String("error", err.Error()), slog.String("filename", filename))
		os.Exit(1)
	}
	remuxFileName := strings.Replace(filename, "_tmp", "", 1)
	args := []string{
		"-hide_banner", "-v", "error", "-stats",
		"-i", filename, "-c", "copy",
		"-movflags", "+faststart",
		remuxFileName,
	}
	_, _ = os.Stdout.WriteString("\n")
	cmd := exec.Command(prog, args...)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stdout
	if err := cmd.Run(); err != nil {
		slog.Error("finalRemux.Run", slog.String("error", err.Error()), slog.String("filename", filename))
		os.Exit(1)
	}
	if err := os.Remove(filename); err != nil {
		slog.Error("fRemoven", slog.String("error", err.Error()), slog.String("filename", filename))
	}
	slog.Info("finalRemux", slog.String("filename", remuxFileName))
}

func RecordStream(ctx context.Context, config utils.Config, username string) (err error) {
	c := client.NewClient(config)
	id, err := c.InitClientConfig(ctx)
	if err != nil {
		slog.Error("main", slog.String("error", err.Error()))
		return
	}
	domains := id.GetDomainsCDN()
	cdnDomain := c.PickBestHTTPS(ctx, domains)
	roomID, online, err := c.GetRoomID(ctx, username)
	if err != nil {
		slog.Error("main", slog.String("error", err.Error()))
		return err
	}
	if !online {
		slog.Warn("main", slog.String("status", "offline"), slog.String("username", username))
		return err
	}
	var usernameURL = config.GetUserURL(username)
	slog.Debug("main", slog.Int("roomID", roomID), slog.String("username", username), slog.String("cdn", cdnDomain))

	defer func() {
		ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
		defer cancel()
		if status, online, err := c.GetRoomStatus(ctx, username); err != nil && !errors.Is(err, context.Canceled) {
			slog.Error("GetRoomStatus", slog.String("error", err.Error()), slog.String("url", usernameURL))
		} else {
			slog.Info("GetRoomStatus", slog.Bool("online", online), slog.String("status", status), slog.String("url", usernameURL))
		}
	}()

	plist, pkey, err := c.GetPlaylistVariants(ctx, cdnDomain, roomID)
	if err != nil {
		slog.Error("main.GetPlaylistVariants", slog.String("error", err.Error()))
		return err
	}
	slog.Debug("main", "plist", plist)
	slog.Info("RecordStream", slog.String("username", username), slog.String("cdn", cdnDomain))

	f, err := config.CreateVideoFile(username)
	if err != nil {
		return err
	}
	defer finalRemux(f.Name())
	defer utils.DeferCloseReader(f)

	var start = time.Now()
	ch := c.StartPlaylistLoop(ctx, plist, pkey)

	for vid := range ch {
		if err := c.Download(ctx, f, vid); err != nil && !utils.IsCancel(err) {
			slog.Error("GetPlaylistVideo", "error", err.Error())
			return err
		}
		logStat(f, username, usernameURL, start)
	}
	// wg.Wait()
	return err
}

func logStat(f *os.File, username, url string, start time.Time) {
	size, err := utils.FileStat(f)
	if err != nil {
		slog.Error("logStat", slog.String("error", err.Error()))
		return
	}
	duration := time.Since(start).Round(time.Second)

	// slog.Info(username,
	// 	slog.String("size", size),
	// 	slog.Duration("duration", duration),
	// 	slog.String("url", url),
	// )
	_ = username
	fmt.Printf("\r%s %s %s\r", url, size, duration)
}

func parseUsername() (string, error) {
	if len(os.Args) < 2 {
		return "", errors.New("please set username")
	}
	username := os.Args[1]
	if strings.HasPrefix(username, "http") {
		URL, err := url.Parse(username)
		if err != nil {
			return "", err
		}
		return filepath.Base(URL.Path), nil
	}
	return username, nil
}

func main() {
	var username, err = parseUsername()
	if err != nil {
		panic(err)
	}

	config, err := utils.LoadConfig()
	if err != nil {
		panic(err)
	}
	outs := []io.Writer{os.Stderr}
	if logLevel == slog.LevelDebug {
		logname := time.Now().Format(time.RFC3339)
		f, err := os.Create(fmt.Sprintf(`%s_%s.log`, username, logname))
		if err != nil {
			panic(err)
		}
		defer utils.DeferCloseReader(f)
		outs = append(outs, f)
	}
	slog.SetDefault(utils.NewSloglogger(logLevel, outs...))
	defer func() {
		slog.Info(config.GetUserURL(username))
	}()
	if config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()
	if err := RecordStream(ctx, config, username); err != nil {
		slog.Error("RecordStream", slog.String("error", err.Error()))
	}
}
