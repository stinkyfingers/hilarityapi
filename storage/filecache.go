package storage

//import (
//	"github.com/stinkyfingers/hilarity/game"
//)
//
//// FileCache combines a File-based storage w/ in-mem cache
//type FileCache struct {
//	File     File
//	InMemory InMemory
//}
//
//func NewFileCache(f File, im InMemory) *FileCache {
//	return &FileCache{
//		File:     f,
//		InMemory: im,
//	}
//}
//
//func (f *FileCache) NameExists(name string) (bool, error) {
//	if ok, err := f.InMemory.NameExists(name); ok && err == nil {
//		return true, nil
//	}
//	return f.File.NameExists(name)
//}
//func (f *FileCache) ListGames() ([]string, error)       {
//	if
//}
//func (f *FileCache) SaveGame(*game.Game) error          {}
//func (f *FileCache) GetGame(string) (*game.Game, error) {}
