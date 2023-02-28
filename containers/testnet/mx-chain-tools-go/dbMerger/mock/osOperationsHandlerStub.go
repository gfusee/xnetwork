package mock

// OsOperationsHandlerStub -
type OsOperationsHandlerStub struct {
	CheckIfDirectoryIsEmptyCalled func(directory string) error
	CopyDirectoryCalled           func(destination string, source string) error
}

// CopyDirectory -
func (stub *OsOperationsHandlerStub) CopyDirectory(destination string, source string) error {
	if stub.CopyDirectoryCalled != nil {
		return stub.CopyDirectoryCalled(destination, source)
	}

	return nil
}

// CheckIfDirectoryIsEmpty -
func (stub *OsOperationsHandlerStub) CheckIfDirectoryIsEmpty(directory string) error {
	if stub.CheckIfDirectoryIsEmptyCalled != nil {
		return stub.CheckIfDirectoryIsEmptyCalled(directory)
	}

	return nil
}

// IsInterfaceNil -
func (stub *OsOperationsHandlerStub) IsInterfaceNil() bool {
	return stub == nil
}
