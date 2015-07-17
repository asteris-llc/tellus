package tf

import (
	"errors"
	"github.com/hashicorp/terraform/terraform"
)

var (
	ErrNoState = errors.New("no state found")
)

type StateGetter interface {
	Get(string) (*terraform.State, error)
}

type StateSetter interface {
	Set(string, *terraform.State) error
}

type StateDeleter interface {
	Delete(string) error
}

type StateManipulator interface {
	StateGetter
	StateSetter
	StateDeleter
}
