package cmd

import (
	"sync"
	"time"
)

// Idea is to create progress bar, that can be like
// a table with details like
//
// +----------+--------+-------+---------+--------+---------+
// |Date      |Time    |Action |File Name|Duration|Bandwidth|
// +----------+--------+-------+---------+--------+---------+
// |XX-XX-XXXX|hh:mm:ss|  $$$  |  ****   |   ^s   |  #kib/s |
// +----------+--------+-------+---------+--------+---------+
// |                                                        |
// +-------------------+-----------------+------------------+
// |Success:           | Fail:           | Total:           |
// +-------------------+-----------------+------------------+
//
// Print errors here
//
//

type BetterProgressBar struct {
	// Current number of bytes shared
	Current int64

	// How many were shared previously
	Previous int64

	// Stats for progress
	Success int64
	Fail    int64
	Total   int64

	// Start time for progress bar
	StartTime time.Time

	// Refresh Rate of progress bar
	RefreshRate time.Duration

	// finishOnce uses sync Once to close the bar
	finishOnce     sync.Once
	isFinishedChan chan struct{}
}
