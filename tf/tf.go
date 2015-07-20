package tf

import (
	"errors"
	"github.com/asteris-llc/tellus/storage"
	"github.com/hashicorp/terraform/terraform"
	"strings"
)

var (
	ErrNoState = errors.New("no state found")
)

type StateGetter interface {
	GetState(string) (*terraform.State, error)
}

type StateSetter interface {
	SetState(string, *terraform.State) error
}

type StateDeleter interface {
	DeleteState(string) error
}

type StateManipulator interface {
	StateGetter
	StateSetter
	StateDeleter
}

type Terraformer struct {
	store storage.BlobStorer
}

func New(s storage.BlobStorer) *Terraformer {
	return &Terraformer{s}
}

func (t *Terraformer) withPrefix(s, prefix string) string {
	if strings.HasPrefix(s, prefix+"/") {
		return s
	} else {
		return prefix + "/" + s
	}
}
