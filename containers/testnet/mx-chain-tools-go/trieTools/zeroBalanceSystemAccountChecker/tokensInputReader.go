package main

import (
	"errors"
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/multiversx/mx-chain-core-go/marshal"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/common"
)

const (
	shardFilePrefix = "shard"
	shardFileSuffix = ".json"
)

type addressTokensMapFileReader struct {
	fileHandler common.FileHandler
	marshaller  marshal.Marshalizer
}

func newAddressTokensMapFileReader(
	fileHandler common.FileHandler,
	marshaller marshal.Marshalizer,
) (*addressTokensMapFileReader, error) {
	if fileHandler == nil {
		return nil, errors.New("nil file handler provided")
	}
	if marshaller == nil {
		return nil, errors.New("nil marshaller provided")
	}

	return &addressTokensMapFileReader{
		fileHandler: fileHandler,
		marshaller:  marshaller,
	}, nil
}

func (atr *addressTokensMapFileReader) readTokensWithNonce(tokensDir string) (trieToolsCommon.AddressTokensMap, map[uint32]trieToolsCommon.AddressTokensMap, error) {
	workingDir, err := atr.fileHandler.Getwd()
	if err != nil {
		return nil, nil, err
	}

	fullPath := filepath.Join(workingDir, tokensDir)
	contents, err := atr.fileHandler.ReadDir(fullPath)
	if err != nil {
		return nil, nil, err
	}

	globalAddressTokensMap := trieToolsCommon.NewAddressTokensMap()
	shardAddressTokensMap := make(map[uint32]trieToolsCommon.AddressTokensMap)
	for _, file := range contents {
		if file.IsDir() {
			continue
		}

		shardID, err := getShardID(file.Name())
		if err != nil {
			return nil, nil, err
		}

		addressTokensMapInCurrFile, err := atr.getFileContent(filepath.Join(fullPath, file.Name()))
		if err != nil {
			return nil, nil, err
		}

		shardAddressTokensMap[shardID] = addressTokensMapInCurrFile.Clone()
		merge(globalAddressTokensMap, addressTokensMapInCurrFile)

		log.Info("read data from",
			"file", file.Name(),
			"shard", shardID,
			"num tokens in shard", shardAddressTokensMap[shardID].NumTokens(),
			"num addresses in shard", shardAddressTokensMap[shardID].NumAddresses(),
			"total num addresses in all shards", globalAddressTokensMap.NumAddresses())
	}

	return globalAddressTokensMap, shardAddressTokensMap, nil
}

func getShardID(file string) (uint32, error) {
	shardIDStr := strings.TrimPrefix(file, shardFilePrefix)
	shardIDStr = strings.TrimSuffix(shardIDStr, shardFileSuffix)
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid file input name: %s; expected tokens shard file name to be <%sX%s>, where X = number(e.g. %s0%s)",
			file, shardFilePrefix, shardFileSuffix, shardFilePrefix, shardFileSuffix)
	}

	return uint32(shardID), nil
}

func (atr *addressTokensMapFileReader) getFileContent(file string) (trieToolsCommon.AddressTokensMap, error) {
	jsonFile, err := atr.fileHandler.Open(file)
	if err != nil {
		return nil, err
	}

	bytesFromJson, err := atr.fileHandler.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	addressTokensMapInCurrFile := make(map[string]map[string]struct{})
	err = atr.marshaller.Unmarshal(&addressTokensMapInCurrFile, bytesFromJson)
	if err != nil {
		return nil, err
	}

	ret := trieToolsCommon.NewAddressTokensMap()
	for address, tokens := range addressTokensMapInCurrFile {
		tokensWithNonce := getTokensWithNonce(tokens)
		ret.Add(address, tokensWithNonce)
	}

	return ret, nil
}

func getTokensWithNonce(tokens map[string]struct{}) map[string]struct{} {
	ret := make(map[string]struct{})

	for token := range tokens {
		addTokenInMapIfHasNonce(token, ret)
	}

	return ret
}

func addTokenInMapIfHasNonce(token string, tokens map[string]struct{}) {
	if hasNonce(token) {
		tokens[token] = struct{}{}
	}
}

func hasNonce(token string) bool {
	return strings.Count(token, "-") == 2
}

func merge(dest, src trieToolsCommon.AddressTokensMap) {
	for addressSrc, tokensSrc := range src.GetMapCopy() {
		if dest.HasAddress(addressSrc) {
			log.Debug("same address found in multiple files", "address", addressSrc)
		}

		dest.Add(addressSrc, tokensSrc)
	}
}
