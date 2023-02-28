package main

import (
	"context"
	"fmt"
	"time"

	"github.com/multiversx/mx-sdk-go/data"
)

type txsSender struct {
	proxy                    proxyProvider
	waitTimeNonceIncremented uint64
}

func (ts *txsSender) send(txs []*data.Transaction, startIdx uint64) error {
	numTxs := uint64(len(txs))
	if startIdx >= numTxs {
		return fmt.Errorf("%w, start index = %d, num txs = %d", errIndexOutOfRange, startIdx, numTxs)
	}

	cfg, err := ts.proxy.GetNetworkConfig(context.Background())
	if err != nil {
		return err
	}

	numTxsSent := 0
	roundDuration := cfg.RoundDuration
	log.Info("found",
		"round duration(ms)", roundDuration,
		"total num of txs in file", numTxs,
		"starting index", startIdx,
		"total num of txs to send", numTxs-startIdx)
	for idx := startIdx; idx < numTxs; idx++ {
		tx := txs[idx]
		err = ts.waitForNonceIncremental(tx.SndAddr, tx.Nonce, ts.waitTimeNonceIncremented)
		if err != nil {
			log.Error("waitForNonceIncremental failed", "tx index", idx, "error", err)
			return err
		}

		hash, err := ts.proxy.SendTransaction(context.Background(), tx)
		if err != nil {
			log.Error("failed to send tx", "tx index", idx, "error", err)
			return err
		}

		numTxsSent++
		log.Info("sent transaction",
			"tx hash", hash,
			"current tx index:", idx,
			"sender nonce", tx.Nonce,
			"num txs sent", numTxsSent,
			"remaining num of txs", numTxs-idx-1)
		time.Sleep(time.Millisecond * time.Duration(roundDuration))
	}

	log.Info("finished sending txs", "num sent txs", numTxsSent)
	return nil
}

func (ts *txsSender) waitForNonceIncremental(address string, expectedNonce uint64, waitTime uint64) error {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for numRetrials := uint64(0); numRetrials < waitTime; <-ticker.C {
		accountNonce, errNonce := ts.getNonce(address)
		if errNonce == nil && accountNonce == expectedNonce {
			return nil
		}

		log.Warn("waitForNonceIncremental",
			"expected nonce", expectedNonce,
			"actual nonce", accountNonce,
			"num retrials", numRetrials,
			"error trying to get nonce", errNonce)
		numRetrials++
	}

	return fmt.Errorf("waitForNonceIncremental: %w of %d seconds", errMaxRetrialsExceeded, waitTime)
}

func (ts *txsSender) getNonce(address string) (uint64, error) {
	addr, err := data.NewAddressFromBech32String(address)
	if err != nil {
		return 0, err
	}

	account, err := ts.proxy.GetAccount(context.Background(), addr)
	if err != nil {
		return 0, err
	}

	return account.Nonce, nil
}
