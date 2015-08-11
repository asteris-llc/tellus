package tf

import (
	"encoding/json"
	"github.com/asteris-llc/tellus/storage"
)

type Module struct {
	Path    string
	Content string
	Link    bool
}

func (t *Terraformer) GetModules(name string) ([]*Module, error) {
	blob, err := t.store.Get(t.withPrefix(name, modulePrefix))
	if err == storage.ErrNoBlob {
		return nil, ErrNoModules
	} else if err != nil {
		return nil, err
	}

	modules := []*Module{}
	err = json.Unmarshal(blob, &modules)
	if err != nil {
		return nil, err
	}

	return modules, nil
}

func (t *Terraformer) SetModules(name string, modules []*Module) error {
	blob, err := json.Marshal(modules)
	if err != nil {
		return err
	}

	err = t.store.Set(t.withPrefix(name, modulePrefix), blob)
	return err
}

func (t *Terraformer) DeleteModules(name string) error {
	err := t.store.Delete(t.withPrefix(name, modulePrefix))

	switch err {
	case storage.ErrNoBlob, nil:
		return nil
	default:
		return err
	}
}
