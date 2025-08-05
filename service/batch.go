package service

import (
	"log/slog"
	"sync"
)

type (
	FileFinder interface {
		List(root string) ([]string, error)
	}
	Resizer interface {
		Resize(path string, max int) error
	}
	Progress func(done, total int)

	Batch struct {
		Logger     *slog.Logger
		Finder     FileFinder
		Resizer    Resizer
		Workers    int
		MaxEdge    int
		OnProgress Progress
	}
)

func (b Batch) Run(root string) error {
	files, err := b.Finder.List(root)
	if err != nil {
		b.Logger.Error("listing jpg files failed", "err", err)
		return err
	}

	total := len(files)
	if total == 0 {
		b.Logger.Info("no jpg files in directory")
		return nil
	}

	sem := make(chan struct{}, b.Workers)
	wg := sync.WaitGroup{}

	var (
		done int
		mu   sync.Mutex
	)

	for _, f := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {

			defer func() {
				<-sem
				wg.Done()
			}()

			if err := b.Resizer.Resize(path, b.MaxEdge); err != nil {
				b.Logger.Error("can't resize file", "path", path, "err", err)
				return
			}

			mu.Lock()
			done++
			if b.OnProgress != nil {
				b.OnProgress(done, total)
			}
			mu.Unlock()
		}(f)
	}
	wg.Wait()
	return nil
}
