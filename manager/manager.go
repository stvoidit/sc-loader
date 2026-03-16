package manager

import (
	"context"
	"fmt"
	"log/slog"
	"os"
	"os/exec"
	"sc-loader/client"
	"sc-loader/utils"
	"strings"
	"time"
)

func NewManager(ctx context.Context, client *client.ClientAPI, config *utils.Config) (*Manager, error) {
	prog, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}
	InitialConfig, err := client.InitClientConfig(ctx)
	if err != nil {
		return nil, err
	}
	domains := InitialConfig.GetDomainsCDN()
	bestCDN := client.PickBestHTTPS(ctx, domains)
	return &Manager{
		client:  client,
		config:  config,
		ffmpeg:  prog,
		bestCDN: bestCDN,
	}, nil
}

type Manager struct {
	client  *client.ClientAPI
	config  *utils.Config
	ffmpeg  string
	bestCDN string
}

func (m *Manager) RecordStream(ctx context.Context, streamer Streamer) (err error) {
	sl := slog.With("streamer", streamer)
	roomID, online, err := m.client.GetRoomID(ctx, streamer.Username)
	if err != nil {
		sl.Error("main", slog.String("error", err.Error()))
		return err
	}
	if !online {
		sl.Warn("main", slog.String("status", "offline"))
		return err
	}
	slog.Debug("main", slog.Int("roomID", roomID), slog.String("cdn", m.bestCDN))

	// defer func() {
	// 	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	// 	defer cancel()
	// 	if status, online, err := m.client.GetRoomStatus(ctx, streamer.Username); err != nil && !errors.Is(err, context.Canceled) {
	// 		slog.Error("GetRoomStatus", slog.String("error", err.Error()))
	// 	} else {
	// 		slog.Info("GetRoomStatus", slog.Bool("online", online), slog.String("status", status))
	// 	}
	// }()

	plist, pkey, err := m.client.GetPlaylistVariants(ctx, m.bestCDN, roomID)
	if err != nil {
		slog.Error("main.GetPlaylistVariants", slog.String("error", err.Error()))
		return err
	}
	slog.Debug("main", "plist", plist)
	slog.Info("RecordStream", slog.String("cdn", m.bestCDN))

	f, err := m.config.CreateVideoFile(streamer.MakeFilename())
	if err != nil {
		return err
	}
	defer m.finalRemux(f.Name())
	defer utils.DeferCloseReader(f)

	var start = time.Now()
	ch := m.client.StartPlaylistLoop(ctx, streamer.Username, plist, pkey)

	var writedBytes int
	for vid := range ch {
		if n, err := m.client.Download(ctx, f, vid); err != nil && !utils.IsCancel(err) {
			slog.Error("GetPlaylistVideo", "error", err.Error())
			return err
		} else {
			writedBytes += n
		}
		logStat(writedBytes, streamer, start)
	}
	return err
}

func logStat(size int, streamer Streamer, start time.Time) {
	hrSize := utils.FormatFileSize(size)
	duration := time.Since(start).Round(time.Second)
	// _ = username
	fmt.Printf("\r%s %s %s\r", streamer.SreamURL, hrSize, duration)
}

func (m Manager) finalRemux(filename string) {
	_, _ = os.Stdout.WriteString("\r\n")
	remuxFileName := strings.Replace(filename, "_tmp", "", 1)
	args := []string{
		"-hide_banner", "-v", "error", "-stats",
		"-i", filename, "-c", "copy",
		"-movflags", "+faststart",
		remuxFileName,
	}
	cmd := exec.Command(m.ffmpeg, args...)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stdout
	if err := cmd.Run(); err != nil {
		slog.Error("finalRemux.Run", slog.String("error", err.Error()), slog.String("filename", filename))
		return
	}
	if err := os.Remove(filename); err != nil {
		slog.Error("fRemoven", slog.String("error", err.Error()), slog.String("filename", filename))
		return
	}
	slog.Info("finalRemux", slog.String("filename", remuxFileName))
}
