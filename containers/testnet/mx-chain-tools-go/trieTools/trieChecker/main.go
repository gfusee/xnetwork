package main

import (
	"context"
	"encoding/hex"
	"fmt"
	"os"
	"path/filepath"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/state"
	"github.com/multiversx/mx-chain-go/storage"
	"github.com/multiversx/mx-chain-go/trie/keyBuilder"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/urfave/cli"
)

const (
	logFilePrefix  = "trie-checker"
	rootHashLength = 32
	addressLength  = 32
)

func main() {
	app := cli.NewApp()
	app.Name = "Trie checker CLI app"
	app.Usage = "This is the entry point for the tool that checks the trie DB"
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

	return checkTrie(flagsConfig, rootHash)
}

func checkTrie(flags trieToolsCommon.ContextFlagsConfig, mainRootHash []byte) error {
	addressConverter, err := pubkeyConverter.NewBech32PubkeyConverter(addressLength, log)
	if err != nil {
		return err
	}

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

	iteratorChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    make(chan error, 1),
	}
	err = tr.GetAllLeavesOnChannel(iteratorChannels, context.Background(), mainRootHash, keyBuilder.NewKeyBuilder())
	if err != nil {
		return err
	}

	numAccountsOnMainTrie := 0
	numCodeNodes := 0
	dataTriesRootHashes := make(map[string][]byte)
	numDataTriesLeaves := 0
	for kv := range iteratorChannels.LeavesChan {
		numAccountsOnMainTrie++

		userAccount := &state.UserAccountData{}
		errUnmarshal := trieToolsCommon.Marshaller.Unmarshal(userAccount, kv.Value())
		if errUnmarshal != nil {
			// probably a code node
			numCodeNodes++
			continue
		}
		if len(userAccount.RootHash) == 0 {
			continue
		}

		address := addressConverter.Encode(kv.Key())
		dataTriesRootHashes[address] = userAccount.RootHash
	}

	err = common.GetErrorFromChanNonBlocking(iteratorChannels.ErrChan)
	if err != nil {
		return err
	}

	log.Info("parsed main trie",
		"num accounts", numAccountsOnMainTrie,
		"num code nodes", numCodeNodes,
		"num data tries", len(dataTriesRootHashes))

	if len(dataTriesRootHashes) == 0 {
		return nil
	}

	for address, dataRootHash := range dataTriesRootHashes {
		log.Debug("iterating data trie", "address", address, "data trie root hash", dataRootHash)

		dataTrieIteratorChannels := &common.TrieIteratorChannels{
			LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
			ErrChan:    make(chan error, 1),
		}
		errGetAllLeaves := tr.GetAllLeavesOnChannel(dataTrieIteratorChannels, context.Background(), dataRootHash, keyBuilder.NewDisabledKeyBuilder())
		if errGetAllLeaves != nil {
			return errGetAllLeaves
		}

		for range dataTrieIteratorChannels.LeavesChan {
			numDataTriesLeaves++
		}

		err = common.GetErrorFromChanNonBlocking(dataTrieIteratorChannels.ErrChan)
		if err != nil {
			return err
		}
	}

	log.Info("parsed all tries",
		"num accounts", numAccountsOnMainTrie,
		"num code nodes", numCodeNodes,
		"num data tries", len(dataTriesRootHashes),
		"num data tries leaves", numDataTriesLeaves)

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
