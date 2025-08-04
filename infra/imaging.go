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
	in.Close() // закрываем сразу!
	if err != nil {
		return err
	}
	// Определим размеры нового изображения
	newWidth, newHeight := dims(img, max)
	//if newWidth+newHeight == 0 {
	//	return nil // изображение меньше, чем max
	//}
	dst := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	// пишем во временный файл
	tmp := path + ".tmp"
	out, err := os.Create(tmp)
	if err != nil {
		return err
	}
	if err := jpeg.Encode(out, dst, nil); err != nil {
		out.Close()
		return err
	}
	out.Close() // тоже закрываем до Rename!

	// замена оригинального изображения обработанным
	return os.Rename(tmp, path)
}

func dims(im image.Image, maxLength int) (uint, uint) {
	bounds := im.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	newWidth := 0
	newHeight := 0

	if width > height {
		newWidth = width
	} else {
		newHeight = height
	}

	if newHeight > maxLength {
		newHeight = maxLength
	}

	if newWidth > maxLength {
		newWidth = maxLength
	}

	return uint(newWidth), uint(newHeight)

}
