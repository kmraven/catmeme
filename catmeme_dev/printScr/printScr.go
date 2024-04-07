package main

import (
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/errnoh/gocurse/curses"
)

func main() {
	// 初期化
	stdscr, err := curses.Initscr()
	if err != nil {
		fmt.Println("Curses init error:", err)
		return
	}
	defer curses.Endwin()
	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)
	curses.Noecho()
	// curses.Start_color()

	// 毎秒20ファイルずつ5秒間出力する
	fps := 20
	delay := time.Duration(float64(time.Second) / float64(fps))
	seconds := 5
	frameDir := "frames/cat_meme"
	for i := 0; i < fps*seconds; i++ {
		inputFilePath := filepath.Join(frameDir, fmt.Sprintf("cat_meme_%04d.jpg", i+1))
		if asciiArt, err := processImage(inputFilePath); err != nil {
			fmt.Printf("Error processing image %s: %v\n", inputFilePath, err)
		} else {
			stdscr.Clear()
			// stdscr.Move(0, 0)
			stdscr.Addstr(0, 0, asciiArt, 0)
			// stdscr.Addstr(0, 0, fmt.Sprintf("%v, %v\n", sec, i), 0)
			stdscr.Refresh()
			// 次の秒へ進むまで待機
			time.Sleep(delay)
		}
	}
}

// 画像ファイルに特定の処理を行う関数（ここでは例としてコピーのみ）
func processImage(inputPath string) (string, error) {
	flags := aic_package.DefaultFlags()
	// flags.Dimensions = []int{100, 30}
	flags.Colored = false
	// flags.Width = 1000
	// flags.Height = 600
	// flags.SaveTxtPath = "."
	// flags.CustomMap = " .-=+#@"
	// flags.SaveBackgroundColor = [4]int{50, 50, 50, 100}
	// flags.SaveGifPath = "ascii_arts"
	// Conversion for an image
	asciiArt, err := aic_package.Convert(inputPath, flags)

	return asciiArt, err
}
