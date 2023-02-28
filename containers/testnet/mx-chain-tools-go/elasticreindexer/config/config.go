package config

// GeneralConfig holds the entire configuration
type GeneralConfig struct {
	Indexers IndexersConfig `toml:"config"`
}

// IndexersConfig holds the configuration related to indexers
type IndexersConfig struct {
	Input         ElasticInstanceConfig `toml:"input"`
	Output        ElasticInstanceConfig `toml:"output"`
	IndicesConfig IndicesConfig         `toml:"indices"`
}

// ElasticInstanceConfig holds the configuration needed for connecting to an Elasticsearch instance
type ElasticInstanceConfig struct {
	URL      string `toml:"url"`
	Username string `toml:"username"`
	Password string `toml:"password"`
}

// IndicesConfig holds the configuration for the indices
type IndicesConfig struct {
	Indices       []string `toml:"indices-no-timestamp"`
	WithTimestamp struct {
		Enabled              bool     `toml:"enabled"`
		BlockchainStartTime  int64    `toml:"blockchain-start-time"`
		NumParallelWrites    int      `toml:"num-parallel-writes"`
		IndicesWithTimestamp []string `toml:"indices-with-timestamp"`
	} `toml:"with-timestamp"`
}
