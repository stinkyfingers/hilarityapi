package main

import (
	"github.com/stinkyfingers/hilarity/server"
	"log"
)

// POSTMAN - don't forget Origin header = http://localhost

func main() {
	s := server.NewServer(8888)
	if err := s.Run(); err != nil {
		log.Fatal(err)
	}
}
