package infra

import (
	"image"
	"image/jpeg"
	"os"

	"github.com/nfnt/resize"
)

type Imaging struct{}

func (Imaging) Resize(path string, max int) error {
	// ---------- читаем исходник ----------
	in, err := os.Open(path)
	if err != nil {
		return err
	}
	img, err := jpeg.Decode(in)
	in.Close() // <--- закрываем сразу!
	if err != nil {
		return err
	}

	w, h := dims(img, max)
	if w == 0 { // уже меньше max
		return nil
	}
	dst := resize.Resize(w, h, img, resize.Lanczos3)

	// ---------- пишем во временный файл ----------
	tmp := path + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(out, dst, nil); err != nil {
		out.Close()
		return err
	}
	out.Close() // <--- тоже закрываем до Rename!

	// ---------- атомарная подмена ----------
	return os.Rename(tmp, path)
}

func dims(im image.Image, max int) (uint, uint) {
	b := im.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= max && h <= max {
		return 0, 0
	}
	if w >= h {
		return uint(max), uint(float64(h) * float64(max) / float64(w))
	}
	return uint(float64(w) * float64(max) / float64(h)), uint(max)
}
