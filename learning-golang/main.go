package main

import (
	"log"
	"io/ioutil"
	"strconv"
)

var id_filename string = "/Users/aca123321/Desktop/Personal/Fly io distriibuted systems challenges/utils/id.txt"

func main() {
	for i:=0;i<10;i++ {
		data, err := ioutil.ReadFile(id_filename)
		if err != nil {
			log.Fatal("File reading error", err)
			return
		}

		last_used_id, err := strconv.Atoi(string(data))
		log.Println("last used id: ", last_used_id)

		err = ioutil.WriteFile(id_filename, []byte(strconv.Itoa(last_used_id+1)), 0644)
    if err != nil {
        log.Fatal("File writing error", err)
				return
    }
	}
}
