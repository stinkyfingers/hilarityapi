package server

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/stinkyfingers/hilarity/game"
	"github.com/stinkyfingers/hilarity/storage"
	"github.com/stinkyfingers/hilarity/user"
	"golang.org/x/net/websocket"
)

// POST
func (s *Server) NewGame(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("name")
	defer r.Body.Close()
	if ok, err := s.Storage.NameExists(name); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	} else if ok {
		s.JSONErr(w, fmt.Errorf("game with that name exists"), http.StatusConflict)
		return
	}
	g, err := game.NewGame(name)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	err = s.Storage.SaveGame(g)
	s.JSON(w, g, err)
}

// POST
func (s *Server) JoinGame(w http.ResponseWriter, r *http.Request) {
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	game := r.URL.Query().Get("game")
	g, err := s.Storage.GetGame(game)
	if err != nil {
		s.JSONErr(w, err, http.StatusBadRequest)
		return
	}
	g.Users[u.Name] = u
	err = s.Storage.SaveGame(g)
	s.JSON(w, g, err)
}

// GET
func (s *Server) ListGames(w http.ResponseWriter, r *http.Request) {
	games, err := s.Storage.ListGames()
	s.JSON(w, games, err)
}

// WS
func (s *Server) PlayGame(ws *websocket.Conn) {
	var gp game.GamePlay
	for {
		err := websocket.JSON.Receive(ws, &gp)
		if err != nil {
			s.WSErr(ws, err)
		}
		g, err := storage.SubmitPlay(s.Storage, gp)
		if err != nil {
			s.WSErr(ws, err)
		}
		err = websocket.JSON.Send(ws, g)
		if err != nil {
			s.WSErr(ws, err)
		}
	}
}
