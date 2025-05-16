package cmd

import (
	"io"
	"runtime"

	"github.com/minio/pkg/v3/console"
)

// Progress - an interface which describes current amount
// of data written.
type Progress interface {
	Get() int64
	SetTotal(int64)
}

// ProgressReader can be used to update the progress of
// an on-going transfer progress.
type ProgressReader interface {
	io.Reader
	Progress
}

func showLastProgressBar(pg ProgressReader, e error) {
	if e != nil {
		// We only erase a line if we are displaying a progress bar
		if !globalQuiet && !globalJSON {
			console.Eraseline()
		}
		return
	}
	if accntReader, ok := pg.(*accounter); ok {
		printMsg(accntReader.Stat())
	}
}

// cursorAnimate - returns a animated rune through read channel for every read.
func cursorAnimate() <-chan string {
	cursorCh := make(chan string)
	var cursors []string

	switch runtime.GOOS {
	case "linux":
		// cursors = "➩➪➫➬➭➮➯➱"
		// cursors = "▁▃▄▅▆▇█▇▆▅▄▃"
		cursors = []string{"◐", "◓", "◑", "◒"}
		// cursors = "←↖↑↗→↘↓↙"
		// cursors = "◴◷◶◵"
		// cursors = "◰◳◲◱"
		// cursors = "⣾⣽⣻⢿⡿⣟⣯⣷"
	case "darwin":
		cursors = []string{"◐", "◓", "◑", "◒"}
	default:
		cursors = []string{"|", "/", "-", "\\"}
	}
	go func() {
		for {
			for _, cursor := range cursors {
				cursorCh <- cursor
			}
		}
	}()
	return cursorCh
}
