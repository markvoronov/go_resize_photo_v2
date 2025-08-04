package infra_test

import (
	"image"
	"image/color"
	"image/jpeg"
	"os"
	"path/filepath"
	"testing"

	"resize/infra"
)

func createJPEG(t *testing.T, dir string, name string, w, h int) string {
	path := filepath.Join(dir, name)
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	// заполним чем-нибудь, чтобы Encode не оптимизировал
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.Set(x, y, color.RGBA{R: 255, A: 255})
		}
	}
	f, err := os.Create(path)
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if err := jpeg.Encode(f, img, nil); err != nil {
		t.Fatalf("encode: %v", err)
	}
	f.Close()
	return path
}

func decodeDims(t *testing.T, path string) (int, int) {
	f, err := os.Open(path)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	defer f.Close()
	img, err := jpeg.Decode(f)
	if err != nil {
		t.Fatalf("decode: %v", err)
	}
	b := img.Bounds()
	return b.Dx(), b.Dy()
}

func TestImaging_Resize(t *testing.T) {
	tmp := t.TempDir()
	path := createJPEG(t, tmp, "big.jpg", 2000, 1000)

	if err := (infra.Imaging{}).Resize(path, 800); err != nil {
		t.Fatalf("resize error: %v", err)
	}

	w, h := decodeDims(t, path)
	if w != 800 {
		t.Errorf("width = %d, want 800", w)
	}
	if h != 400 {
		t.Errorf("height = %d, want 400", h)
	}
}

func TestImaging_Resize_NoChange(t *testing.T) {
	tmp := t.TempDir()
	path := createJPEG(t, tmp, "small.jpg", 400, 300)

	if err := (infra.Imaging{}).Resize(path, 800); err != nil {
		t.Fatalf("resize error: %v", err)
	}

	w, h := decodeDims(t, path)
	if w != 400 || h != 300 {
		t.Errorf("image unexpectedly changed size to %dx%d", w, h)
	}
}
