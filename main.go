package main

import (
	"embed"
	"flag"
	"fmt"
	"image"
	"math/rand"
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
	FPS                   = 20
	FRAME_DIR             = "frames"
	DEFAULT_TIME_DURATION = 3
)

//go:embed frames/*
var files embed.FS

// 元素材を整える (trim, frame数減らす)
// 素材増やす

func main() {
	var (
		coloredFlag  = flag.Bool("c", false, "colored flag")
		timeDuration = flag.Int("t", DEFAULT_TIME_DURATION, "time duration (sec)")
	)
	flag.Parse()

	// get size of screen
	w, h, err := getTerminalSize()
	if err != nil {
		panic(err)
	}

	// setup tview
	app := tview.NewApplication()
	screen := tview.NewTextView().SetScrollable(false).SetTextAlign(tview.AlignCenter)
	if *coloredFlag {
		screen = screen.SetDynamicColors(true)
	}
	screen.SetChangedFunc(func() {
		app.Draw()
	})
	ignoreKeys := func(event *tcell.EventKey) *tcell.EventKey {
		return nil
	}
	app.SetInputCapture(ignoreKeys)

	// setup timer
	ticker := time.NewTicker(time.Second / time.Duration(FPS))
	defer ticker.Stop()
	loopTimer := time.NewTimer(time.Duration(*timeDuration) * time.Second)
	defer loopTimer.Stop()

	// select kind of cat meme
	entries, err := files.ReadDir(FRAME_DIR)
	if err != nil {
		panic(err)
	}
	var dirs []os.DirEntry
	for _, entry := range entries {
		if entry.IsDir() {
			dirs = append(dirs, entry)
		}
	}
	contentName := dirs[rand.Intn(len(dirs))].Name()

	// count filenum
	filenum := 0
	if f, err := files.ReadDir(filepath.Join(FRAME_DIR, contentName)); err != nil {
		panic(err)
	} else {
		filenum = len(f)
	}
	filecounter := 1

	go func() {
		for {
			select {
			case <-ticker.C:
				inputFilePath := filepath.Join(FRAME_DIR, contentName, fmt.Sprintf("frame_%04d.jpg", filecounter))
				if asciiArt, err := processImage(inputFilePath, w, h, *coloredFlag); err != nil {
					fmt.Printf("[Error processingImage func] %s: %v\n", inputFilePath, err)
				} else {
					if *coloredFlag {
						asciiArt = tview.TranslateANSI(asciiArt)
					}
					screen.SetText(asciiArt)
				}
				filecounter++
				if filecounter > filenum {
					filecounter = 1
				}
			case <-loopTimer.C:
				app.Stop()
				return
			}
		}
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

// clone of aic_package, for use embed files
func processImage(inputPath string, w, h int, coloredFlag bool) (string, error) {
	// do not use 'font color'! because temporary 'flattenAscii' func can not work with it.
	flags := aic_package.DefaultFlags()
	flags.Dimensions = []int{w, h}
	if coloredFlag {
		flags.Colored = true
	}

	localImg, err := files.Open(inputPath)
	if err != nil {
		return "", fmt.Errorf("unable to open file: %v", err)
	}
	defer localImg.Close()

	imData, _, err := image.Decode(localImg)
	if err != nil {
		return "", fmt.Errorf("can't decode %v: %v", inputPath, err)
	}

	imgSet, err := imgManip.ConvertToAsciiPixels(
		imData, flags.Dimensions, flags.Width, flags.Height,
		flags.FlipX, flags.FlipY, flags.Full, flags.Braille, flags.Dither,
	)
	if err != nil {
		return "", err
	}

	var asciiSet [][]imgManip.AsciiChar

	if flags.Braille {
		asciiSet, err = imgManip.ConvertToBrailleChars(
			imgSet, flags.Negative, flags.Colored, flags.Grayscale,
			flags.CharBackgroundColor, flags.FontColor, flags.Threshold,
		)
	} else {
		asciiSet, err = imgManip.ConvertToAsciiChars(
			imgSet, flags.Negative, flags.Colored, flags.Grayscale, flags.Complex,
			flags.CharBackgroundColor, flags.CustomMap, flags.FontColor,
		)
	}
	if err != nil {
		return "", err
	}

	ascii := flattenAscii(asciiSet, flags.Colored || flags.Grayscale, false)
	result := strings.Join(ascii, "\n")

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
