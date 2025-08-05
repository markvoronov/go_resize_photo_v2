package logutil

import (
	"io"
	"log/slog"
	"os"
)

// New Возвращает логгер, который дублирует вывод в файл и stdout.
func New(level slog.Level, path string) (*slog.Logger, func() error, error) {
	// открываем (или создаём) файл. Будем туда дописывать логи
	f, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0o644)
	if err != nil {
		return nil, nil, err
	}

	w := io.MultiWriter(os.Stdout, f)

	h := slog.NewTextHandler(w, &slog.HandlerOptions{
		Level: level,
		// красивый формат времени
		ReplaceAttr: func(_ []string, a slog.Attr) slog.Attr {
			if a.Key == slog.TimeKey {
				a.Value = slog.StringValue(
					a.Value.Time().Format("2006-01-02 15:04:05"))
			}
			return a
		},
	})

	return slog.New(h), f.Close, nil
}
