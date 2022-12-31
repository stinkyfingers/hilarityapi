package storage

import (
	"github.com/stinkyfingers/hilarity/game"
)

type Storage interface {
	NameExists(name string) (bool, error)
	ListGames() ([]string, error)
	SaveGame(*game.Game) error
	GetGame(string) (*game.Game, error)
	CleanUpGames() error
	GetQuestions() ([]string, error)
	Lock(string)
	Unlock(string)
}

const questionKey = "questions.json"

func SubmitPlay(storage Storage, gp game.GamePlay) (*game.Game, error) {
	//storage.Lock(gp.GameName)
	//defer storage.Unlock(gp.GameName)
	g, err := storage.GetGame(gp.GameName)
	if err != nil {
		return nil, err
	}

	err = gp.MakePlay(g)
	if err != nil {
		return nil, err
	}

	err = storage.SaveGame(g)
	if err != nil {
		return nil, err
	}
	return g, nil
}
