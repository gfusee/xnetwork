package path

import "errors"

var (
	errInvalidNumOfChainIDsFound = errors.New("invalid number of chain ID directories found")
	errInvalidNumOfShardIDsFound = errors.New("invalid number of shard ID directories found in Static directory")
	errNoEpochDirectoryFound     = errors.New("no epoch directory found")
	errMissingStaticDirectory    = errors.New("missing Static directory")
	errInvalidShardIDDirectory   = errors.New("invalid shard ID directory")
	errDirectoryIsNotEmpty       = errors.New("directory is not empty")
)
