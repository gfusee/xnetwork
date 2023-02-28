package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/fs"
	"io/ioutil"
	"os"
	"strconv"
	"time"

	"github.com/multiversx/mx-chain-crypto-go/signing"
	"github.com/multiversx/mx-chain-crypto-go/signing/ed25519"
	"github.com/multiversx/mx-chain-tools-go/tokensRemover/metaDataRemover/config"
	"github.com/multiversx/mx-sdk-go/blockchain"
	"github.com/multiversx/mx-sdk-go/blockchain/cryptoProvider"
	"github.com/multiversx/mx-sdk-go/builders"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/multiversx/mx-sdk-go/interactors"
)

func createShardTxs(
	outFile string,
	cfg *config.Config,
	shardPemsDataMap map[uint32]*skAddress,
	shardTxsDataMap map[uint32][][]byte,
) error {
	if len(shardPemsDataMap) != len(shardTxsDataMap) {
		return fmt.Errorf("provided invalid input; expected number of pem files = number of shards in tokens input; got num shard tokens = %d, num pem files = %d",
			len(shardPemsDataMap), len(shardTxsDataMap))
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

	txBuilder, err := builders.NewTxBuilder(cryptoProvider.NewSigner())
	if err != nil {
		return err
	}

	ti, err := interactors.NewTransactionInteractor(proxy, txBuilder)
	if err != nil {
		return err
	}

	txc, err := newTxCreator(proxy, ti)
	if err != nil {
		return err
	}

	err = createOutputFileIfDoesNotExist(outFile)
	if err != nil {
		return err
	}

	for shardID, txsData := range shardTxsDataMap {
		pemData, found := shardPemsDataMap[shardID]
		if !found {
			return fmt.Errorf("no pem data provided for shard = %d", shardID)
		}

		log.Info("starting to create txs", "shardID", shardID, "num of txs", len(txsData))
		txsInShard, err := txc.createTxs(pemData, txsData, cfg.AdditionalGasLimit)
		if err != nil {
			return err
		}

		file := outFile + "/txsShard" + strconv.Itoa(int(shardID)) + ".json"
		log.Info("saving txs", "shardID", shardID, "file", file)
		err = saveResult(txsInShard, file)
	}

	return nil
}

type txCreator struct {
	proxy         proxyProvider
	txInteractor  transactionInteractor
	networkConfig *data.NetworkConfig
}

// no need to check for nil pointers since this is unexported and only used internally
func newTxCreator(proxy proxyProvider, txInteractor transactionInteractor) (*txCreator, error) {
	netConfigs, err := proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return nil, err
	}

	return &txCreator{
		proxy:         proxy,
		txInteractor:  txInteractor,
		networkConfig: netConfigs,
	}, nil
}

func (tc *txCreator) createTxs(
	pemData *skAddress,
	txsData [][]byte,
	additionalGasLimit uint64,
) ([]*data.Transaction, error) {
	transactionArguments, err := tc.getDefaultTxsArgs(pemData.address)
	if err != nil {
		return nil, err
	}

	suite := ed25519.NewEd25519()
	keyGen := signing.NewKeyGenerator(suite)
	holder, _ := cryptoProvider.NewCryptoComponentsHolder(keyGen, pemData.secretKey)
	txs := make([]*data.Transaction, 0, len(txsData))
	for _, txData := range txsData {
		transactionArguments.Data = txData
		transactionArguments.GasLimit = tc.computeGasLimit(uint64(len(txData))) + additionalGasLimit
		tx, err := tc.txInteractor.ApplySignatureAndGenerateTx(holder, *transactionArguments)
		if err != nil {
			return nil, err
		}

		txs = append(txs, tx)
		transactionArguments.Nonce++
	}

	return txs, nil
}

func (tc *txCreator) getDefaultTxsArgs(address core.AddressHandler) (*data.ArgCreateTransaction, error) {
	transactionArguments, err := tc.proxy.GetDefaultTransactionArguments(context.Background(), address, tc.networkConfig)
	if err != nil {
		return nil, err
	}

	transactionArguments.RcvAddr = address.AddressAsBech32String() // send to self
	transactionArguments.Value = "0"

	return &transactionArguments, nil
}

func (tc *txCreator) computeGasLimit(dataLen uint64) uint64 {
	return tc.networkConfig.MinGasLimit + tc.networkConfig.GasPerDataByte*dataLen
}

func createOutputFileIfDoesNotExist(outFile string) error {
	_, err := os.Stat(outFile)

	if os.IsNotExist(err) {
		errDir := os.MkdirAll(outFile, os.ModePerm)
		if errDir != nil {
			return err
		}
	}

	return nil
}

func saveResult(txs []*data.Transaction, outfile string) error {
	jsonBytes, err := json.MarshalIndent(txs, "", " ")
	if err != nil {
		return err
	}

	log.Info("writing result in", "file", outfile)
	err = ioutil.WriteFile(outfile, jsonBytes, fs.FileMode(outputFilePerms))
	if err != nil {
		return err
	}
	return nil
}
