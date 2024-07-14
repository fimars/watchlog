package main

import (
	"bufio"
	"fmt"
	"os"
	"testing"
	"time"
)

var (
	testdataPath = "../../testdata/"
	logfile      = "logfile"
)

func TestMain(m *testing.M) {
	code := m.Run()
	os.Exit(code)
}
func TestSplitLogs(t *testing.T) {
	os.RemoveAll(testdataPath)
	os.Mkdir(testdataPath, os.ModePerm)
	time.Sleep(300 * time.Millisecond)

	done := make(chan bool)

	count := 0
	expectSum := 1000

	debugFunc := func(line string) {
		count = count + 1
		if count == expectSum {
			done <- true
		}
		fmt.Printf(`Debug: %v, %v
`, line, count)
	}
	Watcher.DebugLineReader = &debugFunc

	// prepare data
	for i := 0; i < 1000; i++ {
		appendToFile(testdataPath+logfile, fmt.Sprint(`Line `, i, `, time: `, time.Now().Format("2006-01-02 15:04:05")))
		time.Sleep(15 * time.Millisecond)
	}

	go watch(testdataPath)

	for i := 0; i < expectSum; i++ {
		appendToFile(testdataPath+logfile, fmt.Sprint(`Line `, i, `, time: `, time.Now().Format("2006-01-02 15:04:05")))
		time.Sleep(5 * time.Millisecond)

		if (i+1)%100 == 0 {
			time.Sleep(300 * time.Millisecond)
			os.Rename(
				testdataPath+logfile,
				testdataPath+fmt.Sprint(logfile, time.Now().Format("2006-01-02 15:04:05")),
			)
			time.Sleep(300 * time.Millisecond)
		}
	}
	fmt.Println("done")
	<-done
	os.RemoveAll(testdataPath)
}

func appendToFile(filePath, text string) error {
	// 以追加模式打开文件，如果文件不存在则创建
	file, err := os.OpenFile(filePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	// 创建一个带缓冲的写入器
	writer := bufio.NewWriter(file)

	// 写入内容
	_, err = writer.WriteString(text + "\n")
	if err != nil {
		return err
	}

	// 刷新缓冲区，确保内容写入文件
	err = writer.Flush()
	if err != nil {
		return err
	}

	return nil
}
