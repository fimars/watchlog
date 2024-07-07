package main

/// https://medium.com/@arunprabhu.1/tailing-a-file-in-golang-72944204f22b

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"time"
)

func main() {
	file, err := os.Open("./logs/logfile.log")
	if err != nil {
		return
	}
	defer file.Close()
	reader := bufio.NewReader(file)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				// without this sleep you would hogg the CPU
				time.Sleep(500 * time.Millisecond)
				// truncated ?
				truncated, errTruncated := isTruncated(file)
				if errTruncated != nil {
					break
				}
				if truncated {
					// seek from start
					_, errSeekStart := file.Seek(0, io.SeekStart)
					if errSeekStart != nil {
						break
					}
				}
				continue
			}
			break
		}

		if strings.Contains(line, "ERROR") {
			log.Println("-----------------")
			log.Println(line)
		}
	}
}

func isTruncated(file *os.File) (bool, error) {
	// current read position in a file
	currentPos, err := file.Seek(0, io.SeekCurrent)
	if err != nil {
		return false, err
	}
	// file stat to get the size
	fileInfo, err := file.Stat()
	if err != nil {
		return false, err
	}
	return currentPos > fileInfo.Size(), nil
}
