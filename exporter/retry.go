package exporter

import (
	"encoding/json"
	"fmt"
	"os"
)

type BufferEntry struct {
	Message string  `json:"msg"`
	Payload Payload `json:"payload"`
}

func Resend() {
	retryInProgress = true

	messages := ReadLogFile()

	for _, obj := range messages {
		payload := obj.Payload
		PostLoggedData(payload.NodeId, payload.NodeName, payload.Value, payload.Timestamp, payload.LogName, payload.Server)
	}
	bufferSize = 0
	retryInProgress = false

}

func ReadLogFile() []BufferEntry {

	dat, err := os.ReadFile("tmp/logs/buffer.json")

	if err != nil {
		Logs.Error(fmt.Sprint(err))
	}

	jsonString := string(dat)
	jsonString = jsonString[0 : len(jsonString)-1]

	jsonString = "[" + jsonString + "]"

	var jsonArr []BufferEntry

	if err := json.Unmarshal([]byte(jsonString), &jsonArr); err != nil {
		Logs.Error(fmt.Sprint(err))
	}

	os.Truncate("tmp/logs/buffer.json", 0)

	return jsonArr
}
