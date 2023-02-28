package config

import "github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"

// ContextFlagsZeroBalanceSysAccChecker holds all flags required for zeroBalanceSystemAccountChecker tool
type ContextFlagsZeroBalanceSysAccChecker struct {
	trieToolsCommon.ContextFlagsConfig
	TokensDirectory string
	Outfile         string
	CrossCheck      bool
}

// GeneralConfig holds general configs required for zeroBalanceSystemAccountChecker tool if cross check is flag is activated
type GeneralConfig struct {
	Config Config `toml:"config"`
}

// Config holds configs required for zeroBalanceSystemAccountChecker tool if cross check is flag is activated
type Config struct {
	ElasticIndexerConfig ElasticIndexerConfig `toml:"elasticIndexerConfig"`
	Gateway              Gateway              `toml:"gateway"`
}

// ElasticIndexerConfig holds the configuration needed for connecting to an Elasticsearch instance
type ElasticIndexerConfig struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// Gateway holds the configuration needed to cross-check address-token trie
type Gateway struct {
	URL string `toml:"url"`
}
