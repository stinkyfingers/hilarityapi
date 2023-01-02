package main

import (
	"log"
	"os"
	"strconv"

	"github.com/stinkyfingers/hilarity/server"
)

// POSTMAN - don't forget Origin header = http://localhost

var port = 8888

func main() {
	if portStr := os.Getenv("PORT"); portStr != "" {
		var err error
		port, err = strconv.Atoi(portStr)
		if err != nil {
			log.Fatal(err)
		}
	}
	
	s, err := server.NewServer(port)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
