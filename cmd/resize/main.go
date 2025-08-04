package main

import (
	"log/slog"
	"os"
	"resize/logutil"
	"runtime"

	"resize/cli"
	"resize/infra"
	"resize/service"
	"resize/ui"
)

func main() {

	// -- создаём логгер
	logger, closeFn, err := logutil.New(slog.LevelInfo, "resize.log")
	if err != nil {
		// если даже логгер не создался — пишем в ошибки и выходим
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("init logger", "err", err)
		os.Exit(1)
	}
	defer closeFn()

	// ядро-Batch создаём одно
	b := &service.Batch{
		Logger:  logger,
		Finder:  infra.FS{},
		Resizer: infra.Imaging{},
		Workers: runtime.NumCPU() - 2,
		MaxEdge: 3840,
	}

	// распознаём режим
	if len(os.Args) > 1 && os.Args[1] == "--cli" {
		cli.Run(b, os.Args[2:]) // всё после --cli отдаём CLI-флагам
		return
	}
	// по умолчанию GUI
	ui.Run(b)
}
