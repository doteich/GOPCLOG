package exporter

import (
	"encoding/json"
	"fmt"
	"os"
)

type bufferEntry struct {
	Message string  `json:"msg"`
	Payload Payload `json:"payload"`
}

func ReadLogFile() {

	var entry []bufferEntry
	dat, err := os.ReadFile("tmp/logs/buffer.json")

	if err != nil {
		Logs.Error(fmt.Sprint(err))
	}

	if err := json.Unmarshal(dat, &entry); err != nil {

	}

	// fmt.Println(entry)
}
