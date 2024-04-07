package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

const (
	FPS = "10"
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
			err := os.MkdirAll(filepath.Join(frameDir, outputName), 0755)
			if err != nil {
				fmt.Println("Error creating output directory:", err)
			}
			cmd := exec.Command("ffmpeg", "-i", path, "-vf", "fps="+FPS, filepath.Join(frameDir, outputName, outputName+"_%04d.jpg"))
			err = cmd.Run()
			if err != nil {
				fmt.Printf("Error splitting video frames for %s: %v\n", info.Name(), err)
			} else {
				fmt.Printf("Video frames extracted successfully for %s\n", info.Name())
			}
		}
		return nil
	})
	if err != nil {
		fmt.Println("Error walking through directory:", err)
		return
	}
}
