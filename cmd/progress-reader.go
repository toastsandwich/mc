// Copyright (c) 2015-2022 MinIO, Inc.
//
// This file is part of MinIO Object Storage stack
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cmd

import (
	"io"
	"runtime"
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

func showLastProgressBar(pg ProgressReader) {
	if pgbar, ok := pg.(*ProgressBar); ok {
		pgbar.Finish()
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
