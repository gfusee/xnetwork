package main

import (
	"fmt"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/common"
)

type pemsDataReader struct {
	pemProvider pemProvider
	fileHandler common.FileHandler
}

func newPemsDataReader(pemProvider pemProvider, fileHandler common.FileHandler) (*pemsDataReader, error) {
	if pemProvider == nil {
		return nil, errNilPemProvider
	}
	if fileHandler == nil {
		return nil, errNilFileHandler
	}

	return &pemsDataReader{
		pemProvider: pemProvider,
		fileHandler: fileHandler,
	}, nil
}

func (pdr *pemsDataReader) readPemsData(pemsFile string) (map[uint32]*skAddress, error) {
	workingDir, err := pdr.fileHandler.Getwd()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(workingDir, pemsFile)
	contents, err := pdr.fileHandler.ReadDir(fullPath)
	if err != nil {
		return nil, err
	}

	shardPemDataMap := make(map[uint32]*skAddress)
	for _, file := range contents {
		if file.IsDir() {
			continue
		}

		shardID, err := getShardID(file.Name())
		if err != nil {
			return nil, err
		}

		pemData, err := pdr.pemProvider.getPrivateKeyAndAddress(filepath.Join(fullPath, file.Name()))
		if err != nil {
			return nil, err
		}

		shardPemDataMap[shardID] = pemData
	}

	return shardPemDataMap, nil
}

func getShardID(file string) (uint32, error) {
	shardIDStr := strings.TrimPrefix(file, "shard")
	shardIDStr = strings.TrimSuffix(shardIDStr, ".pem")
	shardID, err := strconv.Atoi(shardIDStr)
	if err != nil {
		return 0, fmt.Errorf("invalid file input name = %s; expected pem file name to be <shardX.pem>, where X = number(e.g. shard0.pem)", file)
	}

	return uint32(shardID), nil
}
