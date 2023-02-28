package main

import (
	"errors"
	"fmt"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	vmcommon "github.com/multiversx/mx-chain-vm-common-go"
)

type zeroTokensBalancesExporter struct {
	pubKeyConverter core.PubkeyConverter
}

func newZeroTokensBalancesExporter(pubKeyConverter core.PubkeyConverter) (*zeroTokensBalancesExporter, error) {
	if pubKeyConverter.IsInterfaceNil() {
		return nil, errors.New("nil pub key converter provided")
	}

	return &zeroTokensBalancesExporter{
		pubKeyConverter: pubKeyConverter,
	}, nil
}

func (ztb *zeroTokensBalancesExporter) getExtraTokens(
	globalAddressTokensMap trieToolsCommon.AddressTokensMap,
	shardAddressTokenMap map[uint32]trieToolsCommon.AddressTokensMap,
) (map[string]struct{}, map[uint32]map[string]struct{}, error) {
	systemSCAddress := ztb.pubKeyConverter.Encode(vmcommon.SystemAccountAddress)
	globalExtraTokens, err := getGlobalExtraTokens(globalAddressTokensMap, systemSCAddress)
	if err != nil {
		return nil, nil, err
	}
	log.Info("found", "global num of extra tokens", len(globalExtraTokens))

	shardExtraTokens := make(map[uint32]map[string]struct{})
	for shardID, addressTokensMap := range shardAddressTokenMap {
		if !addressTokensMap.HasAddress(systemSCAddress) {
			return nil, nil, fmt.Errorf("no system account address(%s) found in shard = %v", systemSCAddress, shardID)
		}

		shardTokensInSystemAccAddress := addressTokensMap.GetTokens(systemSCAddress)
		extraTokensInShard := intersection(globalExtraTokens, shardTokensInSystemAccAddress)
		log.Info("found", "shard", shardID, "num tokens in system account", len(shardTokensInSystemAccAddress), "num extra tokens", len(extraTokensInShard))
		shardExtraTokens[shardID] = extraTokensInShard
	}

	if !sanityCheckExtraTokens(shardExtraTokens, globalExtraTokens) {
		return nil, nil, errors.New("sanity check for exported tokens failed")
	}

	return globalExtraTokens, shardExtraTokens, nil
}

func getGlobalExtraTokens(allAddressesTokensMap trieToolsCommon.AddressTokensMap, systemSCAddress string) (map[string]struct{}, error) {
	if !allAddressesTokensMap.HasAddress(systemSCAddress) {
		return nil, fmt.Errorf("no system account address(%s) found", systemSCAddress)
	}

	allTokensInSystemSCAddress := allAddressesTokensMap.GetTokens(systemSCAddress)
	allTokens := getAllTokensWithoutSystemAccount(allAddressesTokensMap, systemSCAddress)
	log.Info("found",
		"global num of tokens in all addresses without system account", len(allTokens),
		"global num of tokens in system sc address", len(allTokensInSystemSCAddress))

	return getExtraTokens(allTokens, allTokensInSystemSCAddress), nil
}

func getAllTokensWithoutSystemAccount(allAddressesTokensMap trieToolsCommon.AddressTokensMap, systemSCAddress string) map[string]struct{} {
	mapCopy := allAddressesTokensMap.Clone()
	mapCopy.Delete(systemSCAddress)

	return mapCopy.GetAllTokens()
}

func getExtraTokens(allTokens, allTokensInSystemSCAddress map[string]struct{}) map[string]struct{} {
	extraTokens := make(map[string]struct{})
	for tokenInSystemSC := range allTokensInSystemSCAddress {
		_, exists := allTokens[tokenInSystemSC]
		if !exists {
			extraTokens[tokenInSystemSC] = struct{}{}
		}
	}

	log.Info("found", "num of sfts/nfts/metaesdts metadata only found in system sc address", len(extraTokens))
	return extraTokens
}

func intersection(globalTokens, shardTokens map[string]struct{}) map[string]struct{} {
	ret := make(map[string]struct{})
	for token := range shardTokens {
		_, found := globalTokens[token]
		if found {
			ret[token] = struct{}{}
		}
	}

	return ret
}

func sanityCheckExtraTokens(shardExtraTokensMap map[uint32]map[string]struct{}, globalExtraTokens map[string]struct{}) bool {
	allMergedExtraTokens := make(map[string]struct{})
	for _, extraTokensInShard := range shardExtraTokensMap {
		for extraToken := range extraTokensInShard {
			allMergedExtraTokens[extraToken] = struct{}{}
		}
	}

	return checkSameMap(allMergedExtraTokens, globalExtraTokens)
}

func checkSameMap(map1, map2 map[string]struct{}) bool {
	if len(map1) != len(map2) {
		return false
	}

	for elemInMap1 := range map1 {
		_, foundInMap2 := map2[elemInMap1]
		if !foundInMap2 {
			return false
		}
	}

	return true
}
