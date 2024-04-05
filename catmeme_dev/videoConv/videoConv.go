package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
)

// 動画を画像に変換して保存する
func main() {
	// 処理対象のディレクトリ
	inputDir := "videos"

	// 出力ディレクトリ
	frameDir := "frames"

	// ディレクトリが存在しない場合は作成する
	err := os.MkdirAll(frameDir, 0755)
	if err != nil {
		fmt.Println("Error creating output directory:", err)
		return
	}

	// 指定されたディレクトリ内の動画ファイルを取得する
	err = filepath.Walk(inputDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mp4") {
			// ffmpegを使って動画をフレームに分割し、画像として保存する
			outputName := strings.TrimSuffix(info.Name(), filepath.Ext(info.Name()))
			// ディレクトリが存在しない場合は作成する
			// err := os.MkdirAll(filepath.Join(frameDir, outputName), 0755)
			if err != nil {
				fmt.Println("Error creating output directory:", err)
			}
			cmd := exec.Command("ffmpeg", "-i", path, "-r", "10", filepath.Join(frameDir, outputName+".gif"))
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error splitting video frames for %s: %v\n", info.Name(), err)
			} else {
				fmt.Printf("Video frames extracted successfully for %s\n", info.Name())
			}
			inputFilePath := filepath.Join(frameDir, outputName+".gif")
			if err := processImage(inputFilePath); err != nil {
				fmt.Printf("Error processing image %s: %v\n", inputFilePath, err)
			} else {
				fmt.Printf("Image processed successfully: %s\n", inputFilePath)
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}
}

// 画像ファイルかどうかをチェックする関数
func isImageFile(filename string) bool {
	ext := filepath.Ext(filename)
	return ext == ".jpg" || ext == ".jpeg" || ext == ".png" || ext == ".gif" || ext == ".bmp"
}

// 画像ファイルに特定の処理を行う関数（ここでは例としてコピーのみ）
func processImage(inputPath string) error {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{100, 30}
	flags.Colored = true
	flags.Width = 1000
	flags.Height = 600
	// flags.SaveTxtPath = "."
	flags.CustomMap = " .-=+#@"
	// flags.SaveBackgroundColor = [4]int{50, 50, 50, 100}
	flags.SaveGifPath = "ascii_arts"
	// Conversion for an image
	asciiArt, err := aic_package.Convert(inputPath, flags)
	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("%v\n", asciiArt)

	return err
}
