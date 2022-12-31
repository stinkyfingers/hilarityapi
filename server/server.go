package server

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"

	"github.com/rs/cors"
	"github.com/stinkyfingers/hilarity/storage"
	"golang.org/x/net/websocket"
)

type Server struct {
	Port         string
	Storage      storage.Storage
	Disseminator *Disseminator
}

type Disseminator struct {
	conns map[string][]*websocket.Conn
	mutex sync.RWMutex
}

func NewServer(port int) (*Server, error) {
	pwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	// TODO S3
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
		Port:         fmt.Sprintf(":%d", port),
		Storage:      store,
		Disseminator: NewDisseminator(),
	}, nil
}

func NewDisseminator() *Disseminator {
	return &Disseminator{
		conns: make(map[string][]*websocket.Conn),
		mutex: sync.RWMutex{},
	}
}

func (s *Server) Run() error {
	mux := http.NewServeMux()
	mux.Handle("/game/play", websocket.Handler(s.PlayGame))
	mux.HandleFunc("/game/new", s.NewGame)
	mux.HandleFunc("/game/join", s.JoinGame)
	mux.HandleFunc("/game/leave", s.LeaveGame)
	mux.HandleFunc("/games/list", s.ListGames)
	mux.HandleFunc("/game/details", s.GetGame)
	return http.ListenAndServe(s.Port, cors.Default().Handler(mux))
}

func (s *Server) JSON(w http.ResponseWriter, data interface{}, dataErr error) {
	if dataErr != nil {
		s.JSONErr(w, dataErr, http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
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

func (d *Disseminator) Disseminate(gameName string, msg interface{}) error {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	conns, ok := d.conns[gameName]
	if !ok {
		return fmt.Errorf("no connections found")
	}
	for _, conn := range conns {
		if err := websocket.JSON.Send(conn, msg); err != nil {
			return err
		}
	}
	return nil
}

func (d *Disseminator) Join(gameName string, conn *websocket.Conn) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	conns, ok := d.conns[gameName]
	if !ok {
		conns = []*websocket.Conn{}
	}
	conns = append(conns, conn)
	d.conns[gameName] = conns
}

func (d *Disseminator) Leave(gameName string, conn *websocket.Conn) {
	d.mutex.Lock()
	defer d.mutex.Unlock()
	conns, ok := d.conns[gameName]
	if !ok {
		return
	}
	for i, c := range conns {
		if c == conn {
			conns = append(conns[:i], conns[i+1:]...)
			break
		}
	}
	d.conns[gameName] = conns
}
