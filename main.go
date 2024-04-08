package main

import (
	"embed"
	"fmt"
	"image"
	"os"
	"path/filepath"
	"strings"
	"syscall"
	"time"
	"unsafe"

	"github.com/TheZoraiz/ascii-image-converter/aic_package"
	imgManip "github.com/TheZoraiz/ascii-image-converter/image_manipulation"
	"github.com/gdamore/tcell/v2"
	"github.com/rivo/tview"
)

const (
	fps          = 20
	frameDir     = "frames"
	tmpVideoName = "cat_meme"
)

//go:embed frames/*
var files embed.FS

// 元素材を整える (trim, fps調整, frame数減らす)
// 存在するフレームを指定秒数ループする処理に変更
// 素材増やす, ランダム処理追加

func main() {
	// get size of screen
	w, h, err := getTerminalSize()
	if err != nil {
		panic(err)
	}

	// setup tview
	app := tview.NewApplication()
	screen := tview.NewTextView().SetScrollable(false).SetTextAlign(tview.AlignCenter)
	// screen = screen.SetDynamicColors(true) // if color
	screen.SetChangedFunc(func() {
		app.Draw()
	})
	ignoreKeys := func(event *tcell.EventKey) *tcell.EventKey {
		return nil
	}
	app.SetInputCapture(ignoreKeys)

	delay := time.Duration(float64(time.Second) / float64(fps))
	seconds := 5
	go func() {
		for i := 0; i < fps*seconds; i++ {
			inputFilePath := filepath.Join(frameDir, tmpVideoName, fmt.Sprintf("cat_meme_%04d.jpg", i+1))
			if asciiArt, err := processImage(inputFilePath, w, h); err != nil {
				fmt.Printf("[Error processingImage func] %s: %v\n", inputFilePath, err)
			} else {
				// asciiArt := tview.TranslateANSI(asciiArt) // if color
				screen.SetText(asciiArt)
			}
			time.Sleep(delay)
		}
		app.Stop()
	}()

	if err := app.SetRoot(screen, true).Run(); err != nil {
		panic(err)
	}
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

func processImage(inputPath string, w, h int) (string, error) {
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{w, h}
	// flags.Colored = true // if color
	// do not use font color! because temporary 'flattenAscii' func can not work with it.

	localImg, err := files.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}
	defer localImg.Close()

	imData, _, err := image.Decode(localImg)
	if err != nil {
		return "", fmt.Errorf("can't decode %v: %v", inputPath, err)
	}

	imgSet, err := imgManip.ConvertToAsciiPixels(imData, flags.Dimensions, flags.Width, flags.Height, flags.FlipX, flags.FlipY, flags.Full, flags.Braille, flags.Dither)
	if err != nil {
		return "", err
	}

	var asciiSet [][]imgManip.AsciiChar

	if flags.Braille {
		asciiSet, err = imgManip.ConvertToBrailleChars(imgSet, flags.Negative, flags.Colored, flags.Grayscale, flags.CharBackgroundColor, flags.FontColor, flags.Threshold)
	} else {
		asciiSet, err = imgManip.ConvertToAsciiChars(imgSet, flags.Negative, flags.Colored, flags.Grayscale, flags.Complex, flags.CharBackgroundColor, flags.CustomMap, flags.FontColor)
	}
	if err != nil {
		return "", err
	}

	ascii := flattenAscii(asciiSet, flags.Colored || flags.Grayscale, false)
	result := strings.Join(ascii, "\n")

	// os.WriteFile("./txt/tmp.txt", []byte(result), 0644)
	return result, nil
}

// flattenAscii flattens a two-dimensional grid of ascii characters into a one dimension
// of lines of ascii
func flattenAscii(asciiSet [][]imgManip.AsciiChar, colored, toSaveTxt bool) []string {
	var ascii []string

	for _, line := range asciiSet {
		var tempAscii string

		for _, char := range line {
			if toSaveTxt {
				tempAscii += char.Simple
				continue
			}

			if colored {
				tempAscii += char.OriginalColor
				// } else if fontColor != [3]int{255, 255, 255} {
				// 	tempAscii += char.SetColor
			} else {
				tempAscii += char.Simple
			}
		}

		ascii = append(ascii, tempAscii)
	}

	return ascii
}
