package game

import (
	"errors"
	"log"
	"sync"

	"github.com/stinkyfingers/hilarity/user"
)

type Game struct {
	ID           int                  `json:"id"`
	Name         string               `json:"name"`
	Users        map[string]user.User `json:"users"` // username:user
	PastRounds   []Round              `json:"pastRounds"`
	CurrentRound Round                `json:"currentRound"`
	mutex        sync.RWMutex         `json:"-"`
}

type Round struct {
	Question string           `json:"question"`
	Plays    map[string]Play  `json:"plays"`   // username:play
	Guesses  map[string]Guess `json:"guesses"` // username:guess
}

type Play struct {
	Answers [3]string `json:"answers"`
}

type Guess struct {
	Responses map[string]string // guessed user:actual user
}

type GamePlay struct {
	GameName string    `json:"gameName"`
	User     user.User `json:"user"`
	Play     *Play     `json:"play"`
	Guess    *Guess    `json:"guess"`
}

var ErrNilPlay = errors.New("play is empty")
var ErrNilGuess = errors.New("guess is empty")

func getQuestion() (string, error) { // TODO
	return "top 3 ice cream flavors", nil
}

func NewGame(name string) (*Game, error) {
	round, err := NewRound()
	if err != nil {
		return nil, err
	}

	return &Game{
		Name:         name,
		CurrentRound: *round,
		Users:        make(map[string]user.User),
		mutex:        sync.RWMutex{},
	}, nil
}

func NewRound() (*Round, error) {
	question, err := getQuestion()
	if err != nil {
		return nil, err
	}
	return &Round{
		Question: question,
		Plays:    make(map[string]Play),
		Guesses:  make(map[string]Guess),
	}, nil
}

func (g *Game) Join(u user.User) {
	g.mutex.Lock()
	g.Users[u.Name] = u
	g.mutex.Unlock()
}

func (g *Game) Leave(u user.User) {
	g.mutex.Lock()
	delete(g.Users, u.Name)
	g.mutex.Unlock()
}

/*
	is play or guess round?
	play:
		if len(play) == len(users):
			switch to guess round
		else:
			update play for user
	guess:
		if len(guess) == len(users):
			switch to NEW round
			save to storage
		else:
			update guess for user
*/

func (gp *GamePlay) MakePlay(g *Game) error {
	if len(g.CurrentRound.Plays) < len(g.Users) {
		return gp.PlayRound(g)
	}
	if len(g.CurrentRound.Guesses) < len(g.Users) {
		return gp.GuessRound(g)
	}
	// next round & play
	round, err := NewRound()
	if err != nil {
		return err
	}
	g.PastRounds = append(g.PastRounds, g.CurrentRound)
	g.CurrentRound = *round
	return gp.PlayRound(g)
}

func (gp *GamePlay) PlayRound(g *Game) error {
	log.Println("PLAY")
	if gp.Play == nil {
		return ErrNilPlay
	}
	g.CurrentRound.Plays[gp.User.Name] = *gp.Play
	return nil
}

func (gp *GamePlay) GuessRound(g *Game) error {
	log.Println("GUESS")
	if gp.Guess == nil {
		return ErrNilGuess
	}
	g.CurrentRound.Guesses[gp.User.Name] = *gp.Guess
	return nil
}
