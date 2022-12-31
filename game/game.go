package game

import (
	"errors"
	"math/rand"
	"sync"
	"time"

	"github.com/stinkyfingers/hilarity/user"
)

type Game struct {
	ID           int                  `json:"id"`
	Name         string               `json:"name"`
	Users        map[string]user.User `json:"users"` // username:user
	PastRounds   []*Round             `json:"pastRounds"`
	CurrentRound *Round               `json:"currentRound"`
	Questions    []string             `json:"questions"`
	TotalRounds  int                  `json:"totalRounds"`
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
	Responses map[string]string `json:"responses"` // guessed user:actual user
}

type GamePlay struct {
	GameName string    `json:"gameName"`
	User     user.User `json:"user"`
	Play     *Play     `json:"play"`
	Guess    *Guess    `json:"guess"`
}

var ErrNilPlay = errors.New("play is empty")
var ErrNilGuess = errors.New("guess is empty")
var ErrTooManyRounds = errors.New("games are limited to 10 rounds")
var ErrNoQuestions = errors.New("no unique questions left")

const DefaultTotalRounds int = 6
const MaxTotalRounds int = 10

func NewGame(name string, totalRounds int, allQuestions []string) (*Game, error) {
	if totalRounds > MaxTotalRounds {
		return nil, ErrTooManyRounds
	}

	questions, err := getQuestions(allQuestions, totalRounds)
	if err != nil {
		return nil, err
	}

	g := &Game{
		Name:        name,
		Users:       make(map[string]user.User),
		TotalRounds: totalRounds,
		Questions:   questions,
		mutex:       sync.RWMutex{},
	}
	g.CurrentRound = NewRound(g)
	return g, nil
}

func NewRound(g *Game) *Round {
	return &Round{
		Question: g.Questions[len(g.PastRounds)],
		Plays:    make(map[string]Play),
		Guesses:  make(map[string]Guess),
	}
}

func getQuestions(questions []string, num int) ([]string, error) { // TODO - move to S3
	i := 0
	var output []string
	usedIndexes := make(map[int]struct{})
	for {
		select {
		case <-time.After(time.Second):
			return nil, ErrNoQuestions
		default:
			if i >= num {
				return output, nil
			}
			index := rand.Intn(len(questions))
			if _, ok := usedIndexes[index]; ok {
				continue
			}
			usedIndexes[index] = struct{}{}
			output = append(output, questions[index])
			i++
		}
	}

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
	if gp.Play != nil {
		return gp.PlayRound(g)
	}
	if gp.Guess != nil {
		if err := gp.GuessRound(g); err != nil {
			return err
		}
	}

	if len(g.CurrentRound.Guesses) < len(g.Users) {
		return nil
	}

	// next round & play
	g.PastRounds = append(g.PastRounds, g.CurrentRound)
	if len(g.PastRounds) == g.TotalRounds {
		g.CurrentRound = nil
		return nil
	}
	round := NewRound(g)
	g.CurrentRound = round
	return nil
}

func (gp *GamePlay) PlayRound(g *Game) error {
	if gp.Play == nil {
		return ErrNilPlay
	}
	g.mutex.Lock()
	g.CurrentRound.Plays[gp.User.Name] = *gp.Play
	g.mutex.Unlock()
	return nil
}

func (gp *GamePlay) GuessRound(g *Game) error {
	if gp.Guess == nil {
		return ErrNilGuess
	}
	g.mutex.Lock()
	g.CurrentRound.Guesses[gp.User.Name] = *gp.Guess
	g.mutex.Unlock()
	return nil
}
