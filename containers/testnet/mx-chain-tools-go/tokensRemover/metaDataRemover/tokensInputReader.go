package main

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
)

func readTokensInput(tokensFile string) (map[uint32]map[string]struct{}, error) {
	workingDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	fullPath := filepath.Join(workingDir, tokensFile)
	jsonFile, err := os.Open(fullPath)
	if err != nil {
		return nil, err
	}

	bytesFromJson, err := ioutil.ReadAll(jsonFile)
	if err != nil {
		return nil, err
	}

	shardTokensMap := make(map[uint32]map[string]struct{})
	err = json.Unmarshal(bytesFromJson, &shardTokensMap)
	if err != nil {
		return nil, err
	}

	log.Info("read from input", "file", tokensFile, "num of shards", len(shardTokensMap), getNumTokens(shardTokensMap))
	return shardTokensMap, nil
}

func getNumTokens(shardTokensMap map[uint32]map[string]struct{}) int {
	numTokensInShard := 0
	for _, tokens := range shardTokensMap {
		numTokensInShard += len(tokens)
	}

	return numTokensInShard
}
