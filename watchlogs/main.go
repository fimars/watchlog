// https://github.com/fsnotify/fsnotify/blob/main/cmd/fsnotify/watch.go
package main

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

var usage = `
watchlogs is a logs monitor application.
This command serves as an example and debugging tool.


Commands:

    watch [paths]  Watch the paths for changes and track the errors.
`[1:]

func exit(format string, a ...interface{}) {
	fmt.Fprintf(os.Stderr, filepath.Base(os.Args[0])+": "+format+"\n", a...)
	fmt.Print("\n" + usage)
	os.Exit(1)
}

func help() {
	fmt.Printf("%s [command] [arguments]\n\n", filepath.Base(os.Args[0]))
	fmt.Print(usage)
	os.Exit(0)
}

// 时间前缀很有用
func printTime(s string, args ...interface{}) {
	fmt.Printf(time.Now().Format("15:04:05.0000")+" "+s+"\n", args...)
}

func main() {
	if len(os.Args) == 1 {
		help()
	}
	// Always show help if -h[elp] appears anywhere before we do anything else.
	for _, f := range os.Args[1:] {
		switch f {
		case "help", "-h", "-help", "--help":
			help()
		}
	}

	cmd, args := os.Args[1], os.Args[2:]
	switch cmd {
	default:
		exit("unknown command: %q", cmd)
	case "watch":
		watch(args...)
	}
}
