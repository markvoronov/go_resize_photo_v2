package main

import (
	"fmt"
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/ncruces/zenity"
	"github.com/nfnt/resize"
	"image"
	"image/jpeg"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"sync"
)

func main() {

	a := app.NewWithID("mark.resize")
	w := a.NewWindow("Folder Selector")

	// Текстовое поле для отображения и редактирования пути
	pathEntry := widget.NewEntry()
	pathEntry.SetPlaceHolder("Selected folder path")

	// Кнопка для выбора папки
	selectFolderBtn := widget.NewButton("Select folder", func() {
		folder, err := zenity.SelectFile(
			zenity.Title("Select a Folder"),
			zenity.Directory(),
		)
		if err != nil {
			if err == zenity.ErrCanceled {
				fmt.Println("User cancelled")
				return
			}
			fmt.Println("Error selecting folder:", err)
			return
		}
		pathEntry.SetText(folder)
	})

	// Поле для ввода количества потоков
	threadsEntry := widget.NewEntry()
	NumCPU := runtime.NumCPU()
	NumFlow := NumCPU
	if NumCPU >= 3 {
		NumFlow = NumCPU - 2
	}

	threadsEntry.SetText(strconv.Itoa(NumFlow)) // Значение по умолчанию — число процессоров
	threadsEntry.SetPlaceHolder("Number of threads")
	threadsLabel := widget.NewLabel("Threads:")
	threadsRow := container.NewHBox(threadsLabel, threadsEntry)

	// Поле для ввода максимального разрешения
	maxResEntry := widget.NewEntry()
	maxResEntry.SetText("3840") // Значение по умолчанию — 3840
	maxResEntry.SetPlaceHolder("Max resolution")
	maxResLabel := widget.NewLabel("Max resolution:")
	maxResRow := container.NewHBox(maxResLabel, maxResEntry)

	// Прогресс-бар
	progressBar := widget.NewProgressBar()
	progressBar.SetValue(0)

	// Кнопка для обработки
	resizeBtn := widget.NewButton("Resize image", func() {
		// Получаем количество потоков
		numThreads, err := strconv.Atoi(threadsEntry.Text)
		if err != nil || numThreads < 1 {
			numThreads = runtime.NumCPU() // Если ввод некорректен, используем число процессоров
		}
		if numThreads > runtime.NumCPU() {
			numThreads = runtime.NumCPU() // Ограничиваем максимумом
		}

		// Получаем максимальное разрешение
		maxRes, err := strconv.Atoi(maxResEntry.Text)
		if err != nil || maxRes < 1 {
			maxRes = 3840 // Если ввод некорректен, используем 3840
		}

		fmt.Println("Processing path:", pathEntry.Text, "with", numThreads, "threads, max resolution:", maxRes)
		go selectAmdResize(pathEntry.Text, numThreads, maxRes, progressBar)
	})

	// Контейнер для поля ввода и кнопки выбора папки (горизонтальное размещение)
	folderRow := container.NewVBox(
		pathEntry, // Поле с заданной шириной
		selectFolderBtn,
	)

	// Главный контейнер (вертикальное размещение строк)
	content := container.NewVBox(
		folderRow,   // Поле и кнопка выбора папки
		threadsRow,  // Поле для числа потоков
		maxResRow,   // Поле для максимального разрешения
		progressBar, // Прогресс-бар
		resizeBtn,   // Кнопка "Resize image"
	)

	// Устанавливаем контейнер в качестве содержимого окна
	w.SetContent(content)

	w.Resize(fyne.NewSize(500, 400))
	w.ShowAndRun()
}

func selectAmdResize(selectedPath string, numThreads int, maxRes int, progressBar *widget.ProgressBar) {
	files, err := getImagesFromPath(selectedPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	totalFiles := len(files)
	if totalFiles == 0 {
		return
	}

	// Канал для файлов
	filesToProcess := make(chan string, 1000)
	// Счётчик обработанных файлов
	processedCount := 0
	var countMutex sync.Mutex

	// Заполняем канал файлами
	wg1 := &sync.WaitGroup{}
	wg1.Add(1)
	go func() {
		defer wg1.Done()
		for _, fileName := range files {
			filesToProcess <- fileName
		}
		close(filesToProcess)
	}()

	// Запускаем рабочие горутины
	wg := &sync.WaitGroup{}
	for i := 0; i < numThreads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			worker(filesToProcess, &countMutex, &processedCount, totalFiles, maxRes, progressBar)
		}()
	}

	// Ждём завершения
	wg1.Wait()
	wg.Wait()
}

func getImagesFromPath(selectedPath string) ([]string, error) {
	images := []string{}
	err := filepath.WalkDir(selectedPath, func(path string, d os.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if !d.IsDir() {
			if strings.HasSuffix(strings.ToLower(path), ".jpg") || strings.HasSuffix(strings.ToLower(path), ".jpeg") {
				images = append(images, path)
			}
		}
		return nil
	})
	if err != nil {
		return images, fmt.Errorf("error during filepath.WalkDir: %v", err)
	}
	if len(images) == 0 {
		return images, fmt.Errorf("no images found")
	}
	return images, nil
}

func worker(ToProcess <-chan string, countMutex *sync.Mutex, processedCount *int, totalFiles int, maxRes int, progressBar *widget.ProgressBar) {
	for value := range ToProcess {
		err := resizeImage(value, maxRes)
		if err != nil {
			fmt.Println("Error processing", value, ":", err)
			continue
		}
		// Обновляем прогресс в главном потоке UI
		countMutex.Lock()
		*processedCount++
		progress := float64(*processedCount) / float64(totalFiles)
		fyne.DoAndWait(func() {
			progressBar.SetValue(progress)
		})
		countMutex.Unlock()
	}
}

func resizeImage(inputPath string, maxRes int) error {
	inputFile, err := os.Open(inputPath)
	if err != nil {
		log.Println("Error opening file:", err)
		return err
	}
	defer inputFile.Close()

	img, err := jpeg.Decode(inputFile)
	if err != nil {
		log.Println("Error decoding image:", err)
		return err
	}

	newWidth, newHeight := getResizedDimensions(img, maxRes)
	if newWidth+newHeight == 0 {
		return nil // Изображение меньше maxRes
	}

	newImg := resize.Resize(newWidth, newHeight, img, resize.Lanczos3)

	outputFile, err := os.Create(inputPath)
	if err != nil {
		log.Println("Error creating file:", err)
		return err
	}
	defer outputFile.Close()

	err = jpeg.Encode(outputFile, newImg, nil)
	if err != nil {
		log.Println("Error encoding image:", err)
		return err
	}
	return nil
}

func getResizedDimensions(img image.Image, maxLength int) (uint, uint) {
	bounds := img.Bounds()
	width := bounds.Dx()
	height := bounds.Dy()

	if width <= maxLength && height <= maxLength {
		return 0, 0 // Не требуется изменение размера
	}

	var newWidth, newHeight float64
	if width > height {
		newWidth = float64(maxLength)
		newHeight = float64(height) * (float64(maxLength) / float64(width))
	} else {
		newHeight = float64(maxLength)
		newWidth = float64(width) * (float64(maxLength) / float64(height))
	}

	return uint(newWidth), uint(newHeight)
}
