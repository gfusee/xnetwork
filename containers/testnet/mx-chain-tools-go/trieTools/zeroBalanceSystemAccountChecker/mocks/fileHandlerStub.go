package mocks

import (
	"io"

	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/common"
)

// FileHandlerStub -
type FileHandlerStub struct {
	OpenCalled    func(name string) (io.Reader, error)
	ReadAllCalled func(r io.Reader) ([]byte, error)
	GetwdCalled   func() (dir string, err error)
	ReadDirCalled func(dirname string) ([]common.FileInfo, error)
}

// Open -
func (fhs *FileHandlerStub) Open(name string) (io.Reader, error) {
	if fhs.OpenCalled != nil {
		return fhs.OpenCalled(name)
	}

	return nil, nil
}

// ReadAll -
func (fhs *FileHandlerStub) ReadAll(r io.Reader) ([]byte, error) {
	if fhs.ReadAllCalled != nil {
		return fhs.ReadAllCalled(r)
	}

	return nil, nil
}

// Getwd -
func (fhs *FileHandlerStub) Getwd() (dir string, err error) {
	if fhs.GetwdCalled != nil {
		return fhs.GetwdCalled()
	}

	return "", nil
}

// ReadDir -
func (fhs *FileHandlerStub) ReadDir(dirname string) ([]common.FileInfo, error) {
	if fhs.ReadDirCalled != nil {
		return fhs.ReadDirCalled(dirname)
	}

	return nil, nil
}
