package main

import (
	"encoding/hex"

	"github.com/multiversx/mx-sdk-go/builders"
)

const esdtDeleteMetadataFunction = "ESDTDeleteMetadata"

func createShardTxsDataMap(shardTokensMap map[uint32]map[string]struct{}, tokensToDeletePerTx uint64) (map[uint32][][]byte, error) {
	shardTxsDataMap := make(map[uint32][][]byte)
	for shardID, tokens := range shardTokensMap {
		log.Info("creating txs data", "shardID", shardID, "num tokens", len(tokens))
		tokensSorted, err := sortTokensIDByNonce(tokens)
		if err != nil {
			return nil, err
		}

		tokensIntervals := groupTokensByIntervals(tokensSorted)
		tokensSortedByNonces := sortTokenIntervalsByMaxConsecutiveNonces(tokensIntervals)
		tokensInBulks := groupTokenIntervalsInBulks(tokensSortedByNonces, tokensToDeletePerTx)

		txsData, err := createTxsData(tokensInBulks)
		if err != nil {
			return nil, err
		}

		log.Info("created", "num of txs", len(txsData), "shardID", shardID, "num of nonces per tx", tokensToDeletePerTx)
		shardTxsDataMap[shardID] = txsData
	}

	return shardTxsDataMap, nil
}

func createTxsData(bulks [][]*tokenData) ([][]byte, error) {
	txsData := make([][]byte, 0, len(bulks))
	for _, bulk := range bulks {
		txData, err := tokensBulkAsOnData(bulk)
		if err != nil {
			return nil, err
		}

		txsData = append(txsData, txData)
	}

	return txsData, nil
}

func tokensBulkAsOnData(bulk []*tokenData) ([]byte, error) {
	txDataBuilder := builders.NewTxDataBuilder().Function(esdtDeleteMetadataFunction)
	for _, tkData := range bulk {
		tokenIDHex := hex.EncodeToString([]byte(tkData.tokenID))
		txDataBuilder.ArgHexString(tokenIDHex)

		addIntervalsAsOnData(txDataBuilder, tkData.intervals)
	}

	return txDataBuilder.ToDataBytes()
}

func addIntervalsAsOnData(builder builders.TxDataBuilder, intervals []*interval) {
	builder.ArgInt64(int64(len(intervals)))

	for _, interval := range intervals {
		builder.
			ArgInt64(int64(interval.start)).
			ArgInt64(int64(interval.end))
	}
}
