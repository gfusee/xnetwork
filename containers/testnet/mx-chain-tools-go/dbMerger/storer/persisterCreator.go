package storer

import (
	"github.com/multiversx/mx-chain-storage-go/leveldb"
	"github.com/multiversx/mx-chain-storage-go/types"
)

const (
	batchDelaySeconds = 2
	maxBatchSize      = 10000
	maxOpenFiles      = 10
)

type persisterCreator struct {
}

// NewPersisterCreator will create a new persister creator instance
func NewPersisterCreator() *persisterCreator {
	return &persisterCreator{}
}

// CreatePersister will try to create a new persister instance provided the directory path
func (creator *persisterCreator) CreatePersister(path string) (types.Persister, error) {
	return leveldb.NewDB(path, batchDelaySeconds, maxBatchSize, maxOpenFiles)
}

// IsInterfaceNil returns true if there is no value under the interface
func (creator *persisterCreator) IsInterfaceNil() bool {
	return creator == nil
}
