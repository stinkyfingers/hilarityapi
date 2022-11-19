package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/stinkyfingers/hilarity/storage"
	"golang.org/x/net/websocket"
)

type Server struct {
	Port    string
	Storage storage.Storage
}

func NewServer(port int) (*Server, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	var store storage.Storage
	if os.Getenv("STORAGE") == "mem" {
		store = storage.NewInMemory()
	} else {
		store, err = storage.NewFile(filepath.Join(pwd, "gamestorage"))
		if err != nil {
			return nil, err
		}
	}
	return &Server{
		Port:    fmt.Sprintf(":%d", port),
		Storage: store,
	}, nil
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/game", websocket.Handler(s.PlayGame))
	mux.HandleFunc("/game/new", s.NewGame)
	mux.HandleFunc("/game/join", s.JoinGame)
	mux.HandleFunc("/games/list", s.ListGames)
	return http.ListenAndServe(s.Port, mux)
}

func (s *Server) JSON(w http.ResponseWriter, data interface{}, dataErr error) {
	if dataErr != nil {
		s.JSONErr(w, dataErr, http.StatusInternalServerError)
		return
	}
	err := json.NewEncoder(w).Encode(data)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
	}
}

func (s *Server) JSONErr(w http.ResponseWriter, dataErr error, code int) {
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": dataErr.Error(),
		"code":  code,
	})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}

func (s *Server) WSErr(w io.Writer, dataErr error) {
	err := json.NewEncoder(w).Encode(map[string]interface{}{
		"error": dataErr.Error(),
	})
	if err != nil {
		w.Write([]byte(err.Error()))
	}
}
