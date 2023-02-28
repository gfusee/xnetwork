package mock

import (
	"bytes"
)

// ElasticClientStub -
type ElasticClientStub struct {
	GetMappingCalled                  func(index string) (*bytes.Buffer, error)
	CreateIndexWithMappingCalled      func(targetIndex string, body *bytes.Buffer) error
	DoScrollRequestAllDocumentsCalled func(index string, body []byte, handlerFunc func(responseBytes []byte) error) error
	GetCountCalled                    func(index string) (uint64, error)
	DoesAliasExistCalled              func(alias string) bool
	DoBulkRequestCalled               func(buff *bytes.Buffer, index string) error
	DoesIndexExistCalled              func(index string) bool
	PutAliasCalled                    func(index string, alias string) error
}

// GetMapping -
func (e *ElasticClientStub) GetMapping(index string) (*bytes.Buffer, error) {
	if e.GetMappingCalled != nil {
		return e.GetMappingCalled(index)
	}

	return nil, nil
}

// CreateIndexWithMapping -
func (e *ElasticClientStub) CreateIndexWithMapping(targetIndex string, body *bytes.Buffer) error {
	if e.CreateIndexWithMappingCalled != nil {
		return e.CreateIndexWithMappingCalled(targetIndex, body)
	}

	return nil
}

// DoScrollRequestAllDocuments -
func (e *ElasticClientStub) DoScrollRequestAllDocuments(index string, body []byte, handlerFunc func(responseBytes []byte) error) error {
	if e.DoScrollRequestAllDocumentsCalled != nil {
		return e.DoScrollRequestAllDocumentsCalled(index, body, handlerFunc)
	}

	return nil
}

// GetCount -
func (e *ElasticClientStub) GetCount(index string) (uint64, error) {
	if e.GetCountCalled != nil {
		return e.GetCountCalled(index)
	}

	return 0, nil
}

// DoesAliasExist -
func (e *ElasticClientStub) DoesAliasExist(alias string) bool {
	if e.DoesAliasExistCalled != nil {
		return e.DoesAliasExistCalled(alias)
	}

	return false
}

// DoBulkRequest -
func (e *ElasticClientStub) DoBulkRequest(buff *bytes.Buffer, index string) error {
	if e.DoBulkRequestCalled != nil {
		return e.DoBulkRequestCalled(buff, index)
	}

	return nil
}

// DoesIndexExist -
func (e *ElasticClientStub) DoesIndexExist(index string) bool {
	if e.DoesIndexExistCalled != nil {
		return e.DoesIndexExistCalled(index)
	}

	return false
}

// PutAlias -
func (e *ElasticClientStub) PutAlias(index string, alias string) error {
	if e.PutAliasCalled != nil {
		return e.PutAliasCalled(index, alias)
	}

	return nil
}

// IsInterfaceNil -
func (e *ElasticClientStub) IsInterfaceNil() bool {
	return e == nil
}
