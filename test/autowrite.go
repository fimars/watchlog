package main

import (
	"fmt"
	"math/rand"
	"os"
	"time"
)

const (
	fileName         = "./logs/logfile.log"
	errorProbability = 0.05
)

func main() {
	rand.Seed(time.Now().UnixNano())

	file, err := os.OpenFile(fileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()

	for {
		randomString := generateRandomString(10)
		if rand.Float64() < errorProbability {
			randomString += " ERROR"
		}

		_, err := file.WriteString(randomString + "\n")
		if err != nil {
			fmt.Println("Error writing to file:", err)
			return
		}

		time.Sleep(10 * time.Microsecond)
	}
}

func generateRandomString(length int) string {
	const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, length)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}
