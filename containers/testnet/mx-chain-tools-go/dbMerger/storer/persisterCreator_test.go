package storer

import (
	"fmt"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/stretchr/testify/assert"
)

func TestNewPersisterCreator(t *testing.T) {
	t.Parallel()

	creator := NewPersisterCreator()
	assert.False(t, check.IfNil(creator))
}

func TestPersisterCreator_CreatePersister(t *testing.T) {
	t.Parallel()

	creator := NewPersisterCreator()

	persister, err := creator.CreatePersister("test")
	assert.False(t, check.IfNil(persister))
	assert.Nil(t, err)
	assert.Equal(t, "*leveldb.DB", fmt.Sprintf("%T", persister))

	_ = persister.Destroy()
}
