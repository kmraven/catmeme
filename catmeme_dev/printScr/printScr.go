package main

import (
	"fmt"
	"os"
	"path/filepath"
	"syscall"
	"time"
	"unsafe"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

// sigintを無効にする
// 色付き出力

func main() {
	// get size of screen
	w, h, err := getTerminalSize()
	if err != nil {
		panic(err)
	}

	// setup tview
	app := tview.NewApplication()
	screen := tview.NewTextView().SetScrollable(false).SetTextAlign(tview.AlignCenter).SetDynamicColors(true)
	screen.SetChangedFunc(func() {
		app.Draw()
	})
	ignoreKeys := func(event *tcell.EventKey) *tcell.EventKey {
		return nil
	}
	app.SetInputCapture(ignoreKeys)

	fps := 20
	delay := time.Duration(float64(time.Second) / float64(fps))
	seconds := 5
	frameDir := "frames/cat_meme"
	go func() {
		for i := 0; i < fps*seconds; i++ {
			inputFilePath := filepath.Join(frameDir, fmt.Sprintf("cat_meme_%04d.jpg", i+1))
			if asciiArt, err := processImage(inputFilePath, w-1, h-1); err != nil {
				fmt.Printf("Error processing image %s: %v\n", inputFilePath, err)
			} else {
				tmp := tview.TranslateANSI(asciiArt)
				screen.SetText(tmp)
			}
			time.Sleep(delay)
		}
		app.Stop()
	}()

	if err := app.SetRoot(screen, true).Run(); err != nil {
		panic(err)
	}
}

// 画像ファイルに特定の処理を行う関数（ここでは例としてコピーのみ）
func processImage(inputPath string, w, h int) (string, error) {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{w, h}
	// flags.Colored = true
	asciiArt, err := aic_package.Convert(inputPath, flags)
	// os.WriteFile("./txt/tmp.txt", []byte(asciiArt), 0644)

	return asciiArt, err
}

func getTerminalSize() (width, height int, err error) {
	var dimensions [4]uint16
	_, _, errno := syscall.Syscall(
		syscall.SYS_IOCTL,
		os.Stdout.Fd(),
		syscall.TIOCGWINSZ,
		uintptr(unsafe.Pointer(&dimensions)),
	)
	if errno != 0 {
		err = errno
		return
	}
	width = int(dimensions[1])
	height = int(dimensions[0])
	return
}
