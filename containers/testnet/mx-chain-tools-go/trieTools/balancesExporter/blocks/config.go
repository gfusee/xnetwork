package blocks

import (
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
)

func getCacheConfig() storageUnit.CacheConfig {
	return storageUnit.CacheConfig{
		Type:        "SizeLRU",
		Capacity:    500000,
		SizeInBytes: 314572800, // 300MB
	}
}

func getDbConfig(filePath string) storageUnit.DBConfig {
	return storageUnit.DBConfig{
		FilePath:          filePath,
		Type:              "LvlDBSerial",
		BatchDelaySeconds: 2,
		MaxBatchSize:      45000,
		MaxOpenFiles:      10,
	}
}
