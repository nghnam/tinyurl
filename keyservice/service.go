package keyservice

import (
	"sync"
)

type KeyService struct {
	sync.Mutex
	remainKeys []string
	usedKeys   []string
	keyCh      chan chan string
}

var once sync.Once
var Service *KeyService
var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")

func NewKeyService(keyCh chan chan string) *KeyService {
	once.Do(func() {
		Service = &KeyService{
			keyCh: keyCh,
		}
	})
	return Service
}

func (ks *KeyService) Run() {
	var wg sync.WaitGroup
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func(i int) {
			for {
				select {
				case reqCh := <-ks.keyCh:
					ks.Lock()
					if len(ks.remainKeys) > 0 {
						key := ks.remainKeys[len(ks.remainKeys)-1]
						ks.remainKeys = ks.remainKeys[:len(ks.remainKeys)-1]
						ks.usedKeys = append(ks.usedKeys, key)
						reqCh <- key
					} else {
						close(reqCh)
					}
					ks.Unlock()
				}
			}
		}(i)
	}
	wg.Wait()
}

func (ks *KeyService) CreateKeys(n int, length int) {
	b := make([]rune, length)
	for i := 0; i < n; i++ {
		m := i
		for j := 0; j < length; j++ {
			r := m % len(letters)
			m = m / len(letters)
			b[j] = letters[r]
		}
		ks.remainKeys = append(ks.remainKeys, string(b))
	}
}
