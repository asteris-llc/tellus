package storage

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMemoryStoreLifecycle(t *testing.T) {
	t.Parallel()

	ms := NewMemoryStore()
	project := "test"
	state := State("test")
	// put and get
	err := ms.Set(project, state)
	assert.Nil(t, err)

	state2, err := ms.Get(project)
	assert.Nil(t, err)
	assert.Equal(t, state, state2)

	// delete and get
	err = ms.Delete(project)
	assert.Nil(t, err)

	state2, err = ms.Get(project)
	assert.Equal(t, err, ErrNoProject)
	assert.Equal(t, state2, State{})
}
