package storage

import (
	"errors"
)

var (
	ErrNoProject = errors.New("project not found")
)

type State []byte

type StateGetter interface {
	Get(string) (State, error)
}

type StateSetter interface {
	Set(string, State) error
}

type StateDeleter interface {
	Delete(string) error
}
