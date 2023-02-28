package config

import "github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"

// ContextFlagsConfigAddr the configuration for flags
type ContextFlagsConfigAddr struct {
	trieToolsCommon.ContextFlagsConfig
	Address string
}
