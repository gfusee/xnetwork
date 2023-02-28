package common

import (
	"fmt"
	"path/filepath"
)

const staticIdentifier = "Static"

type simplePathManager struct {
	dbPath string
}

// NewSimplePathManager creates a new instance of the simplePathManager
func NewSimplePathManager(dbPath string) *simplePathManager {
	return &simplePathManager{
		dbPath: dbPath,
	}
}

// PathForEpoch returns the path for epoch taking into account the epoch and the shard
func (spm *simplePathManager) PathForEpoch(shardId string, epoch uint32, identifier string) string {
	epochPart := fmt.Sprintf("Epoch_%d", epoch)
	shardPart := fmt.Sprintf("Shard_%s", shardId)
	path := filepath.Join(spm.dbPath, epochPart, shardPart, identifier)
	return path
}

// PathForStatic returns a dummy static path
func (spm *simplePathManager) PathForStatic(_ string, _ string) string {
	panic("not implemented")
}

// DatabasePath returns the raw working dir
func (spm *simplePathManager) DatabasePath() string {
	return spm.dbPath
}

// IsInterfaceNil returns true if there is no value under the interface
func (spm *simplePathManager) IsInterfaceNil() bool {
	return spm == nil
}
