package storer

import (
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core/check"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-storage-go/types"
)

var log = logger.GetOrCreate("storer")

// dataMerger is able to copy key by key all values from the provided sources persisters into the destination persister
type dataMerger struct {
}

// NewDataMerger returns a new instance of a data merger
func NewDataMerger() *dataMerger {
	return &dataMerger{}
}

// MergeDBs will iterate over all provided sources and take all key-value pairs and write them in the destination persister
func (dm *dataMerger) MergeDBs(dest types.Persister, sources ...types.Persister) error {
	err := checkArgs(dest, sources...)
	if err != nil {
		return err
	}

	numKeys := 0

	for _, source := range sources {
		copiedKeys, errMerge := mergeDB(dest, source)
		if errMerge != nil {
			return errMerge
		}

		numKeys += copiedKeys
	}

	log.Debug("finished copying data",
		"num source persisters", len(sources), "num key-values copied", numKeys)

	return nil
}

func checkArgs(dest types.Persister, sources ...types.Persister) error {
	if check.IfNil(dest) {
		return fmt.Errorf("%w for the destination persister", errNilPersister)
	}
	for idx, source := range sources {
		if check.IfNil(source) {
			return fmt.Errorf("%w for the source persister, index %d", errNilPersister, idx)
		}
	}

	return nil
}

func mergeDB(dest types.Persister, source types.Persister) (int, error) {
	var foundErr error
	numKeysCopied := 0
	source.RangeKeys(func(key []byte, val []byte) bool {
		numKeysCopied++
		foundErr = dest.Put(key, val)
		if foundErr != nil {
			return false
		}

		return true
	})

	return numKeysCopied, foundErr
}

// IsInterfaceNil returns true if there is no value under the interface
func (dm *dataMerger) IsInterfaceNil() bool {
	return dm == nil
}
