package main

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-go/common"
	"github.com/multiversx/mx-chain-go/state"
	"github.com/multiversx/mx-chain-go/trie/keyBuilder"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-tools-go/trieTools/tokensExporter/config"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
	"github.com/urfave/cli"
)

const (
	logFilePrefix   = "accounts-tokens-exporter"
	rootHashLength  = 32
	addressLength   = 32
	outputFilePerms = 0644
)

func main() {
	app := cli.NewApp()
	app.Name = "Tokens exporter CLI app"
	app.Usage = "This is the entry point for the tool that exports all tokens for a given root hash"
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

	maxDBValue, err := trieToolsCommon.GetMaxDBValue(filepath.Join(flagsConfig.WorkingDir, flagsConfig.DbDir), log)
	if err != nil {
		return err
	}

	log.Info("starting processing trie", "pid", os.Getpid())

	return exportTokens(flagsConfig, rootHash, maxDBValue)
}

func exportTokens(flags config.ContextFlagsTokensExporter, mainRootHash []byte, maxDBValue int) error {
	addressConverter, err := pubkeyConverter.NewBech32PubkeyConverter(addressLength, log)
	if err != nil {
		return err
	}

	db, err := trieToolsCommon.CreatePruningStorer(flags.ContextFlagsConfig, maxDBValue)
	if err != nil {
		return err
	}

	tr, err := trieToolsCommon.CreateTrie(db)
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

	accDb, err := trieToolsCommon.NewAccountsAdapter(tr)
	if err != nil {
		return err
	}

	err = accDb.RecreateTrie(mainRootHash)
	if err != nil {
		return err
	}

	numAccountsOnMainTrie := 0
	addressTokensMap := make(map[string]map[string]struct{})
	for keyValue := range iteratorChannels.LeavesChan {
		address, found := getAddress(keyValue)
		if !found {
			continue
		}

		numAccountsOnMainTrie++

		account, errGetAccount := accDb.GetExistingAccount(address)
		if errGetAccount != nil {
			return errGetAccount
		}

		esdtTokens, errGetESDT := getAllESDTTokens(account, addressConverter)
		if errGetESDT != nil {
			return errGetESDT
		}

		if len(esdtTokens) > 0 {
			encodedAddress := addressConverter.Encode(address)
			addressTokensMap[encodedAddress] = esdtTokens
		}
	}

	err = common.GetErrorFromChanNonBlocking(iteratorChannels.ErrChan)
	if err != nil {
		return err
	}

	encodedSysAccAddress := addressConverter.Encode(vmcommon.SystemAccountAddress)
	log.Info("parsed main trie",
		"num accounts", numAccountsOnMainTrie,
		"num accounts with tokens", len(addressTokensMap),
		"num tokens in all accounts", trieToolsCommon.GetNumTokens(addressTokensMap),
		"num tokens in system account address", len(addressTokensMap[encodedSysAccAddress]))

	_, found := addressTokensMap[encodedSysAccAddress]
	if !found {
		log.Warn(fmt.Sprintf("system account address(%s) not found, input dbs might be incomplete/corrupted", encodedSysAccAddress))
	}

	return saveResult(addressTokensMap, flags.Outfile)
}

func getAddress(kv core.KeyValueHolder) ([]byte, bool) {
	userAccount := &state.UserAccountData{}
	errUnmarshal := trieToolsCommon.Marshaller.Unmarshal(userAccount, kv.Value())
	if errUnmarshal != nil {
		// probably a code node
		return nil, false
	}
	if len(userAccount.RootHash) == 0 {
		return nil, false
	}

	return kv.Key(), true
}

func saveResult(addressTokensMap map[string]map[string]struct{}, outfile string) error {
	jsonBytes, err := json.MarshalIndent(addressTokensMap, "", " ")
	if err != nil {
		return err
	}

	log.Info("writing result in", "file", outfile)
	err = ioutil.WriteFile(outfile, jsonBytes, fs.FileMode(outputFilePerms))
	if err != nil {
		return err
	}

	log.Info("finished exporting address-tokens map")
	return nil
}

func getAllESDTTokens(account vmcommon.AccountHandler, pubKeyConverter core.PubkeyConverter) (map[string]struct{}, error) {
	userAccount, ok := account.(state.UserAccountHandler)
	if !ok {
		return nil, fmt.Errorf("could not convert account to user account, address = %s",
			pubKeyConverter.Encode(account.AddressBytes()))
	}

	allESDTs := make(map[string]struct{})
	if check.IfNil(userAccount.DataTrie()) {
		return allESDTs, nil
	}

	rootHash, err := userAccount.DataTrie().RootHash()
	if err != nil {
		return nil, err
	}

	iteratorChannels := &common.TrieIteratorChannels{
		LeavesChan: make(chan core.KeyValueHolder, common.TrieLeavesChannelDefaultCapacity),
		ErrChan:    make(chan error, 1),
	}
	err = userAccount.DataTrie().GetAllLeavesOnChannel(iteratorChannels, context.Background(), rootHash, keyBuilder.NewKeyBuilder())
	if err != nil {
		return nil, err
	}

	esdtPrefix := []byte(core.ProtectedKeyPrefix + core.ESDTKeyIdentifier)
	for leaf := range iteratorChannels.LeavesChan {
		if !bytes.HasPrefix(leaf.Key(), esdtPrefix) {
			continue
		}

		// TODO: Try to unmarshal it when the new meta data storage model will be live
		tokenKey := leaf.Key()
		lenESDTPrefix := len(esdtPrefix)
		tokenName := getPrettyTokenName(tokenKey[lenESDTPrefix:])

		allESDTs[tokenName] = struct{}{}
	}

	err = common.GetErrorFromChanNonBlocking(iteratorChannels.ErrChan)
	if err != nil {
		return nil, err
	}

	return allESDTs, nil
}

func getPrettyTokenName(tokenName []byte) string {
	token, nonce := common.ExtractTokenIDAndNonceFromTokenStorageKey(tokenName)
	if nonce != 0 {
		tokens := bytes.Split(token, []byte("-"))

		token = append(tokens[0], []byte("-")...)          // ticker-
		token = append(token, tokens[1]...)                // ticker-randSequence
		token = append(token, []byte("-")...)              // ticker-randSequence-
		token = append(token, getPrettyHexNonce(nonce)...) // ticker-randSequence-nonce
	}

	return string(token)
}

func getPrettyHexNonce(nonce uint64) []byte {
	nonceStr := fmt.Sprintf("%x", nonce)
	if len(nonceStr)%2 != 0 {
		nonceStr = "0" + nonceStr
	}

	return []byte(nonceStr)
}
