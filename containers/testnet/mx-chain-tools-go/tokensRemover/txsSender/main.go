package main

import (
	"io/ioutil"
	"os"
	"time"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-tools-go/tokensRemover/txsSender/config"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/pelletier/go-toml"
	"github.com/urfave/cli"
)

const (
	logFilePrefix   = "txs-sender"
	tomlFile        = "./config.toml"
	outputFilePerms = 0644
)

func main() {
	app := cli.NewApp()
	app.Name = "Transaction sender tool"
	app.Usage = "This is the entry point for the tool that sends one tx/block from multiple shards"
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

	cfg, err := loadConfig()
	if err != nil {
		return err
	}
	args := blockchain.ArgsProxy{
		ProxyURL:            cfg.ProxyUrl,
		CacheExpirationTime: time.Minute,
		EntityType:          core.Proxy,
	}

	proxy, err := blockchain.NewProxy(args)
	if err != nil {
		return err
	}

	txs, err := readTxsInput(flagsConfig.TxsInput)
	if err != nil {
		return err
	}

	ts := &txsSender{
		proxy:                    proxy,
		waitTimeNonceIncremented: cfg.WaitTimeNonceIncremented,
	}

	return ts.send(txs, flagsConfig.StartIndex)
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
