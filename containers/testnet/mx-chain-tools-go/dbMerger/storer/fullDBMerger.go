package storer

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-go/storage"
	"github.com/multiversx/mx-chain-storage-go/types"
)

const minNumOfPersisters = 2

// ArgsFullDBMerger is the DTO used in the NewFullDBMerger constructor function
type ArgsFullDBMerger struct {
	DataMergerInstance  DataMerger
	PersisterCreator    PersisterCreator
	OsOperationsHandler OsOperationsHandler
}

type fullDBMerger struct {
	dataMergerInstance  DataMerger
	persisterCreator    PersisterCreator
	osOperationsHandler OsOperationsHandler
}

// NewFullDBMerger creates a new instance of type fullDBMerger
func NewFullDBMerger(args ArgsFullDBMerger) (*fullDBMerger, error) {
	if check.IfNil(args.DataMergerInstance) {
		return nil, fmt.Errorf("%w, DataMergerInstance", errNilComponent)
	}
	if check.IfNil(args.PersisterCreator) {
		return nil, fmt.Errorf("%w, PersisterCreator", errNilComponent)
	}
	if check.IfNil(args.OsOperationsHandler) {
		return nil, fmt.Errorf("%w, OsOperationsHandler", errNilComponent)
	}

	return &fullDBMerger{
		dataMergerInstance:  args.DataMergerInstance,
		persisterCreator:    args.PersisterCreator,
		osOperationsHandler: args.OsOperationsHandler,
	}, nil
}

// MergeDBs will merge all data from the source persister paths into a new storage persister
func (fdm *fullDBMerger) MergeDBs(destinationPath string, sourcePaths ...string) (storage.Persister, error) {
	if len(sourcePaths) < minNumOfPersisters {
		return nil, fmt.Errorf("%w, provided %d, minimum %d", errInvalidNumberOfPersisters, len(sourcePaths), minNumOfPersisters)
	}

	err := fdm.osOperationsHandler.CheckIfDirectoryIsEmpty(destinationPath)
	if err != nil {
		return nil, err
	}

	err = fdm.osOperationsHandler.CopyDirectory(destinationPath, sourcePaths[0])
	if err != nil {
		return nil, err
	}

	destPersister, err := fdm.persisterCreator.CreatePersister(destinationPath)
	if err != nil {
		return nil, fmt.Errorf("%w for destination persister", err)
	}

	sourcePersisters, err := fdm.createSourcePersisters(sourcePaths...)
	if err != nil {
		return nil, err
	}

	err = fdm.dataMergerInstance.MergeDBs(destPersister, sourcePersisters...)
	if err != nil {
		return nil, err
	}

	err = fdm.closeSourcePersisters(sourcePersisters)
	if err != nil {
		return nil, err
	}

	return destPersister, nil
}

func (fdm *fullDBMerger) createSourcePersisters(sourcePaths ...string) ([]types.Persister, error) {
	sourcePersisters := make([]types.Persister, 0, len(sourcePaths)-1)
	for i := 1; i < len(sourcePaths); i++ {
		srcPersister, errPersister := fdm.persisterCreator.CreatePersister(sourcePaths[i])
		if errPersister != nil {
			return nil, fmt.Errorf("%w for source persister with index %d", errPersister, i)
		}

		sourcePersisters = append(sourcePersisters, srcPersister)
	}

	return sourcePersisters, nil
}

func (fdm *fullDBMerger) closeSourcePersisters(sourcePersisters []types.Persister) error {
	var lastErrFound error

	for _, persister := range sourcePersisters {
		err := persister.Close()
		if err != nil {
			lastErrFound = err
		}
	}

	return lastErrFound
}

// IsInterfaceNil returns true if there is no value under the interface
func (fdm *fullDBMerger) IsInterfaceNil() bool {
	return fdm == nil
}
