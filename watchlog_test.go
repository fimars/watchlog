package watchdog

import (
	"fmt"
	"math/rand"
	"os"
	"strings"
	"testing"

	"github.com/fsnotify/fsnotify"
	"github.com/stretchr/testify/assert"
)

const (
	fileName         = "./testdata/logfile.log"
	errorProbability = 0.05
)

func TestErrorCounter(t *testing.T) {
	done := make(chan bool)
	count := 0

	go Tail(fileName, func(line string) {
		if strings.Contains(line, "ERROR") {
			count += 1
		}
		if count == 50 {
			done <- true
		}
	}, fsnotify.Write)
	go generateLogs(50)
	<-done

	os.Remove(fileName)
	assert.Equal(t, count, 50)
}

func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func generateLogs(errorLimit uint) {
	file, _ := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	defer file.Close()

	var count uint = 0

	for {
		randomString := generateRandomString(100)
		if rand.Float64() < errorProbability {
			randomString += " ERROR"
			count += 1
			fmt.Println("Error count:", count)
		}
		_, err := file.WriteString(randomString + "\n")
		if err != nil {
			return
		}
		if errorLimit == count {
			return
		}
	}

}
