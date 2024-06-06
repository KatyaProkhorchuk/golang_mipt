//go:build !solution

package keylock

import (
	"sort"
	"sync"
)

type KeyLock struct{
	mu sync.Mutex
	locks map[string]chan struct{}
}

func New() *KeyLock {
	return &KeyLock {
		locks: make(map[string]chan struct{}),
	}
}

func (l *KeyLock) LockKeys(keys []string, cancel <-chan struct{}) (canceled bool, unlock func()) {
	passedKeys := make([]string, len(keys))
	copy(passedKeys, keys)
	sort.Strings(passedKeys)
	wait := make([]chan struct{}, 0)
	unlock = func() {
		l.mu.Lock()
		for key := range l.locks {
			delete(l.locks, key)
		}
		l.mu.Unlock()
		for _, lock := range wait {
			close(lock)
		}
	}
	for _, key := range passedKeys {
		for {
			l.mu.Lock()
			other, already := l.locks[key]
			lock := make(chan struct{})
			if !already {
				l.locks[key] = lock
			}
			l.mu.Unlock()
			if !already {
				wait = append(wait, lock)
				break
			}
			select {
			case <- cancel:
				unlock()
				return true, unlock
			case <-other:
				continue
			} 
		}
	}
	return false, unlock
}
