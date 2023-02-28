package main

import (
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/storage"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/urfave/cli"
)

var log = logger.GetOrCreate("trie")

const (
	logFilePrefix  = "trie"
	rootHashLength = 32
)

type StateStatsCollector interface {
	GetStatsForRootHash(rootHash []byte) (common.TriesStatisticsCollector, error)
}

func main() {
	app := cli.NewApp()
	app.Name = "Trie stats CLI app"
	app.Usage = "This is the entry point for the tool that prints stats about the state"
	app.Flags = trieToolsCommon.GetFlags()
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

	log.Info("finished processing trie")
}

func startProcess(c *cli.Context) error {
	flagsConfig := trieToolsCommon.GetFlagsConfig(c)

	_, errLogger := trieToolsCommon.AttachFileLogger(log, logFilePrefix, flagsConfig)
	if errLogger != nil {
		return errLogger
	}

	log.Info("sanity checks...")

	err := logger.SetLogLevel(flagsConfig.LogLevel)
	if err != nil {
		return err
	}

	rootHash, err := hex.DecodeString(flagsConfig.HexRootHash)
	if err != nil {
		return fmt.Errorf("%w when decoding the provided hex root hash", err)
	}
	if len(rootHash) != rootHashLength {
		return fmt.Errorf("wrong root hash length: expected %d, got %d", rootHashLength, len(rootHash))
	}

	log.Info("starting processing trie", "pid", os.Getpid())

	return printTrieStats(flagsConfig, rootHash)
}

func printTrieStats(flags trieToolsCommon.ContextFlagsConfig, mainRootHash []byte) error {
	storer, err := createStorer(flags, log)
	if err != nil {
		return err
	}

	tr, err := trieToolsCommon.CreateTrie(storer)
	if err != nil {
		return err
	}

	defer func() {
		errNotCritical := tr.Close()
		log.LogIfError(errNotCritical)
	}()

	accDb, err := trieToolsCommon.NewAccountsAdapter(tr)
	if err != nil {
		return err
	}

	err = accDb.RecreateTrie(mainRootHash)
	if err != nil {
		return err
	}

	stateStatsCollector, ok := accDb.(StateStatsCollector)
	if !ok {
		return fmt.Errorf("invalid type assertion")
	}

	log.Info("get stats for rootHash", "root hash", mainRootHash)
	stats, err := stateStatsCollector.GetStatsForRootHash(mainRootHash)
	if err != nil {
		return err
	}

	stats.Print()
	return nil
}

func createStorer(flags trieToolsCommon.ContextFlagsConfig, log logger.Logger) (storage.Storer, error) {
	maxDBValue, err := trieToolsCommon.GetMaxDBValue(filepath.Join(flags.WorkingDir, flags.DbDir), log)
	if err == nil {
		return trieToolsCommon.CreatePruningStorer(flags, maxDBValue)
	}

	log.Info("no ordered DBs for a pruning storer operation, will switch to single directory operation...")

	return trieToolsCommon.CreateStorer(flags)
}
