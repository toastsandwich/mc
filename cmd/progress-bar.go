package cmd

import (
	"fmt"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/cheggaaa/pb"
	"github.com/fatih/color"
	"github.com/minio/pkg/v3/console"
)

const REFRESH_RATE = time.Millisecond * 125

// [dd time] [action] [dur / bandwidth] [src -> dst]
type ProgressBar struct {
	total   int64
	current int64

	src    string // src path
	dst    string // dst path
	action string // action performed

	complete string

	mu     sync.Mutex
	finish sync.Once

	finishCh chan struct{}

	start       time.Time
	refreshRate time.Duration

	p pb.ProgressBar
}

func NewProgressBar(total int64) *ProgressBar {
	return &ProgressBar{
		start:       time.Now(),
		finishCh:    make(chan struct{}),
		total:       total,
		refreshRate: REFRESH_RATE,
	}
}

func (pg *ProgressBar) refresher() {
	pg.Update()
	for {
		select {
		case <-pg.finishCh:
			return
		case <-time.After(pg.refreshRate):
			pg.Update()
		}
	}
}

func (pg *ProgressBar) Update() {
	// total := pg.GetTotal()
	pg.write()
}

func (pg *ProgressBar) Start() *ProgressBar {
	// pg.start = time.Now()
	// fmt.Printf("\rstarting operation %s", pg.action)
	go pg.refresher()
	return pg
}

func (pg *ProgressBar) once() {
	taken := time.Since(pg.start)
	console.Eraseline()
	conc := fmt.Sprintf("\r[ done ] [ %s ]", taken.String())
	color.Green(conc)
	close(pg.finishCh)
}

func (pg *ProgressBar) Finish() {
	pg.finish.Do(pg.once)
}

func (pg *ProgressBar) SetTotal(v int64) {
	atomic.StoreInt64(&pg.total, v)
}

func (pg *ProgressBar) SetCurrent(v int64) {
	atomic.StoreInt64(&pg.current, v)
}

func (pg *ProgressBar) Get() int64 {
	return atomic.LoadInt64(&pg.current)
}

func (pg *ProgressBar) Add(n int64) int64 {
	return atomic.AddInt64(&pg.current, n)
}

func (pg *ProgressBar) GetTotal() int64 {
	return atomic.LoadInt64(&pg.total)
}

func (pg *ProgressBar) Read(p []byte) (int, error) {
	defer func() {
		if c, t := pg.Get(), pg.GetTotal(); t > 0 && c > t {
			pg.SetCurrent(t)
		}
	}()
	n := len(p)
	pg.Add(int64(n))
	return n, nil
}

func (pg *ProgressBar) write() {
	pg.mu.Lock()

	console.Eraseline()
	fmt.Printf("\r[ %s ] [ %s > %s ] [%d] [%d]", pg.action, pg.src, pg.dst, pg.current, pg.total)

	pg.mu.Unlock()
}

func (pg *ProgressBar) SetCaption(s string) *ProgressBar {
	pg.mu.Lock()
	sp := strings.Split(s, "\x00")
	if len(sp) == 3 {
		pg.setSrc(fixateBarCaption(sp[0], 18))
		pg.setDst(fixateBarCaption(sp[1], 18))
		pg.setAction(sp[2])
	} else {
		pg.setSrc("?")
		pg.setDst("?")
		pg.setAction("?")
	}
	pg.mu.Unlock()
	return pg
}

// SetSrc sets source path
func (pg *ProgressBar) setSrc(src string) {
	pg.src = src
}

// SetDst sets destination path
func (pg *ProgressBar) setDst(dst string) {
	pg.dst = dst
}

// SetAction sets action being performed
func (pg *ProgressBar) setAction(action string) {
	pg.action = action
}

// fixateBarCaption - fancify bar caption based on the terminal width.
func fixateBarCaption(caption string, width int) string {
	switch {
	case len(caption) > width:
		// Trim caption to fit within the screen
		trimSize := len(caption) - width + 3
		if trimSize < len(caption) {
			caption = "..." + caption[trimSize:]
		}
	case len(caption) < width:
		caption += strings.Repeat(" ", width-len(caption))
	}
	return caption
}
