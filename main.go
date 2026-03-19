package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"sc-loader/manager"
	"sc-loader/utils"
	"time"
)

var logLevel = slog.LevelInfo

func init() {
	const logLevelVar = "LOG_LEVEL"
	logLevel = utils.SlogParseEnvLevel(logLevelVar)
	slog.SetLogLoggerLevel(logLevel)
}

func main() {
	config, err := utils.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}
	if len(os.Args) < 2 {
		log.Fatal("please set username")
	}
	streamer, err := manager.ParseUsername(os.Args[1], config.Host)
	if err != nil {
		log.Fatal(err)
	}

	outs := []io.Writer{os.Stderr}
	if logLevel == slog.LevelDebug {
		logname := time.Now().Format(time.RFC3339)
		f, err := os.Create(fmt.Sprintf(`%s_%s.log`, streamer.Username, logname))
		if err != nil {
			panic(err)
		}
		defer utils.DeferCloseReader(f)
		outs = append(outs, f)
	}
	slog.SetDefault(utils.NewSloglogger(logLevel, outs...))
	if config.Debug {
		slog.SetLogLoggerLevel(slog.LevelDebug)
	}
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	manager, err := manager.NewManager(ctx, config)
	if err != nil {
		log.Fatal(err)
	}
	writedBytes, err := manager.RecordStream(ctx, streamer)
	if err != nil && !utils.IsCancel(err) {
		slog.Error("RecordStream", slog.String("error", err.Error()))
	}
	slog.Info("file size", slog.String("size", utils.FormatFileSize(writedBytes)))
	slog.Info("process done", "streamer", streamer)
}
