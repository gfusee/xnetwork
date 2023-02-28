package mock

import (
	"errors"
	"sync"
)

type persisterMock struct {
	mut         sync.RWMutex
	data        map[string][]byte
	CloseCalled func() error
}

// NewPersisterMock -
func NewPersisterMock() *persisterMock {
	return &persisterMock{
		data: make(map[string][]byte),
	}
}

// Put -
func (mock *persisterMock) Put(key, val []byte) error {
	mock.mut.Lock()
	defer mock.mut.Unlock()

	mock.data[string(key)] = val

	return nil
}

// Get -
func (mock *persisterMock) Get(key []byte) ([]byte, error) {
	mock.mut.RLock()
	defer mock.mut.RUnlock()

	val, ok := mock.data[string(key)]
	if ok {
		return val, nil
	}

	return nil, errors.New("key not found")
}

// Has -
func (mock *persisterMock) Has(key []byte) error {
	mock.mut.RLock()
	defer mock.mut.RUnlock()

	_, ok := mock.data[string(key)]
	if !ok {
		return errors.New("key not found")
	}

	return nil
}

// Close -
func (mock *persisterMock) Close() error {
	if mock.CloseCalled != nil {
		return mock.CloseCalled()
	}

	return nil
}

// Remove -
func (mock *persisterMock) Remove(key []byte) error {
	mock.mut.Lock()
	defer mock.mut.Unlock()

	delete(mock.data, string(key))

	return nil
}

// Destroy -
func (mock *persisterMock) Destroy() error {
	return nil
}

// DestroyClosed -
func (mock *persisterMock) DestroyClosed() error {
	return nil
}

// RangeKeys -
func (mock *persisterMock) RangeKeys(handler func(key []byte, val []byte) bool) {
	if handler == nil {
		return
	}

	mock.mut.RLock()
	defer mock.mut.RUnlock()

	for key, data := range mock.data {
		shouldContinue := handler([]byte(key), data)
		if !shouldContinue {
			return
		}
	}
}

// IsInterfaceNil -
func (mock *persisterMock) IsInterfaceNil() bool {
	return mock == nil
}
