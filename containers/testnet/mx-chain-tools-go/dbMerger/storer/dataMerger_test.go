package storer

import (
	"errors"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-tools-go/dbmerger/mock"
	"github.com/stretchr/testify/assert"
)

func TestNewDataMerger(t *testing.T) {
	t.Parallel()

	dm := NewDataMerger()
	assert.False(t, check.IfNil(dm))
}

func TestMergeDBs(t *testing.T) {
	t.Parallel()

	t.Run("nil destination should error", func(t *testing.T) {
		t.Parallel()

		dm := NewDataMerger()
		err := dm.MergeDBs(nil)
		assert.True(t, errors.Is(err, errNilPersister))
		assert.True(t, strings.Contains(err.Error(), "for the destination persister"))
	})
	t.Run("sources contains a nil persister should error", func(t *testing.T) {
		t.Parallel()

		dm := NewDataMerger()
		err := dm.MergeDBs(&mock.PersisterStub{}, nil)
		assert.True(t, errors.Is(err, errNilPersister))
		assert.True(t, strings.Contains(err.Error(), "for the source persister, index 0"))

		err = dm.MergeDBs(&mock.PersisterStub{}, &mock.PersisterStub{}, nil)
		assert.True(t, errors.Is(err, errNilPersister))
		assert.True(t, strings.Contains(err.Error(), "for the source persister, index 1"))
	})
	t.Run("empty sources list should not put", func(t *testing.T) {
		t.Parallel()

		dm := NewDataMerger()
		err := dm.MergeDBs(&mock.PersisterStub{
			PutCalled: func(key, val []byte) error {
				assert.Fail(t, "should have not called put")
				return nil
			},
		})
		assert.Nil(t, err)
	})
	t.Run("put works, should merge", func(t *testing.T) {
		t.Parallel()

		src1 := map[string]string{
			"key1": "val1",
			"key2": "val2",
		}

		src2 := map[string]string{
			"key3": "val3",
			"key4": "val4",
		}

		src3 := map[string]string{
			"key5": "val5",
			"key6": "val6",
		}

		result := make(map[string]string)

		dm := NewDataMerger()
		err := dm.MergeDBs(&mock.PersisterStub{
			PutCalled: func(key, val []byte) error {
				result[string(key)] = string(val)
				return nil
			},
		},
			createPersisterStub(src1),
			createPersisterStub(src2),
			createPersisterStub(src3),
		)

		assert.Nil(t, err)
		assert.Equal(t, len(src1)+len(src2)+len(src3), len(result))
		checkMapContained(t, result, src1)
		checkMapContained(t, result, src2)
		checkMapContained(t, result, src3)
	})
	t.Run("put errors, should error", func(t *testing.T) {
		t.Parallel()

		expectedErr := errors.New("expected error")
		src1 := map[string]string{
			"key1": "val1",
			"key2": "val2",
		}

		dm := NewDataMerger()
		err := dm.MergeDBs(&mock.PersisterStub{
			PutCalled: func(key, val []byte) error {
				return expectedErr
			},
		},
			createPersisterStub(src1),
		)

		assert.Equal(t, expectedErr, err)
	})
}

func createPersisterStub(rangeMap map[string]string) *mock.PersisterStub {
	return &mock.PersisterStub{
		RangeKeysCalled: func(handler func(key []byte, val []byte) bool) {
			for key, val := range rangeMap {
				handler([]byte(key), []byte(val))
			}
		},
	}
}

func checkMapContained(tb testing.TB, mainMap map[string]string, subset map[string]string) {
	for key, val := range subset {
		recovered, found := mainMap[key]
		assert.True(tb, found)
		assert.Equal(tb, val, recovered)
	}
}
