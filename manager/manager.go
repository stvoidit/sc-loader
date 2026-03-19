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

func NewManager(ctx context.Context, config *utils.Config) (*Manager, error) {
	prog, err := exec.LookPath("ffmpeg")
	if err != nil {
		return nil, err
	}
	client := client.NewClient(config)
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

func (m *Manager) RecordStream(ctx context.Context, streamer Streamer) (writedBytes int, err error) {
	roomID, online, err := m.client.GetRoomID(ctx, streamer.Username)
	if err != nil {
		slog.Error("main", "streamer", streamer, slog.String("error", err.Error()))
		return writedBytes, err
	}
	if !online {
		slog.Warn("main", "streamer", streamer, slog.String("status", "offline"))
		return writedBytes, err
	}
	slog.Debug("main", slog.Int("roomID", roomID), slog.String("cdn", m.bestCDN))
	plist, pkey, err := m.client.GetPlaylistVariants(ctx, m.bestCDN, roomID)
	if err != nil {
		slog.Error("main.GetPlaylistVariants", slog.String("error", err.Error()))
		return writedBytes, err
	}
	slog.Debug("main", "plist", plist)
	slog.Info("RecordStream", slog.String("cdn", m.bestCDN))
	outputFilename := m.GetFullPathVideoFile(streamer)
	ch := m.client.StartPlaylistLoop(ctx, streamer.Username, plist, pkey)
	return m.StartPipeFFmpeg(ctx, outputFilename, ch)
}

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
	remuxFileName := strings.Replace(filename, "_tmp", "", 1)
	args := []string{
		"-vaapi_device", "/dev/dri/renderD128",
		"-hwaccel", "vaapi",
		"-hide_banner", "-v", "error", "-stats",
		"-i", "-",
		"-c:a", "copy",
		"-vf", "format=nv12|vaapi,hwupload",
		"-c:v", "hevc_vaapi",
		"-qp", "26",
		"-profile:v", "main",
		// "-sei", "+timing+recovery_point",
		remuxFileName,
	}
	cmd := exec.Command(m.ffmpeg, args...)
	cmd.WaitDelay = time.Second * 10
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	pw, err := cmd.StdinPipe()
	if err != nil {
		return writedBytes, err
	}
	if err := cmd.Start(); err != nil {
		return writedBytes, err
	}
loop:
	for vid := range ch {
		buf, err := m.client.DownloadBuf(ctx, vid)
		if err != nil {
			slog.Error("DownloadBuf", "error", err.Error())
			break loop
		}
		n, err := pw.Write(buf)
		if err != nil {
			slog.Error("GetPlaylistVideo", "error", err.Error())
			break loop
		}
		writedBytes += n
	}
	if err := pw.Close(); err != nil {
		slog.Warn("pw.Close", slog.String("error", err.Error()))
	}
	return writedBytes, cmd.Wait()
}
