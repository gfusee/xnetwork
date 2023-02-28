package main

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/multiversx/mx-chain-tools-go/tokensRemover/metaDataRemover/mocks"
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/require"
)

func TestTxsSender_SendTxs(t *testing.T) {
	t.Parallel()

	getAccountCt := 0
	nonce := uint64(4)
	currTxIndex := 0
	roundDuration := int64(100)
	txs := []*data.Transaction{
		{
			Nonce:   nonce,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
		{
			Nonce:   nonce + 1,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
	}
	proxy := &mocks.ProxyStub{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			return &data.NetworkConfig{
				RoundDuration: roundDuration,
			}, nil
		},
		GetAccountCalled: func(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
			defer func() {
				getAccountCt++
				nonce++
			}()

			switch getAccountCt {
			case 0, 1:
				return &data.Account{Nonce: nonce}, nil
			default:
				require.Fail(t, "should not request account anymore")
			}

			return nil, nil
		},
		SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
			require.Equal(t, txs[currTxIndex], tx)

			currTxIndex++
			return fmt.Sprintf("txhash%d", currTxIndex), nil
		},
	}

	ts := txsSender{
		proxy:                    proxy,
		waitTimeNonceIncremented: 60,
	}

	start := time.Now()
	err := ts.send(txs, 0)
	elapsed := time.Since(start)

	require.Nil(t, err)
	blockTimeWaiting := time.Duration(roundDuration) * time.Millisecond
	// 2 txs should take > 2*roundDuration time to send, but < ticker time in waitForNonceIncremental, since no retrial is needed
	require.True(t, elapsed > 2*blockTimeWaiting)
	require.True(t, elapsed < time.Second)
	require.Equal(t, 2, getAccountCt)
	require.Equal(t, 2, currTxIndex)
	require.Equal(t, uint64(6), nonce)
}

func TestTxsSender_SendTxsFromIndex(t *testing.T) {
	t.Parallel()

	getAccountCt := 0
	sendTxsCt := 0
	nonce := uint64(4)
	roundDuration := int64(100)
	txs := []*data.Transaction{
		{
			Nonce:   nonce,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
		{
			Nonce:   nonce + 1,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
	}
	proxy := &mocks.ProxyStub{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			return &data.NetworkConfig{
				RoundDuration: roundDuration,
			}, nil
		},
		GetAccountCalled: func(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
			getAccountCt++
			return &data.Account{Nonce: nonce + 1}, nil
		},
		SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
			sendTxsCt++
			require.Equal(t, txs[1], tx)
			return "txHash", nil
		},
	}

	ts := txsSender{
		proxy:                    proxy,
		waitTimeNonceIncremented: 60,
	}

	start := time.Now()
	err := ts.send(txs, 1)
	elapsed := time.Since(start)

	require.Nil(t, err)
	blockTimeWaiting := time.Duration(roundDuration) * time.Millisecond
	// 1 tx should take > roundDuration time to send, but < 2*roundDuration, since only one out of 2 txs will be sent
	require.True(t, elapsed > blockTimeWaiting)
	require.True(t, elapsed < 2*blockTimeWaiting)
	require.Equal(t, 1, getAccountCt)
	require.Equal(t, 1, sendTxsCt)
}

func TestTxsSender_SendTxsAfterRetrials(t *testing.T) {
	t.Parallel()

	getAccountCt := 0
	nonce := uint64(4)
	currTxIndex := 0
	txs := []*data.Transaction{
		{
			Nonce:   nonce,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
		{
			Nonce:   nonce + 1,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
	}

	roundDuration := int64(100)
	proxy := &mocks.ProxyStub{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			return &data.NetworkConfig{
				RoundDuration: roundDuration,
			}, nil
		},
		GetAccountCalled: func(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
			getAccountCt++

			localErr := fmt.Errorf("error get account %d", getAccountCt)
			switch getAccountCt {
			case 1:
				return nil, localErr
			case 2:
				return nil, localErr
			case 3, 4:
				defer func() {
					nonce++
				}()

				return &data.Account{Nonce: nonce}, nil
			default:
				require.Fail(t, "should not request account anymore")
			}

			return nil, nil
		},
		SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
			require.Equal(t, txs[currTxIndex], tx)

			currTxIndex++
			return fmt.Sprintf("txhash%d", currTxIndex), nil
		},
	}

	ts := txsSender{
		proxy:                    proxy,
		waitTimeNonceIncremented: 60,
	}

	start := time.Now()
	err := ts.send(txs, 0)
	elapsed := time.Since(start)

	require.Nil(t, err)
	blockTimeWaiting := time.Duration(roundDuration) * time.Millisecond
	retrialsTimeWaiting := 2 * time.Second
	waitTimeDuringExecution := 2*blockTimeWaiting + retrialsTimeWaiting
	require.True(t, elapsed > waitTimeDuringExecution)
	require.True(t, elapsed < waitTimeDuringExecution+blockTimeWaiting)
	require.Equal(t, 4, getAccountCt)
	require.Equal(t, 2, currTxIndex)
	require.Equal(t, uint64(6), nonce)
}

func TestTxsSender_SendTxsFailedAfterWaitingForNonce(t *testing.T) {
	t.Parallel()

	getAccountCt := 0
	transactionWasSend := false
	txs := []*data.Transaction{
		{
			Nonce:   0,
			SndAddr: "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th",
		},
	}

	roundDuration := int64(100)
	proxy := &mocks.ProxyStub{
		GetNetworkConfigCalled: func(ctx context.Context) (*data.NetworkConfig, error) {
			return &data.NetworkConfig{
				RoundDuration: roundDuration,
			}, nil
		},
		GetAccountCalled: func(ctx context.Context, address core.AddressHandler) (*data.Account, error) {
			getAccountCt++
			return nil, fmt.Errorf("error get account %d", getAccountCt)
		},
		SendTransactionCalled: func(ctx context.Context, tx *data.Transaction) (string, error) {
			transactionWasSend = true
			return "", nil
		},
	}

	ts := txsSender{
		proxy:                    proxy,
		waitTimeNonceIncremented: 2,
	}

	start := time.Now()
	err := ts.send(txs, 0)
	elapsed := time.Since(start)

	require.NotNil(t, err)
	blockTimeWaiting := time.Duration(roundDuration) * time.Millisecond
	retrialsTimeWaiting := 2 * time.Second
	require.True(t, elapsed > retrialsTimeWaiting)
	require.True(t, elapsed < retrialsTimeWaiting+blockTimeWaiting)
	require.Equal(t, 2, getAccountCt)
	require.False(t, transactionWasSend)
}
