package main

import (
	"time"
	"log"
	"io/ioutil"
	"strconv"
	"os"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	filelock "github.com/zbiljic/go-filelock"
)

var id_filename string = "/Users/aca123321/Desktop/Personal/Fly io distriibuted systems challenges/utils/id.txt"

func main() {
	file, err := os.OpenFile("/Users/aca123321/Desktop/Personal/Fly io distriibuted systems challenges/logs/unique_id_generation_logs.txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)

	lock, err := filelock.New(id_filename)
	if err != nil {
		log.Fatal("Error creating lock:", err)
		return
	}

	n := maelstrom.NewNode()
	n.Handle("generate", func(msg maelstrom.Message) error {
		body := make(map[string]any)

		for true {
			isLockObtained, _ := lock.TryLock()
			if !isLockObtained {
				time.Sleep(10 * time.Millisecond)
			} else {
				// read from file
				data, err := ioutil.ReadFile(id_filename)
				if err != nil {
					log.Fatal("File reading error", err)
					return err
				}

				// get last used id as integer
				last_used_id, err := strconv.Atoi(string(data))
				log.Println("LAST USED ID: ", last_used_id)

				// populate body
				body["type"] = "generate_ok"
				body["id"] = last_used_id + 1
				log.Println("Sending ", body["id"])

				// write to file
				err = ioutil.WriteFile(id_filename, []byte(strconv.Itoa(last_used_id+1)), 0644)
				if err != nil {
						log.Fatal("File writing error", err)
						return err
				}

				lock.Unlock()
				return n.Reply(msg, body)
			}
		}
		return n.Reply(msg, body)
	})

	if err := n.Run(); err != nil {
		log.Fatal(err)
	}
}
