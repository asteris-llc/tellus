package tf

import (
	"bytes"
	"encoding/json"
	"github.com/asteris-llc/tellus/storage"
	"github.com/hashicorp/terraform/terraform"
)

var (
	statePrefix = "state"
)

func (t *Terraformer) GetState(name string) (*terraform.State, error) {
	blob, err := t.store.Get(t.withPrefix(name, statePrefix))
	if err == storage.ErrNoBlob {
		return nil, ErrNoState
	} else if err != nil {
		return nil, err
	}

	state, err := terraform.ReadState(bytes.NewBuffer(blob))
	if err != nil {
		return nil, err
	}

	return state, nil
}

func (t *Terraformer) SetState(name string, state *terraform.State) error {
	blob, err := json.Marshal(state)
	if err != nil {
		return err
	}

	err = t.store.Set(t.withPrefix(name, statePrefix), blob)
	return err
}

func (t *Terraformer) DeleteState(name string) error {
	err := t.store.Delete(t.withPrefix(name, statePrefix))

	switch err {
	case storage.ErrNoBlob, nil:
		return nil
	default:
		return err
	}
}
