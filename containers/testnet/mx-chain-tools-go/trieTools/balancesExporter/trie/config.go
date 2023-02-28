package trie

import (
	nodeConfig "github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
)

func getCacheConfig() storageUnit.CacheConfig {
	return storageUnit.CacheConfig{
		Type:        "SizeLRU",
		Capacity:    500000,
		SizeInBytes: 314572800, // 300MB
	}
}

func getDbConfig(filePath string) nodeConfig.DBConfig {
	return nodeConfig.DBConfig{
		FilePath:          filePath,
		Type:              "LvlDBSerial",
		BatchDelaySeconds: 2,
		MaxBatchSize:      45000,
		MaxOpenFiles:      10,
	}
}
