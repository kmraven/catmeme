package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"
	"time"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	"github.com/gdamore/tcell"
)

func main() {
	// 初期化
	s, err := tcell.NewScreen()
	if err != nil {
		log.Fatalf("%+v", err)
	}
	if err := s.Init(); err != nil {
		log.Fatalf("%+v", err)
	}
	defer s.Fini()
	w, h := s.Size()

	sigint := make(chan os.Signal, 1)
	signal.Notify(sigint, syscall.SIGINT)

	// 毎秒20ファイルずつ5秒間出力する
	fps := 20
	delay := time.Duration(float64(time.Second) / float64(fps))
	seconds := 2
	frameDir := "frames/cat_meme"
	for i := 0; i < fps*seconds; i++ {
		inputFilePath := filepath.Join(frameDir, fmt.Sprintf("cat_meme_%04d.jpg", i+1))
		if asciiArt, err := processImage(inputFilePath, w-1, h-1); err != nil {
			fmt.Printf("Error processing image %s: %v\n", inputFilePath, err)
		} else {
			s.Clear()
			drawText(s, 0, 0, w, h, tcell.StyleDefault, asciiArt)
			s.Show()
			time.Sleep(delay)
		}
	}
}

func drawText(s tcell.Screen, x1, y1, x2, y2 int, style tcell.Style, text string) {
	row := y1
	col := x1
	for _, r := range []rune(text) {
		s.SetContent(col, row, r, nil, style)
		col++
		if col >= x2 {
			row++
			col = x1
		}
		if row > y2 {
			break
		}
	}
}

// 画像ファイルに特定の処理を行う関数（ここでは例としてコピーのみ）
func processImage(inputPath string, w, h int) (string, error) {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{w, h}
	// flags.Colored = true
	// flags.SaveTxtPath = "."
	// flags.CustomMap = " .-=+#@"
	// flags.SaveBackgroundColor = [4]int{}
	// flags.SaveGifPath = "ascii_arts"
	asciiArt, err := aic_package.Convert(inputPath, flags)

	return asciiArt, err
}
