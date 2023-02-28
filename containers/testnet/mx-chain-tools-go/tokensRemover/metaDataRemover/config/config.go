package config

import "github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"

// ContextFlagsMetaDataRemover is the flags config for meta data remover
type ContextFlagsMetaDataRemover struct {
	trieToolsCommon.ContextFlagsConfig
	Outfile string
	Tokens  string
	Pems    string
}

// Config holds the config for meta data remover tool
type Config struct {
	ProxyUrl                     string `toml:"ProxyUrl"`
	TokensToDeletePerTransaction uint64 `toml:"TokensToDeletePerTransaction"`
	AdditionalGasLimit           uint64 `toml:"AdditionalGasLimit"`
}
