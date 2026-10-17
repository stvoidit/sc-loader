package manager

import (
	"context"
	"io"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"path/filepath"
	"slices"
	"strings"
	"testing"

	"sc-loader/client"
	"sc-loader/models"
	"sc-loader/utils"
)

// TestMain глушит slog-вывод тестируемого кода, чтобы не засорять лог теста.
func TestMain(m *testing.M) {
	slog.SetDefault(slog.New(slog.NewTextHandler(io.Discard, nil)))
	os.Exit(m.Run())
}

func TestGetFullPathVideoFile(t *testing.T) {
	folder := t.TempDir()
	m := Manager{config: &utils.Config{Folder: folder}}

	got := m.GetFullPathVideoFile(models.Streamer{Username: "alex"})
	if filepath.Dir(got) != folder {
		t.Errorf("GetFullPathVideoFile() dir = %q, want %q", filepath.Dir(got), folder)
	}
	base := filepath.Base(got)
	if !strings.HasPrefix(base, "alex_") || !strings.HasSuffix(base, "_tmp.mp4") {
		t.Errorf("GetFullPathVideoFile() = %q, want alex_*_tmp.mp4", got)
	}
}

func TestCreateFFmpegProcess(t *testing.T) {
	truePath, err := exec.LookPath("true")
	if err != nil {
		t.Skip("true не найден")
	}

	m := Manager{ffmpeg: truePath}
	filename := filepath.Join(t.TempDir(), "alex_20260101-000000_tmp.mp4")

	cmd, pw, err := m.createFFmpegProcess(filename)
	if err != nil {
		t.Fatalf("createFFmpegProcess() error = %v", err)
	}

	// "_tmp" из имени файла должен быть убран
	if got := cmd.Args[len(cmd.Args)-1]; filepath.Base(got) != "alex_20260101-000000.mp4" {
		t.Errorf("выходной файл = %q, want alex_20260101-000000.mp4", got)
	}
	want := []string{
		"-vaapi_device", "/dev/dri/renderD128",
		"-hwaccel", "vaapi",
		"-i", "-",
		"-c:a", "copy",
		"-c:v", "hevc_vaapi",
		"-qp", "31",
		"-profile:v", "main",
	}
	for _, arg := range want {
		if !slices.Contains(cmd.Args, arg) {
			t.Errorf("в аргументах ffmpeg нет %q: %v", arg, cmd.Args)
		}
	}

	if err := pw.Close(); err != nil {
		t.Logf("pw.Close() error = %v", err)
	}
	if err := cmd.Wait(); err != nil {
		t.Errorf("cmd.Wait() error = %v", err)
	}

	t.Run("missing binary", func(t *testing.T) {
		broken := Manager{ffmpeg: filepath.Join(t.TempDir(), "no-such-ffmpeg")}
		if _, _, err := broken.createFFmpegProcess("out_tmp.mp4"); err == nil {
			t.Error("createFFmpegProcess() error = nil, want error")
		}
	})
}

func TestStartPipeFFmpegError(t *testing.T) {
	m := Manager{ffmpeg: filepath.Join(t.TempDir(), "no-such-ffmpeg")}
	ch := make(chan string)
	close(ch)

	if _, err := m.StartPipeFFmpeg(context.Background(), "out_tmp.mp4", ch); err == nil {
		t.Error("StartPipeFFmpeg() error = nil, want error")
	}
}

// TestStartPipeFFmpegStub проверяет связку createFFmpegProcess + streamToFFmpeg
// на любом бинарнике: cat не понимает аргументы ffmpeg и завершается с ошибкой.
func TestStartPipeFFmpegStub(t *testing.T) {
	catPath, err := exec.LookPath("cat")
	if err != nil {
		t.Skip("cat не найден")
	}
	m := Manager{
		ffmpeg: catPath,
		client: client.NewClient(&utils.Config{Host: "example.com"}),
	}

	ch := make(chan string)
	close(ch)
	out := filepath.Join(t.TempDir(), "alex_20260101-000000_tmp.mp4")

	if _, err := m.StartPipeFFmpeg(context.Background(), out, ch); err == nil {
		t.Error("StartPipeFFmpeg() error = nil, want ошибку процесса-заглушки")
	}
}

// partsMux отдаёт фиксированные тела сегментов.
func partsMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/p1", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "AAA") //nolint:errcheck
	})
	mux.HandleFunc("/p2", func(w http.ResponseWriter, r *http.Request) {
		io.WriteString(w, "BBBB") //nolint:errcheck
	})
	return mux
}

