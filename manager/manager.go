package manager

import (
	"context"
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

func (m Manager) GetFullPathVideoFile(streamer Streamer) string {
	filename := streamer.MakeFilename()
	return m.config.MakeFilePath(filename)
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
	plist, pkey, err := m.client.GetPlaylistVariants(ctx, m.bestCDN, roomID)
	if err != nil {
		slog.Error("main.GetPlaylistVariants", slog.String("error", err.Error()))
		return err
	}
	slog.Debug("main", "plist", plist)
	slog.Info("RecordStream", slog.String("cdn", m.bestCDN))
	outputFilename := m.GetFullPathVideoFile(streamer)
	ch := m.client.StartPlaylistLoop(ctx, streamer.Username, plist, pkey)
	if _, err := m.StartPipeFFmpeg(ctx, outputFilename, ch); err != nil {
		return err
	}
	return nil
}

// func logStat(size int, streamer Streamer, start time.Time) {
// 	hrSize := utils.FormatFileSize(size)
// 	duration := time.Since(start).Round(time.Second)
// 	// _ = username
// 	fmt.Printf("\r%s %s %s\r", streamer.SreamURL, hrSize, duration)
// }

// func finalLog(size int, streamer Streamer, start time.Time) {
// 	hrSize := utils.FormatFileSize(size)
// 	duration := time.Since(start).Round(time.Second)
// 	_, _ = os.Stdout.WriteString("\r\n")
// 	slog.Info("stop recording",
// 		"streamer", streamer,
// 		slog.String("size", hrSize),
// 		slog.String("duration", duration.String()))
// }

// func (m Manager) finalRemux(filename string) {
// 	_, _ = os.Stdout.WriteString("\r\n")
// 	remuxFileName := strings.Replace(filename, "_tmp", "", 1)
// 	args := []string{
// 		"-hide_banner", "-v", "error", "-stats",
// 		"-i", filename, "-c", "copy",
// 		"-movflags", "+faststart",
// 		remuxFileName,
// 	}
// 	cmd := exec.Command(m.ffmpeg, args...)
// 	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stdout
// 	if err := cmd.Run(); err != nil {
// 		slog.Error("finalRemux.Run", slog.String("error", err.Error()), slog.String("filename", filename))
// 		return
// 	}
// 	if err := os.Remove(filename); err != nil {
// 		slog.Error("fRemoven", slog.String("error", err.Error()), slog.String("filename", filename))
// 		return
// 	}
// 	slog.Info("finalRemux", slog.String("filename", remuxFileName))
// }

func (m Manager) StartPipeFFmpeg(ctx context.Context, filename string, ch <-chan string) (writedBytes int, err error) {
	_, _ = os.Stdout.WriteString("\r\n")
	remuxFileName := strings.Replace(filename, "_tmp", "", 1)
	args := []string{
		"-hide_banner", "-v", "error", "-stats",
		"-i", "-", "-c", "copy",
		"-movflags", "+faststart",
		remuxFileName,
	}
	cmd := exec.Command(m.ffmpeg, args...)
	cmd.WaitDelay = time.Second * 5
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stdout
	pw, err := cmd.StdinPipe()
	if err != nil {
		return writedBytes, err
	}
	if err := cmd.Start(); err != nil {
		return writedBytes, err
	}
	defer func() {
		if err := pw.Close(); err != nil {
			slog.Error("pw.Close", slog.String("error", err.Error()))
		}
		if err := cmd.Wait(); err != nil && !strings.Contains(err.Error(), "exit status") {
			slog.Error("cmd.Wait", slog.String("error", err.Error()))
		}
	}()
	for vid := range ch {
		if !m.client.IsNewURL(vid) {
			continue
		}
		buf, err := m.client.DownloadBuf(ctx, vid)
		if err != nil {
			return writedBytes, err
		}
		if n, err := pw.Write(buf); err != nil {
			slog.Error("GetPlaylistVideo", "error", err.Error())
			return writedBytes, err
		} else {
			writedBytes += n
			m.client.SetCheckURL(vid)
		}
	}
	return writedBytes, nil
}
