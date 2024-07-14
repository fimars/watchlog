package watchdog

import (
	"fmt"
	"io"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/nxadm/tail"
	"github.com/nxadm/tail/ratelimiter"
)

func Tail(
	file string,
	readFunc func(line string),
	op fsnotify.Op,
) {
	seek := &tail.SeekInfo{
		Offset: 0,
		Whence: io.SeekEnd,
	}

	if (op & fsnotify.Create) != 0 {
		seek.Whence = io.SeekStart
	}

	config := tail.Config{
		Location: seek,
		Follow:   true,
		Poll:     true,
		ReOpen:   false,
		// MustExist: true,
		RateLimiter: ratelimiter.NewLeakyBucket(
			4*1000, // 1000 characters per Millisecond
			time.Millisecond,
		),
	}
	t, err := tail.TailFile(file, config)

	if err != nil {
		panic(err)
	}

	fmt.Println(len(t.Lines))

	for lines := range t.Lines {
		readFunc(lines.Text)
	}

	<-make(chan struct{}) // Block forever
}
