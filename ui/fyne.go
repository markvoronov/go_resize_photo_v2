package ui

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/validation"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	"resize/service" // если пакет batch — кастомный
	"runtime"
	"strconv"
)

func Run(batch *service.Batch) {
	a := app.NewWithID("app.resize")
	w := a.NewWindow("Resize JPGs")

	maxThreads := runtime.NumCPU()
	defaultThreads := max(1, maxThreads-2)

	// Путь
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Select folder with JPG/JPEG files")

	selectBtn := widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
		if p, err := zenity.SelectFile(zenity.Directory()); err == nil {
			pathEntry.SetText(p)
		}
	})
	selectBtn.Importance = widget.LowImportance

	// Потоки
	thrEntry := widget.NewEntry()
	thrEntry.SetText(strconv.Itoa(defaultThreads))
	thrEntry.Validator = validation.NewRegexp(`^\d+$`, "Number only")

	// Максимальный размер
	maxEntry := widget.NewEntry()
	maxEntry.SetText("3840")
	maxEntry.Validator = validation.NewRegexp(`^\d+$`, "Number only")

	// Прогресс
	bar := widget.NewProgressBar()
	bar.SetValue(0)

	// Кнопка запуска
	runBtn := widget.NewButton("▶ Run", func() {
		t, _ := strconv.Atoi(thrEntry.Text)
		m, _ := strconv.Atoi(maxEntry.Text)
		batch.Workers = t
		batch.MaxEdge = m
		batch.OnProgress = func(d, tot int) {
			fyne.DoAndWait(func() {
				bar.SetValue(float64(d) / float64(tot))
			})
		}
		go func() {
			if err := batch.Run(pathEntry.Text); err != nil {
				fmt.Println("Error:", err)
			}
		}()
	})

	// Сборка интерфейса
	form := container.NewVBox(
		container.NewBorder(nil, nil, nil, selectBtn, pathEntry),
		widget.NewLabel("Supported formats: .jpg, .jpeg"),
		container.NewGridWithColumns(2,
			container.NewVBox(widget.NewLabel(fmt.Sprintf("Threads (max %d):", maxThreads)), thrEntry),
			container.NewVBox(widget.NewLabel("Max image dimension (px):"), maxEntry),
		),
		bar,
		runBtn,
	)

	w.SetContent(container.NewPadded(form))
	w.Resize(fyne.NewSize(440, 270))
	w.ShowAndRun()
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
