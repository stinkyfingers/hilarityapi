package server

import (
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"net/http"
)

type Server struct {
	Port string
}

func NewServer(port int) *Server {
	return &Server{
		Port: fmt.Sprintf(":%d", port),
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/game", websocket.Handler(s.gameHandler))
	return http.ListenAndServe(s.Port, mux)
}

func (s *Server) gameHandler(ws *websocket.Conn) {
	type Data struct {
		Message string
	}
	for {
		var d Data
		err := websocket.JSON.Receive(ws, &d)
		if err != nil {
			log.Print(err)
			break
		}
		log.Print(d)
	}
}
