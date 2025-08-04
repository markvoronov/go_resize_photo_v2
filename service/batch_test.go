package service_test

import (
	"sync"
	"sync/atomic"
	"testing"

	"resize/service"
)

// ------------------- заглушки ----------------------------------------------

type stubFinder struct{ paths []string }

func (s stubFinder) List(root string) ([]string, error) { return s.paths, nil }

type stubResizer struct {
	mu    sync.Mutex
	calls []string
}

func (s *stubResizer) Resize(path string, _ int) error {
	s.mu.Lock()
	s.calls = append(s.calls, path)
	s.mu.Unlock()
	return nil
}

// ---------------------------------------------------------------------------

func TestBatch_Run(t *testing.T) {
	paths := []string{"1.jpg", "2.jpg", "3.jpg", "4.jpg"}

	var (
		cbCalls  int32 // сколько раз вызвали колбэк
		lastDone int32 // какой done прислали в последний раз
	)
	batch := service.Batch{
		Finder:  stubFinder{paths},
		Resizer: &stubResizer{},
		Workers: 2,
		MaxEdge: 800,
		OnProgress: func(done, total int) {
			atomic.AddInt32(&cbCalls, 1) // просто считаем вызовы
			atomic.StoreInt32(&lastDone, int32(done))
			if total != len(paths) {
				t.Errorf("total=%d, want %d", total, len(paths))
			}
		},
	}

	if err := batch.Run("/dummy"); err != nil {
		t.Fatalf("batch run error: %v", err)
	}

	if int(cbCalls) != len(paths) {
		t.Errorf("OnProgress calls = %d, want %d", cbCalls, len(paths))
	}
	if int(lastDone) != len(paths) {
		t.Errorf("last done = %d, want %d", lastDone, len(paths))
	}
}
