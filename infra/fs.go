package infra

import (
	"io/fs"
	"path/filepath"
	"strings"
)

type FS struct{}

func (FS) List(root string) ([]string, error) {
	var out []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		e := strings.ToLower(filepath.Ext(path))
		if e == ".jpg" || e == ".jpeg" {
			out = append(out, path)
		}
		return nil
	})
	return out, err
}
