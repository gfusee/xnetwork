package components

import (
	"fmt"
	"path/filepath"
)

const staticIdentifier = "Static"

type simplePathManager struct {
	workingDir string
}

// NewSimplePathManager creates a new instance of the simplePathManager
func NewSimplePathManager(workingDir string) *simplePathManager {
	return &simplePathManager{
		workingDir: workingDir,
	}
}

// PathForEpoch returns the path for epoch taking into account only the epoch
func (spm *simplePathManager) PathForEpoch(_ string, epoch uint32, _ string) string {
	return filepath.Join(spm.workingDir, fmt.Sprintf("%d", epoch))
}

// PathForStatic returns a dummy static path
func (spm *simplePathManager) PathForStatic(_ string, _ string) string {
	return filepath.Join(spm.workingDir, staticIdentifier)
}

// DatabasePath returns the raw working dir
func (spm *simplePathManager) DatabasePath() string {
	return spm.workingDir
}

// IsInterfaceNil returns true if there is no value under the interface
func (spm *simplePathManager) IsInterfaceNil() bool {
	return spm == nil
}
