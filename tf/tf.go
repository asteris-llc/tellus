package tf

import (
	"errors"
	"github.com/asteris-llc/tellus/storage"
	"github.com/hashicorp/terraform/terraform"
	"strings"
)

var (
	ErrNoState   = errors.New("no state found")
	ErrNoModules = errors.New("no module found")
	statePrefix  = "state"
	modulePrefix = "module"
)

// state

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

// modules

type ModuleGetter interface {
	GetModules(string) ([]*Module, error)
}

type ModuleSetter interface {
	SetModules(string, []*Module) error
}

type ModuleDeleter interface {
	DeleteModules(string) error
}

type ModuleManipulator interface {
	ModuleGetter
	ModuleSetter
	ModuleDeleter
}

// all together now!

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