func partsServer(t *testing.T) *httptest.Server {
	t.Helper()
	srv := httptest.NewServer(partsMux())
	t.Cleanup(srv.Close)
	return srv
}

// startCat подменяет ffmpeg на cat: пишет вывод в outPath и возвращает
// процесс вместе с pipe для stdin.
func startCat(t *testing.T, catPath, outPath string) (*exec.Cmd, io.WriteCloser) {
	t.Helper()
	out, err := os.Create(outPath)
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { out.Close() }) //nolint:errcheck

	cmd := exec.Command(catPath)
	cmd.Stdout = out
	pw, err := cmd.StdinPipe()
	if err != nil {
		t.Fatal(err)
	}
	if err := cmd.Start(); err != nil {
		t.Fatal(err)
	}
	return cmd, pw
}

func TestStreamToFFmpeg(t *testing.T) {
	catPath, err := exec.LookPath("cat")
	if err != nil {
		t.Skip("cat не найден")
	}
	srv := partsServer(t)
	m := Manager{client: client.NewClient(&utils.Config{Host: "example.com"}), ffmpeg: catPath}

	t.Run("writes all parts", func(t *testing.T) {
		outPath := filepath.Join(t.TempDir(), "out.mp4")
		cmd, pw := startCat(t, catPath, outPath)

		ch := make(chan string, 2)
		ch <- srv.URL + "/p1"
		ch <- srv.URL + "/p2"
		close(ch)

		n, err := m.streamToFFmpeg(context.Background(), cmd, pw, ch)
		if err != nil {
			t.Fatalf("streamToFFmpeg() error = %v", err)
		}
		if n != 7 {
			t.Errorf("streamToFFmpeg() n = %d, want 7", n)
		}
		got, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatal(err)
		}
		if string(got) != "AAABBBB" {
			t.Errorf("файл = %q, want %q", got, "AAABBBB")
		}
	})

	t.Run("download error breaks loop", func(t *testing.T) {
		dead := httptest.NewServer(partsMux())
		deadURL := dead.URL
		dead.Close() // хост недоступен -> ошибка загрузки

		cmd, pw := startCat(t, catPath, filepath.Join(t.TempDir(), "out.mp4"))
		ch := make(chan string, 2)
		ch <- srv.URL + "/p1"
		ch <- deadURL + "/p1"
		close(ch)

		n, err := m.streamToFFmpeg(context.Background(), cmd, pw, ch)
		if err != nil {
			t.Fatalf("streamToFFmpeg() error = %v", err)
		}
		if n != 3 {
			t.Errorf("streamToFFmpeg() n = %d, want 3 (второй сегмент не загружен)", n)
		}
	})

	t.Run("write error breaks loop", func(t *testing.T) {
		truePath, err := exec.LookPath("true")
		if err != nil {
			t.Skip("true не найден")
		}
		cmd := exec.Command(truePath)
		if err := cmd.Start(); err != nil {
			t.Fatal(err)
		}

		// закрытые концы pipe дают гарантированную ошибку и на Write, и на Close
		pr, pw, err := os.Pipe()
		if err != nil {
			t.Fatal(err)
		}
		if err := pr.Close(); err != nil {
			t.Fatal(err)
		}
		if err := pw.Close(); err != nil {
			t.Fatal(err)
		}

		ch := make(chan string, 1)
		ch <- srv.URL + "/p1"
		close(ch)

		n, err := m.streamToFFmpeg(context.Background(), cmd, pw, ch)
		if err != nil {
			t.Fatalf("streamToFFmpeg() error = %v", err)
		}
		if n != 0 {
			t.Errorf("streamToFFmpeg() n = %d, want 0", n)
		}
	})

	t.Run("empty channel", func(t *testing.T) {
		outPath := filepath.Join(t.TempDir(), "empty.mp4")
		cmd, pw := startCat(t, catPath, outPath)

		ch := make(chan string)
		close(ch)

		n, err := m.streamToFFmpeg(context.Background(), cmd, pw, ch)
		if err != nil {
			t.Fatalf("streamToFFmpeg() error = %v", err)
		}
		if n != 0 {
			t.Errorf("streamToFFmpeg() n = %d, want 0", n)
		}
		got, err := os.ReadFile(outPath)
		if err != nil {
			t.Fatal(err)
		}
		if len(got) != 0 {
			t.Errorf("файл = %q, want empty", got)
		}
	})
}

func TestNewManagerError(t *testing.T) {
	cnf := &utils.Config{Host: "127.0.0.1:1", Folder: t.TempDir()}
	if _, err := NewManager(context.Background(), cnf); err == nil {
		t.Error("NewManager() error = nil, want error")
	}
}
