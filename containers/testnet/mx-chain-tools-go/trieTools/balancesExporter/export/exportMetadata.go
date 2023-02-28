package export

type exportMetadata struct {
	ChainID                  string `json:"chainID"`
	ActualShardID            uint32 `json:"actualShardID"`
	Epoch                    uint32 `json:"epoch"`
	BlockNonce               uint64 `json:"blockNonce"`
	BlockRootHash            string `json:"blockRootHash"`
	Format                   string `json:"format"`
	Currency                 string `json:"currency"`
	CurrencyDecimals         uint   `json:"currencyDecimals"`
	WithContracts            bool   `json:"withContracts"`
	WithZero                 bool   `json:"withZero"`
	ByProjectedShardID       uint32 `json:"byProjectedShardID"`
	ByProjectedShardHasValue bool   `json:"byProjectedShardHasValue"`
	NumAccounts              int    `json:"numAccounts"`
}
