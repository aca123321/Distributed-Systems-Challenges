package main

import (
	"time"
	"log"
	"io/ioutil"
	"strconv"
	"os"
	"path/filepath"
	"runtime"
	"fmt"

	maelstrom "github.com/jepsen-io/maelstrom/demo/go"
	filelock "github.com/zbiljic/go-filelock"
)

var curfilename string
var ok bool

func initLogging() {
	_, curfilename, _, ok = runtime.Caller(0)
	if !ok {
		fmt.Println("Error: Unable to get the current file path")
		return
	}
	logfile_path := filepath.Join(filepath.Dir(curfilename), "../logs/unique_id_generation_log.txt")
	file, err := os.OpenFile(logfile_path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
}

func main() {
	initLogging()
	id_filename := filepath.Join(filepath.Dir(curfilename), "../utils/id.txt")

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
