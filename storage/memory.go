package storage

import (
	"sync"
)

type MemoryStore struct {
	store map[string]State
	sync  sync.RWMutex
}

func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		map[string]State{},
		sync.RWMutex{},
	}
}

func (ms *MemoryStore) Get(key string) (State, error) {
	ms.sync.RLock()
	defer ms.sync.RUnlock()

	state, present := ms.store[key]
	if !present {
		return State{}, ErrNoProject
	}

	return state, nil
}

func (ms *MemoryStore) Set(key string, state State) error {
	ms.sync.Lock()
	defer ms.sync.Unlock()

	ms.store[key] = state
	return nil
}

func (ms *MemoryStore) Delete(key string) error {
	ms.sync.Lock()
	defer ms.sync.Unlock()

	// this isn't *strictly* necessary here because delete is a no-op if the
	// project doesn't exist, but the rest of the stores should raise an error in
	// this case so we do here too.
	_, present := ms.store[key]
	if !present {
		return ErrNoProject
	}

	delete(ms.store, key)
	return nil
}
