package ui

import (
	"fmt"
	"runtime"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"

	"resize/service"
)

func Run(batch *service.Batch) {
	a := app.NewWithID("app.resize")
	w := a.NewWindow("Resize JPGs")

	pathEntry := widget.NewEntry()

	pick := widget.NewButton("Select", func() {
		if p, err := zenity.SelectFile(zenity.Directory()); err == nil {
			pathEntry.SetText(p)
		}
	})

	thrEntry := widget.NewEntry()
	thrEntry.SetText(strconv.Itoa(max(1, runtime.NumCPU()-2)))

	maxEntry := widget.NewEntry()
	maxEntry.SetText("3840")

	bar := widget.NewProgressBar()

	runBtn := widget.NewButton("Run", func() {
		t, _ := strconv.Atoi(thrEntry.Text)
		m, _ := strconv.Atoi(maxEntry.Text)
		batch.Workers = t
		batch.MaxEdge = m
		batch.OnProgress = func(d, tot int) {
			fyne.DoAndWait(func() { bar.SetValue(float64(d) / float64(tot)) })
		}
		go func() {
			if err := batch.Run(pathEntry.Text); err != nil {
				fmt.Println("error:", err)
			}
		}()
	})

	w.SetContent(container.NewVBox(
		container.NewHBox(pathEntry, pick),
		container.NewHBox(widget.NewLabel("Threads:"), thrEntry),
		container.NewHBox(widget.NewLabel("Max px:"), maxEntry),
		bar,
		runBtn,
	))
	w.Resize(fyne.NewSize(420, 280))
	w.ShowAndRun()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
