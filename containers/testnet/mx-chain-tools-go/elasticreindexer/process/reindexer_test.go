package process

import (
	"bytes"
	"testing"

	"github.com/multiversx/mx-chain-tools-go/elasticreindexer/process/mock"
	"github.com/stretchr/testify/require"
)

const testIndex = "index"

func TestCopyMapping(t *testing.T) {
	t.Run("no overwrite - no alias - no index => should create",
		testCopyMappingNoOverwriteShouldCreate)
	t.Run("no overwrite - no alias - index => should err",
		testCopyMappingNoOverwriteAliasDoesNotExistsIndexExistsShouldErr)
	t.Run("no overwrite - alias - no index => should err",
		testCopyMappingNoOverwriteAliasExistsIndexDoesNotExistShouldErr)
	t.Run("no overwrite - alias - index => should err",
		testCopyMappingNoOverwriteAliasAndIndexExistShouldErr)

	t.Run("overwrite - no alias - no index => should copy/create",
		testCopyMappingOverwriteAliasAndIndexExistShouldCreate)
	t.Run("overwrite - no alias - index => should create alias",
		testCopyMappingOverwriteAliasDoesNotExistShouldCreateAlias)
	t.Run("overwrite - alias - no index => should create alias",
		testCopyMappingOverwriteAliasExistsIndexDoesNotExistShouldCreateIndex)
	t.Run("overwrite - alias - index => should not copy/create",
		testCopyMappingOverwriteAliasAndIndexExistShouldNotCreate)
	t.Run("skip-mappings",
		testSkipMappings)
}

func testCopyMappingNoOverwriteShouldCreate(t *testing.T) {
	getMappingCalled := false
	putAliasCalled := false
	createIndexCalled := false
	sourceClient := &mock.ElasticClientStub{
		GetMappingCalled: func(_ string) (*bytes.Buffer, error) {
			getMappingCalled = true
			return nil, nil
		},
	}

	destinationClient := &mock.ElasticClientStub{
		DoesIndexExistCalled: func(_ string) bool {
			return false
		},
		DoesAliasExistCalled: func(_ string) bool {
			return false
		},
		PutAliasCalled: func(_ string, _ string) error {
			putAliasCalled = true
			return nil
		},
		CreateIndexWithMappingCalled: func(_ string, _ *bytes.Buffer) error {
			createIndexCalled = true
			return nil
		},
	}

	r, _ := newReindexer(sourceClient, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, false, false)
	require.NoError(t, err)
	require.True(t, getMappingCalled)
	require.True(t, putAliasCalled)
	require.True(t, createIndexCalled)
}

func testCopyMappingNoOverwriteAliasExistsIndexDoesNotExistShouldErr(t *testing.T) {
	destinationClient := &mock.ElasticClientStub{
		DoesAliasExistCalled: func(_ string) bool {
			return true
		},
		DoesIndexExistCalled: func(_ string) bool {
			return false
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, false, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "index with alias index already exists.")
}

func testCopyMappingNoOverwriteAliasAndIndexExistShouldErr(t *testing.T) {
	destinationClient := &mock.ElasticClientStub{
		DoesAliasExistCalled: func(_ string) bool {
			return true
		},
		DoesIndexExistCalled: func(_ string) bool {
			return true
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, false, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "index with alias index already exists.")
}

func testCopyMappingNoOverwriteAliasDoesNotExistsIndexExistsShouldErr(t *testing.T) {
	destinationClient := &mock.ElasticClientStub{
		DoesIndexExistCalled: func(_ string) bool {
			return true
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"test-index"})

	err := r.copyMappingIfNecessary("test-index", false, false)
	require.Error(t, err)
	require.Contains(t, err.Error(), "index test-index already exists.")
}

func testCopyMappingOverwriteAliasAndIndexExistShouldNotCreate(t *testing.T) {
	getMappingCalled := false
	putAliasCalled := false
	createIndexCalled := false
	sourceClient := &mock.ElasticClientStub{
		GetMappingCalled: func(_ string) (*bytes.Buffer, error) {
			getMappingCalled = true
			return nil, nil
		},
	}

	destinationClient := &mock.ElasticClientStub{
		DoesIndexExistCalled: func(_ string) bool {
			return true
		},
		DoesAliasExistCalled: func(_ string) bool {
			return true
		},
		PutAliasCalled: func(_ string, _ string) error {
			putAliasCalled = true
			return nil
		},
		CreateIndexWithMappingCalled: func(_ string, _ *bytes.Buffer) error {
			createIndexCalled = true
			return nil
		},
	}

	r, _ := newReindexer(sourceClient, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, true, false)
	require.NoError(t, err)
	require.False(t, getMappingCalled)
	require.False(t, putAliasCalled)
	require.False(t, createIndexCalled)
}

func testCopyMappingOverwriteAliasDoesNotExistShouldCreateAlias(t *testing.T) {
	putAlias := false
	createIndexCalled := false
	destinationClient := &mock.ElasticClientStub{
		DoesAliasExistCalled: func(_ string) bool {
			return false
		},
		DoesIndexExistCalled: func(_ string) bool {
			return true
		},
		CreateIndexWithMappingCalled: func(_ string, _ *bytes.Buffer) error {
			createIndexCalled = true
			return nil
		},
		PutAliasCalled: func(_ string, _ string) error {
			putAlias = true
			return nil
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, true, false)
	require.NoError(t, err)
	require.True(t, putAlias)
	require.False(t, createIndexCalled)
}

func testCopyMappingOverwriteAliasExistsIndexDoesNotExistShouldCreateIndex(t *testing.T) {
	putAlias := false
	createIndexCalled := false
	destinationClient := &mock.ElasticClientStub{
		DoesAliasExistCalled: func(_ string) bool {
			return true
		},
		DoesIndexExistCalled: func(_ string) bool {
			return false
		},
		CreateIndexWithMappingCalled: func(_ string, _ *bytes.Buffer) error {
			createIndexCalled = true
			return nil
		},
		PutAliasCalled: func(_ string, _ string) error {
			putAlias = true
			return nil
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, true, false)
	require.NoError(t, err)
	require.False(t, putAlias)
	require.True(t, createIndexCalled)
}

func testCopyMappingOverwriteAliasAndIndexExistShouldCreate(t *testing.T) {
	putAlias := false
	createIndexCalled := false
	destinationClient := &mock.ElasticClientStub{
		DoesAliasExistCalled: func(_ string) bool {
			return false
		},
		DoesIndexExistCalled: func(_ string) bool {
			return false
		},
		CreateIndexWithMappingCalled: func(_ string, _ *bytes.Buffer) error {
			createIndexCalled = true
			return nil
		},
		PutAliasCalled: func(_ string, _ string) error {
			putAlias = true
			return nil
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, true, false)
	require.NoError(t, err)
	require.True(t, putAlias)
	require.True(t, createIndexCalled)
}

func testSkipMappings(t *testing.T) {
	called := false
	destinationClient := &mock.ElasticClientStub{
		DoesAliasExistCalled: func(_ string) bool {
			called = true
			return false
		},
		DoesIndexExistCalled: func(_ string) bool {
			called = true
			return false
		},
		CreateIndexWithMappingCalled: func(_ string, _ *bytes.Buffer) error {
			called = true
			return nil
		},
		PutAliasCalled: func(_ string, _ string) error {
			called = true
			return nil
		},
	}

	r, _ := newReindexer(&mock.ElasticClientStub{}, destinationClient, []string{"index"})

	err := r.copyMappingIfNecessary(testIndex, false, true)
	require.NoError(t, err)
	require.False(t, called)
}
