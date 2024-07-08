package main

import (
	"flag"
	"io"
	"log"
	"strings"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/hpcloud/tail"
	"github.com/hpcloud/tail/ratelimiter"
)

var logfilePath = flag.String("p", "/temp", "log files path")

var (
	tailingFiles = []string{}
)

func main() {
	flag.Parse()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	done := make(chan bool)

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}

				if event.Op&fsnotify.Write == fsnotify.Write {
					included := false
					for _, file := range tailingFiles {
						if file == event.Name {
							included = true
							break
						}
					}
					if !included {
						log.Println("New File" + event.Name + " is created.")
						tailingFiles = append(tailingFiles, event.Name)
						err := tailLog(event.Name)
						if err != nil {
							log.Println("读取文件内容错误:", err)
							continue
						}
					}
				}
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Println("WatchLogs Error:", err)
			}

			time.Sleep(1000 * time.Millisecond)
		}
	}()

	watchPath := *logfilePath // 替换为你要监听的日志文件路径
	err = watcher.Add(watchPath)
	if err != nil {
		log.Fatal(err)
	}

	<-done
}

func tailLog(filename string) error {
	log.Println(filename)

	seek := &tail.SeekInfo{
		Offset: 0,
		Whence: io.SeekEnd,
	}
	config := tail.Config{
		Location:  seek,
		Follow:    true,
		Poll:      true,
		ReOpen:    true,
		MustExist: true,
		RateLimiter: ratelimiter.NewLeakyBucket(
			4*1000, // 1000 characters per Millisecond
			time.Millisecond,
		),
	}
	t, err := tail.TailFile(filename, config)

	if err != nil {
		return err
	}

	for lines := range t.Lines {
		// TODO: logic here
		if strings.Contains(lines.Text, "ERROR") {
			log.Println(lines.Text)
		}
	}

	return nil
}
