package mocks

// FileStub -
type FileStub struct {
	NameCalled  func() string
	IsDirCalled func() bool
}

// Name -
func (fs *FileStub) Name() string {
	if fs.NameCalled != nil {
		return fs.NameCalled()
	}

	return ""
}

// IsDir -
func (fs *FileStub) IsDir() bool {
	if fs.IsDirCalled != nil {
		return fs.IsDirCalled()
	}

	return false
}
