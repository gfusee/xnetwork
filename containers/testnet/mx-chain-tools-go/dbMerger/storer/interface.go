package storer

import (
	"github.com/multiversx/mx-chain-storage-go/types"
)

// DataMerger specify the operations supported by a component able to merge data between persisters
type DataMerger interface {
	MergeDBs(dest types.Persister, sources ...types.Persister) error
	IsInterfaceNil() bool
}

// PersisterCreator is able to create a persister instance based on the provided path
type PersisterCreator interface {
	CreatePersister(path string) (types.Persister, error)
	IsInterfaceNil() bool
}

// OsOperationsHandler is able to handle the os-level functions
type OsOperationsHandler interface {
	CheckIfDirectoryIsEmpty(directory string) error
	CopyDirectory(destination string, source string) error
	IsInterfaceNil() bool
}
