package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

const (
	FPS = "10"
)

// 動画を画像に変換して保存する
func main() {
	// 処理対象のディレクトリ
	frameDir := "frames/cat_meme"

	// 指定されたディレクトリ内の動画ファイルを取得する
	err := filepath.Walk(frameDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".jpg") {
			inputFilePath := filepath.Join(frameDir, info.Name())
			if _, err := processImage(inputFilePath); err != nil {
				fmt.Printf("Error processing image %s: %v\n", inputFilePath, err)
			} else {

			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}
}

// 画像ファイルに特定の処理を行う関数（ここでは例としてコピーのみ）
func processImage(inputPath string) (string, error) {
	flags := aic_package.DefaultFlags()
	// flags.Dimensions = []int{100, 30}
	flags.Colored = true
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
