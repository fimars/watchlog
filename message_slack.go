package watchdog

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type SlackMessage struct {
	Text string `json:"text"`
}

// sends a message to a Slack channel
func SendToSlackChannel(message string) {
	url := "https://hooks.slack.com/services/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"

	msg := SlackMessage{Text: message}
	msgBytes, _ := json.Marshal(msg)
	body := bytes.NewReader(msgBytes)

	_, err := http.Post(url, "application/json;charset=utf-8", body)
	if err != nil {
		fmt.Println(err)
	}
}
