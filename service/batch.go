package service

import "sync"

type (
	FileFinder interface {
		List(root string) ([]string, error)
	}
	Resizer interface {
		Resize(path string, max int) error
	}
	Progress func(done, total int)

	Batch struct {
		Finder     FileFinder
		Resizer    Resizer
		Workers    int
		MaxEdge    int
		OnProgress Progress
	}
)

func (b Batch) Run(root string) error {
	files, err := b.Finder.List(root)
	if err != nil || len(files) == 0 {
		return err
	}

	total := len(files)
	sem := make(chan struct{}, b.Workers)
	wg := sync.WaitGroup{}

	done := 0
	mu := sync.Mutex{}

	for _, f := range files {
		wg.Add(1)
		sem <- struct{}{}
		go func(path string) {
			defer func() { <-sem; wg.Done() }()
			_ = b.Resizer.Resize(path, b.MaxEdge)

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
