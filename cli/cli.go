package cli

import (
	"flag"
	"fmt"
	"os"
	"runtime"

	"resize/service"
)

func Run(batch *service.Batch, args []string) {
	fs := flag.NewFlagSet("resize", flag.ExitOnError)
	var (
		root    = fs.String("path", "", "папка с JPG")
		maxEdge = fs.Int("max", 3840, "максимальная сторона, px")
		threads = fs.Int("threads", max(1, runtime.NumCPU()-2), "число потоков")
	)
	_ = fs.Parse(args)

	if *root == "" {
		fmt.Fprintln(os.Stderr, "не указана --path")
		fs.Usage()
		os.Exit(1)
	}

	// конфигурируем Batch
	batch.MaxEdge = *maxEdge
	batch.Workers = *threads
	var done, total int
	batch.OnProgress = func(d, t int) {
		done = d
		total = t
		fmt.Printf("\r[%d/%d] %.1f%%", done, total, float64(done)*100/float64(t))
	}

	if err := batch.Run(*root); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка:", err)
		os.Exit(1)
	}
	fmt.Println("\nГотово")
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
