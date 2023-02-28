package path

import (
	"errors"
	"fmt"
	"strings"
	"testing"

	"github.com/multiversx/mx-chain-tools-go/dbmerger/mock"
	"github.com/stretchr/testify/assert"
)

func createMockDirectoryStructure() *mock.DirectoryStructure {
	ds := &mock.DirectoryStructure{}
	ds.AddPath("db/1/Epoch_724/Shard_1", true)
	ds.AddPath("db/1/Epoch_724/Shard_2", true)
	ds.AddPath("db/1/Epoch_725/Shard_1", true)
	ds.AddPath("db/1/Epoch_726/Shard_1", true)
	ds.AddPath("db/1/Static/Shard_1", true)

	return ds
}

func TestDirectoryStructure(t *testing.T) {
	ds := createMockDirectoryStructure()

	fmt.Println(ds.String())

	dirs, err := ds.ListDirectory("db/1/Epoch_724")
	assert.Nil(t, err)

	for _, dir := range dirs {
		fmt.Println(dir.Name())
	}
}

func TestParser_ParseDirectory(t *testing.T) {
	t.Parallel()

	t.Run("more chain IDs should error", func(t *testing.T) {
		t.Parallel()

		ds := &mock.DirectoryStructure{}
		ds.AddPath("db/1/Epoch_0", true)
		ds.AddPath("db/T/Epoch_0", true)

		p := NewParser("db")
		p.dirReadHandler = ds.ListDirectory

		err := p.ParseDirectory()
		assert.Equal(t, errInvalidNumOfChainIDsFound, err)
	})
	t.Run("no correct epochs directories found should error", func(t *testing.T) {
		t.Parallel()

		ds := &mock.DirectoryStructure{}
		ds.AddPath("db/1/Static", true)

		p := NewParser("db")
		p.dirReadHandler = ds.ListDirectory

		t.Run("new epoch directory", func(t *testing.T) {
			err := p.ParseDirectory()
			assert.Equal(t, errNoEpochDirectoryFound, err)
		})
		t.Run("incorrect epoch value", func(t *testing.T) {
			ds.AddPath("db/1/random_directory", true)
			ds.AddPath("db/1/Epoch_", true)
			ds.AddPath("db/1/Epoch_A", true)
			err := p.ParseDirectory()
			assert.Equal(t, errNoEpochDirectoryFound, err)
		})
	})
	t.Run("missing static directory should error", func(t *testing.T) {
		t.Parallel()

		ds := &mock.DirectoryStructure{}
		ds.AddPath("db/1/Epoch_0", true)

		p := NewParser("db")
		p.dirReadHandler = ds.ListDirectory

		err := p.ParseDirectory()
		assert.Equal(t, errMissingStaticDirectory, err)
	})
	t.Run("invalid number of shard ID directories found should error", func(t *testing.T) {
		t.Parallel()

		ds := &mock.DirectoryStructure{}
		ds.AddPath("db/1/Epoch_0", true)
		ds.AddPath("db/1/Static", true)

		p := NewParser("db")
		p.dirReadHandler = ds.ListDirectory

		t.Run("no shard ID directory", func(t *testing.T) {
			err := p.ParseDirectory()
			assert.Equal(t, errInvalidNumOfShardIDsFound, err)
		})
		t.Run("2 shard ID directories", func(t *testing.T) {
			ds.AddPath("db/1/Static/1", true)
			ds.AddPath("db/1/Static/2", true)

			err := p.ParseDirectory()
			assert.Equal(t, errInvalidNumOfShardIDsFound, err)
		})
		t.Run("missing shard path part", func(t *testing.T) {
			ds = &mock.DirectoryStructure{}
			ds.AddPath("db/1/Epoch_0", true)
			ds.AddPath("db/1/Static/random", true)
			p.dirReadHandler = ds.ListDirectory

			err := p.ParseDirectory()
			assert.True(t, errors.Is(err, errInvalidShardIDDirectory))
			assert.True(t, strings.Contains(err.Error(), "found path random"))
		})
		t.Run("invalid shard value part", func(t *testing.T) {
			ds = &mock.DirectoryStructure{}
			ds.AddPath("db/1/Epoch_0", true)
			ds.AddPath("db/1/Static/Shard_A", true)
			p.dirReadHandler = ds.ListDirectory

			err := p.ParseDirectory()
			assert.True(t, errors.Is(err, errInvalidShardIDDirectory))
			assert.True(t, strings.Contains(err.Error(), "found path Shard_A"))
		})
	})
	t.Run("should work with continuous epochs", func(t *testing.T) {
		t.Parallel()

		ds := createMockDirectoryStructure()
		p := NewParser("db")
		p.dirReadHandler = ds.ListDirectory

		err := p.ParseDirectory()

		assert.Nil(t, err)
		assert.Equal(t, "1", p.ChainID())
		assert.Equal(t, 726, int(p.HighestEpoch()))
		assert.Equal(t, 724, int(p.LowestContinuousEpoch()))
		assert.Equal(t, uint64(1), p.ShardID())
	})
	t.Run("should work with not continuous epochs", func(t *testing.T) {
		t.Parallel()

		ds := createMockDirectoryStructure()
		ds.AddPath("db/1/Epoch_723/Shard_1", true)
		ds.AddPath("db/1/Epoch_721/Shard_1", true)
		ds.AddPath("db/1/Epoch_1/Shard_1", true)
		ds.AddPath("db/1/Epoch_0/Shard_1", true)

		p := NewParser("db")
		p.dirReadHandler = ds.ListDirectory

		err := p.ParseDirectory()

		assert.Nil(t, err)
		assert.Equal(t, "1", p.ChainID())
		assert.Equal(t, 726, int(p.HighestEpoch()))
		assert.Equal(t, 723, int(p.LowestContinuousEpoch()))
		assert.Equal(t, uint64(1), p.ShardID())
	})
}
