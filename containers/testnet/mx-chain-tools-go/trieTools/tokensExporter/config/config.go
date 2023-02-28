package config

import "github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"

// ContextFlagsTokensExporter is the flags config for tokens exporter
type ContextFlagsTokensExporter struct {
	trieToolsCommon.ContextFlagsConfig
	Outfile string
}
