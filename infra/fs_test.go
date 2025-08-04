package infra_test

import (
	"os"
	"path/filepath"
	"testing"

	"resize/infra"
)

func TestFS_List(t *testing.T) {
	tmp := t.TempDir()

	// создаём тестовую структуру файлов
	files := []string{
		"root.jpg",
		"root.jpeg",
		"sub/inner.jpg",
		"skip.png",
	}
	for _, f := range files {
		full := filepath.Join(tmp, f)
		if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(full, []byte("x"), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
	}

	got, err := (infra.FS{}).List(tmp)
	if err != nil {
		t.Fatalf("fs list error: %v", err)
	}

	want := map[string]bool{
		filepath.Join(tmp, "root.jpg"):      true,
		filepath.Join(tmp, "root.jpeg"):     true,
		filepath.Join(tmp, "sub/inner.jpg"): true,
	}
	if len(got) != len(want) {
		t.Fatalf("got %d paths, want %d: %#v", len(got), len(want), got)
	}
	for _, p := range got {
		if !want[p] {
			t.Errorf("unexpected path %s", p)
		}
	}
}
