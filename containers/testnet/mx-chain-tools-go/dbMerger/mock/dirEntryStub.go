package mock

import (
	"errors"
	"io/fs"
)

// DirEntryStub -
type DirEntryStub struct {
	NameValue  string
	IsDirValue bool
	TypeValue  fs.FileMode
	InfoCalled func() (fs.FileInfo, error)
}

// Name -
func (stub *DirEntryStub) Name() string {
	return stub.NameValue
}

// IsDir -
func (stub *DirEntryStub) IsDir() bool {
	return stub.IsDirValue
}

// Type -
func (stub *DirEntryStub) Type() fs.FileMode {
	return stub.TypeValue
}

// Info -
func (stub *DirEntryStub) Info() (fs.FileInfo, error) {
	if stub.InfoCalled != nil {
		return stub.InfoCalled()
	}

	return nil, errors.New("not implemented")
}
