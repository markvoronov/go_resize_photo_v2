package main

import (
	"log/slog"
	"os"
	"resize/cli"
	"resize/infra"
	"resize/logutil"
	"resize/service"
)

func main3() {

	// -- создаём логгер
	logger, closeFn, err := logutil.New(slog.LevelInfo, "resize.log")
	if err != nil {
		// если даже логгер не создался — пишем в ошибки и выходим
		slog.New(slog.NewTextHandler(os.Stderr, nil)).Error("init logger", "err", err)
		os.Exit(1)
	}
	defer closeFn()

	b := &service.Batch{
		Logger:  logger,
		Finder:  infra.FS{},
		Resizer: infra.Imaging{},
		Workers: 4,
		MaxEdge: 3840,
	}
	cli.Run(b, os.Args[1:])
}
