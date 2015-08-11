package tf

import (
	"encoding/json"
	"github.com/asteris-llc/tellus/storage"
	"github.com/stretchr/testify/assert"
	"testing"
)

var testModules = []*Module{
	&Module{"test.tf", "// test", false},
}

func TestGetModules(t *testing.T) {
	t.Parallel()

	name := "test"
	key := modulePrefix + "/" + name
	mem := storage.NewMemoryStore()
	tf := New(mem)

	// empty case
	state, err := tf.GetModules(name)
	assert.Nil(t, state)
	assert.Equal(t, err, ErrNoModules)

	// good serialization
	blob, err := json.Marshal(testModules)
	assert.Nil(t, err)
	mem.Set(key, blob)

	state, err = tf.GetModules(name)
	assert.Nil(t, err)
	assert.Equal(t, state, testModules)
}

func TestSetModules(t *testing.T) {
	t.Parallel()

	name := "test"
	key := modulePrefix + "/" + name
	mem := storage.NewMemoryStore()
	tf := New(mem)

	// the current implementation of this code (2015-07-20) relies heavily on
	// upstream terraform's code being valid, and downstream being reasonable. So
	// we only are testing the "happy path" here, since the state storage is
	// tested separately.
	_, err := mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)

	err = tf.SetModules(name, testModules)
	assert.Nil(t, err)

	_, err = mem.Get(key)
	assert.Nil(t, err)
}

func TestDeleteModules(t *testing.T) {
	t.Parallel()

	name := "test"
	key := modulePrefix + "/" + name
	mem := storage.NewMemoryStore()
	tf := New(mem)

	// deleting a key that isn't there doesn't do anything
	_, err := mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)

	err = tf.DeleteModules(name)
	assert.Nil(t, err)

	_, err = mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)

	// deleting a key that *is* there does something
	err = mem.Set(key, []byte("test"))
	assert.Nil(t, err)

	err = tf.DeleteModules(name)
	assert.Nil(t, err)

	_, err = mem.Get(key)
	assert.Equal(t, err, storage.ErrNoBlob)
}
