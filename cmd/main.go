package main

import (
	"resize/infra"
	"resize/service"
	"resize/ui"
)

func main() {
	ui.Run(&service.Batch{
		Finder:  infra.FS{},
		Resizer: infra.Imaging{},
		Workers: 4,
		MaxEdge: 3840,
	})
}
