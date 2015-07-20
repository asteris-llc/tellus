package tf

import (
	"encoding/json"
	"github.com/asteris-llc/tellus/storage"
	"github.com/hashicorp/terraform/terraform"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetState(t *testing.T) {
	t.Parallel()

	name := "test"
	key := statePrefix + "/" + name
	mem := storage.NewMemoryStore()
	tf := New(mem)

	// empty case
	state, err := tf.GetState(name)
	assert.Nil(t, state)
	assert.Equal(t, err, ErrNoState)

	// bad serialization
	mem.Set(key, []byte("blah"))
	state, err = tf.GetState(name)
	assert.Nil(t, state)
	assert.Equal(t, err.Error(), "Failed to check for magic bytes: EOF")

	// good serialization
	goodState := terraform.NewState()
	blob, err := json.Marshal(goodState)
	assert.Nil(t, err)
	mem.Set(key, blob)

	state, err = tf.GetState(name)
	assert.Nil(t, err)
	assert.Equal(t, state, goodState)
}

func TestSetState(t *testing.T) {
	t.Parallel()

	name := "test"
	key := statePrefix + "/" + name
	mem := storage.NewMemoryStore()
	tf := New(mem)

	// the current implementation of this code (2015-07-20) relies heavily on
	// upstream terraform's code being valid, and downstream being reasonable. So
	// we only are testing the "happy path" here, since the state storage is
	// tested separately.
	_, err := mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)

	err = tf.SetState(name, terraform.NewState())
	assert.Nil(t, err)

	_, err = mem.Get(key)
	assert.Nil(t, err)
}

func TestDeleteState(t *testing.T) {
	t.Parallel()

	name := "test"
	key := statePrefix + "/" + name
	mem := storage.NewMemoryStore()
	tf := New(mem)

	// deleting a key that isn't there doesn't do anything
	_, err := mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)

	err = tf.DeleteState(name)
	assert.Nil(t, err)

	_, err = mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)

	// deleting a key that *is* there does something
	err = mem.Set(key, []byte("test"))
	assert.Nil(t, err)

	err = tf.DeleteState(name)
	assert.Nil(t, err)

	_, err = mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)
}
