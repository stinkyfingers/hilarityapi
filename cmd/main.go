package main

import (
	"bufio"
	"encoding/json"
	"log"
	"os"

	"golang.org/x/net/websocket"
)

func main() {
	ws, err := websocket.Dial("ws://localhost:8888/game", "", "http://localhost")
	if err != nil {
		log.Fatal(err)
	}

	type Data struct {
		Message string
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		log.Println(scanner.Text())
		d := Data{scanner.Text()}
		j, err := json.Marshal(d)
		if err != nil {
			log.Fatal(err)
		}
		ws.Write(j)
	}
}
