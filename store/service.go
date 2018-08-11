package store

import (
	"errors"
	"sync"
)

var once sync.Once

type MapDB struct {
	sync.Mutex
	keyUrl map[string]string
}

var DB *MapDB

func NewMapDB() *MapDB {
	once.Do(func() {
		DB = &MapDB{
			keyUrl: make(map[string]string),
		}
	})
	return DB
}

func (db *MapDB) Update(key string, url string) {
	db.Lock()
	defer db.Unlock()
	db.keyUrl[key] = url
}

func (db *MapDB) Lookup(key string) (string, error) {
	if url, ok := db.keyUrl[key]; ok {
		return url, nil
	}
	return "", errors.New("Key not found")
}
