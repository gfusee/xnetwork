package main

import (
	"io/ioutil"
	"os"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-tools-go/tokensRemover/metaDataRemover/config"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/common"
	"github.com/pelletier/go-toml"
	"github.com/urfave/cli"
)

const (
	logFilePrefix   = "meta-data-remover"
	tomlFile        = "./config.toml"
	outputFilePerms = 0644
)

func main() {
	app := cli.NewApp()
	app.Name = "Tokens exporter CLI app"
	app.Usage = "This is the entry point for the tool that deletes tokens meta-data"
	app.Flags = getFlags()
	app.Authors = []cli.Author{
		{
			Name:  "The MultiversX Team",
			Email: "contact@multiversx.com",
		},
	}

	app.Action = func(c *cli.Context) error {
		return startProcess(c)
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
		return
	}
}

func startProcess(c *cli.Context) error {
	flagsConfig := getFlagsConfig(c)

	_, errLogger := trieToolsCommon.AttachFileLogger(log, logFilePrefix, flagsConfig.ContextFlagsConfig)
	if errLogger != nil {
		return errLogger
	}

	err := logger.SetLogLevel(flagsConfig.LogLevel)
	if err != nil {
		return err
	}

	log.Info("starting processing", "pid", os.Getpid())

	shardTokensMap, err := readTokensInput(flagsConfig.Tokens)
	if err != nil {
		return err
	}

	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	shardTxsDataMap, err := createShardTxsDataMap(shardTokensMap, cfg.TokensToDeletePerTransaction)
	if err != nil {
		return err
	}

	shardPemsDataMap, err := getShardPemsDataMap(flagsConfig.Pems)
	if err != nil {
		return err
	}

	return createShardTxs(flagsConfig.Outfile, cfg, shardPemsDataMap, shardTxsDataMap)
}

func getShardPemsDataMap(pemsFile string) (map[uint32]*skAddress, error) {
	osFileHandler := common.NewOSFileHandler()
	pemsReader, err := newPemsDataReader(&pemDataProvider{}, osFileHandler)
	if err != nil {
		return nil, err
	}

	return pemsReader.readPemsData(pemsFile)
}

func loadConfig() (*config.Config, error) {
	tomlBytes, err := ioutil.ReadFile(tomlFile)
	if err != nil {
		return nil, err
	}

	var cfg config.Config
	err = toml.Unmarshal(tomlBytes, &cfg)
	if err != nil {
		return nil, err
	}

	return &cfg, nil
}
