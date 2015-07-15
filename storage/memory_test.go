package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStoreLifecycle(t *testing.T) {
	t.Parallel()

	ms := NewMemoryStore()
	key := "test"
	state := State("test")
	// put and get
	err := ms.Set(key, state)
	assert.Nil(t, err)

	state2, err := ms.Get(key)
	assert.Nil(t, err)
	assert.Equal(t, state, state2)

	// delete and get
	err = ms.Delete(key)
	assert.Nil(t, err)

	state2, err = ms.Get(key)
	assert.Equal(t, err, ErrNoProject)
	assert.Equal(t, state2, State{})
}
