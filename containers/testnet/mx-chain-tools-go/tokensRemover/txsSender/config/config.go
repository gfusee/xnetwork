package config

import "github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"

// ContextFlagsTxsSender is the flags config for txs sender tool
type ContextFlagsTxsSender struct {
	trieToolsCommon.ContextFlagsConfig
	TxsInput   string
	StartIndex uint64
}

// Config holds the config for txs sender tool
type Config struct {
	ProxyUrl                 string `toml:"ProxyUrl"`
	WaitTimeNonceIncremented uint64 `toml:"WaitTimeNonceIncremented"`
}
