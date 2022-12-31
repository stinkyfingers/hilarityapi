package server

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/stinkyfingers/hilarity/game"
	"github.com/stinkyfingers/hilarity/storage"
	"github.com/stinkyfingers/hilarity/user"
	"golang.org/x/net/websocket"
)

// POST
func (s *Server) NewGame(w http.ResponseWriter, r *http.Request) {
	name := r.URL.Query().Get("game")
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	if ok, err := s.Storage.NameExists(name); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	} else if ok {
		s.JSONErr(w, fmt.Errorf("game with that name exists"), http.StatusConflict)
		return
	}
	roundsStr := r.URL.Query().Get("rounds")
	totalRounds := game.DefaultTotalRounds
	if roundsStr != "" {
		var err error
		totalRounds, err = strconv.Atoi(roundsStr)
		if err != nil {
			s.JSONErr(w, err, http.StatusBadRequest)
			return
		}
	}
	questions, err := s.Storage.GetQuestions()
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	g, err := game.NewGame(name, totalRounds, questions)
	if err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	g.Join(u)
	err = s.Storage.SaveGame(g)
	s.JSON(w, g, err)
}

// POST
func (s *Server) JoinGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	game := r.URL.Query().Get("game")
	if game == "" {
		s.JSONErr(w, fmt.Errorf("no game specified"), http.StatusBadRequest)
		return
	}
	g, err := s.Storage.GetGame(game)
	if err != nil {
		s.JSONErr(w, err, http.StatusBadRequest)
		return
	}
	g.Join(u)
	err = s.Storage.SaveGame(g)
	s.JSON(w, g, err)
}

// DELETE
func (s *Server) LeaveGame(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	var u user.User
	if err := json.NewDecoder(r.Body).Decode(&u); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	game := r.URL.Query().Get("game")
	if game == "" {
		s.JSONErr(w, fmt.Errorf("no game specified"), http.StatusBadRequest)
		return
	}
	g, err := s.Storage.GetGame(game)
	if err != nil {
		s.JSONErr(w, err, http.StatusBadRequest)
		return
	}
	g.Leave(u)
	err = s.Storage.SaveGame(g)
	s.JSON(w, g, err)
}

// GET
func (s *Server) ListGames(w http.ResponseWriter, r *http.Request) {
	if err := s.Storage.CleanUpGames(); err != nil {
		s.JSONErr(w, err, http.StatusInternalServerError)
		return
	}
	games, err := s.Storage.ListGames()
	s.JSON(w, games, err)
}

// GET
func (s *Server) GetGame(w http.ResponseWriter, r *http.Request) {
	game := r.URL.Query().Get("game")
	g, err := s.Storage.GetGame(game)
	s.JSON(w, g, err)
}

// WS
func (s *Server) PlayGame(ws *websocket.Conn) {
	gameName := ws.Request().URL.Query().Get("game")
	s.Disseminator.Join(gameName, ws)
	defer func() {
		s.Disseminator.Leave(gameName, ws)
		ws.Close()
	}()
	for {
		var gp game.GamePlay
		err := websocket.JSON.Receive(ws, &gp)
		if err != nil {
			s.WSErr(ws, err)
			return
		}

		g, err := storage.SubmitPlay(s.Storage, gp)
		if err != nil {
			s.WSErr(ws, err)
			return
		}
		err = s.Disseminator.Disseminate(gameName, g)
	}
}
