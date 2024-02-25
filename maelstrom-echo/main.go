package main

import (
    "encoding/json"
    "log"
		"path/filepath"
		"os"
		"runtime"
		"fmt"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func initLogging() {
	_, curfilename, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Error: Unable to get the current file path")
		return
	}
	logfile_path := filepath.Join(filepath.Dir(curfilename), "../logs/echo_log.txt")
	file, err := os.OpenFile(logfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func main() {
	initLogging()

	n := maelstrom.NewNode()
	n.Handle("echo", func(msg maelstrom.Message) error {
			// Unmarshal the message body as an loosely-typed map.
			var body map[string]any
			if err := json.Unmarshal(msg.Body, &body); err != nil {
					return err
			}

			// Update the message type to return back.
			body["type"] = "echo_ok"

			// Echo the original message back with the updated message type.
			return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
    log.Fatal(err)
	}
}
