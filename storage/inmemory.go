package storage

import (
	"fmt"
	"sync"

	"github.com/stinkyfingers/hilarity/game"
)

var (
	ErrNoGamesFound = fmt.Errorf("no games found")
	ErrGamenotFound = fmt.Errorf("game not found")
)

type InMemory struct {
	UserGamesMap map[string][]game.Game
	Games        map[string]game.Game
	mutex        sync.RWMutex
}

func NewInMemory() *InMemory {
	return &InMemory{
		UserGamesMap: make(map[string][]game.Game),
		Games:        make(map[string]game.Game),
		mutex:        sync.RWMutex{},
	}
}

var _ Storage = &InMemory{}

func (i *InMemory) NameExists(name string) (bool, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	_, ok := i.Games[name]
	return ok, nil
}

func (i *InMemory) ListGames() ([]string, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	var names []string
	for name := range i.Games {
		names = append(names, name)
	}
	return names, nil
}
func (i *InMemory) SaveGame(g *game.Game) error {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	i.Games[g.Name] = *g
	return nil
}

func (i *InMemory) GetGame(name string) (*game.Game, error) {
	i.mutex.Lock()
	defer i.mutex.Unlock()
	g, ok := i.Games[name]
	if !ok {
		return nil, ErrNoGamesFound
	}
	return &g, nil
}

func (c *InMemory) Lock(name string) {
	c.mutex.Lock()
}

func (c *InMemory) Unlock(name string) {
	c.mutex.Unlock()
}
