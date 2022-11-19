package main

import (
	"log"

	"github.com/stinkyfingers/hilarity/server"
)

// POSTMAN - don't forget Origin header = http://localhost

func main() {
	s, err := server.NewServer(8888)
	if err != nil {
		log.Fatal(err)
	}

	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
