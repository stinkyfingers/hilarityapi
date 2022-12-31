package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/stinkyfingers/hilarity/game"
)

type File struct {
	Path    string
	mutexes map[string]*sync.RWMutex
}

var _ Storage = &File{}

func NewFile(path string) (*File, error) {
	if err := os.Mkdir(path, os.ModePerm); err != nil && !os.IsExist(err) {
		return nil, err
	}
	mutexes := make(map[string]*sync.RWMutex)
	infos, err := os.ReadDir(path)
	if err != nil {
		return nil, err
	}
	for _, info := range infos {
		mutexes[info.Name()] = &sync.RWMutex{}
	}
	return &File{
		Path:    path,
		mutexes: mutexes,
	}, nil
}

// not threadsafe - call lock/unlock
func (f *File) read(name string) (*game.Game, error) {
	file, err := os.Open(filepath.Join(f.Path, fmt.Sprintf("%s.json", name)))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var g game.Game
	err = json.NewDecoder(file).Decode(&g)
	if err != nil {
		return nil, err
	}
	return &g, nil
}

// not threadsafe - call lock/unlock
func (f *File) write(g *game.Game) error {
	path := filepath.Join(f.Path, fmt.Sprintf("%s.json", g.Name))
	var err error
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	err = json.NewEncoder(file).Encode(g)
	return err
}

func (f *File) NameExists(name string) (bool, error) {
	dirs, err := os.ReadDir(f.Path)
	if err != nil {
		return false, err
	}
	for _, dir := range dirs {
		info, err := dir.Info()
		if err != nil {
			return true, err
		}
		if filepath.Base(info.Name()) == fmt.Sprintf("%s.json", name) {
			return true, nil
		}
	}
	return false, nil
}

func (f *File) GetGame(name string) (*game.Game, error) {
	return f.read(name)
}

func (f *File) ListGames() ([]string, error) {
	infos, err := os.ReadDir(f.Path)
	if err != nil {
		return nil, err
	}
	var games []string
	for _, info := range infos {
		if info.Name() == questionKey {
			continue
		}
		games = append(games, strings.TrimSuffix(info.Name(), filepath.Ext(info.Name())))
	}
	return games, nil
}

func (f *File) SaveGame(game *game.Game) error {
	ok, err := f.NameExists(game.Name)
	if err != nil {
		return err
	}
	if !ok {
		f.mutexes[game.Name] = &sync.RWMutex{}
	}
	return f.write(game)
}

func (f *File) CleanUpGames() error {
	infos, err := os.ReadDir(f.Path)
	if err != nil {
		return err
	}

	for _, info := range infos {
		if info.Name() == questionKey {
			continue
		}
		fileInfo, err := info.Info()
		if err != nil {
			return err
		}
		if fileInfo.ModTime().Before(time.Now().Add(-24 * time.Hour)) {
			err = os.Remove(fileInfo.Name())
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (f *File) Lock(name string) {
	if _, ok := f.mutexes[name]; ok {
		f.mutexes[name].Lock()
	}
}

func (f *File) Unlock(name string) {
	if _, ok := f.mutexes[name]; ok {
		f.mutexes[name].Unlock()
	}
}

func (f *File) GetQuestions() ([]string, error) {
	file, err := os.Open(filepath.Join(f.Path, questionKey))
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var output []string
	err = json.NewDecoder(file).Decode(&output)
	return output, err
}
