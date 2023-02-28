package process

import (
	"errors"

	"github.com/multiversx/mx-chain-tools-go/elasticreindexer/config"
	"github.com/multiversx/mx-chain-tools-go/elasticreindexer/elastic"
)

// CreateReindexer will create the source and destination elastic handlers and create a reindexer based on them
func CreateReindexer(cfg *config.GeneralConfig) (*reindexer, error) {
	if cfg.Indexers.Input.URL == "" {
		return nil, errors.New("empty url for the input cluster")
	}
	if cfg.Indexers.Output.URL == "" {
		return nil, errors.New("empty url for the output cluster")
	}

	sourceElastic, err := elastic.NewElasticClient(cfg.Indexers.Input)
	if err != nil {
		return nil, err
	}

	destinationElastic, err := elastic.NewElasticClient(cfg.Indexers.Output)
	if err != nil {
		return nil, err
	}

	return newReindexer(sourceElastic, destinationElastic, cfg.Indexers.IndicesConfig.Indices)
}
