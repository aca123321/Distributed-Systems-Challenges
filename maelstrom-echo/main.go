package main

import (
    "encoding/json"
    "log"
		"os"

    maelstrom "github.com/jepsen-io/maelstrom/demo/go"
)

func main() {
	n := maelstrom.NewNode()
	file, err := os.OpenFile("/Users/aca123321/Desktop/Personal/Fly io distriibuted systems challenges/logs/echo_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	n.Handle("echo", func(msg maelstrom.Message) error {
			// Unmarshal the message body as an loosely-typed map.
			var body map[string]any
			if err := json.Unmarshal(msg.Body, &body); err != nil {
					return err
			}

			log.Println("REPLYING from node 1")
			// Update the message type to return back.
			body["type"] = "echo_ok"

			// Echo the original message back with the updated message type.
			return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
    log.Fatal(err)
	}
}
