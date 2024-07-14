package watchdog

import (
	"fmt"
	"strings"
)

func ReadErrorLine(line string) {
	if strings.Contains(line, "ERROR") {
		msg := "[WatchDogs] " + line
		SendToSlackChannel(msg)
	}

	if strings.Contains(line, "queue.INFO") {
		fmt.Println(
			fmt.Errorf("[WatchDogs *Queue line] " + line),
		)
	}

	if strings.Contains(line, "app.NOTICE") {
		fmt.Println("[WatchDogs *Notice line] " + line)
	}

	if strings.Contains(line, "adjustTier") {
		msg := "[WatchDogs someone adjustTier.] " + line
		fmt.Println(msg)
		SendToSlackChannel(msg)
	}
}
