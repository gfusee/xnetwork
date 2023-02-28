package main

import (
	"encoding/hex"
	"testing"

	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/common"
	"github.com/multiversx/mx-sdk-go/data"
	"github.com/stretchr/testify/require"
)

func TestNewPemsDataReader(t *testing.T) {
	t.Parallel()

	t.Run("nil pem data provider, should err", func(t *testing.T) {
		t.Parallel()

		pdr, err := newPemsDataReader(nil, common.NewOSFileHandler())
		require.Nil(t, pdr)
		require.Equal(t, errNilPemProvider, err)
	})

	t.Run("nil file handler, should err", func(t *testing.T) {
		t.Parallel()

		pdr, err := newPemsDataReader(&pemDataProvider{}, nil)
		require.Nil(t, pdr)
		require.Equal(t, errNilFileHandler, err)
	})

	t.Run("should work", func(t *testing.T) {
		t.Parallel()

		pdr, err := newPemsDataReader(&pemDataProvider{}, common.NewOSFileHandler())
		require.NotNil(t, pdr)
		require.Nil(t, err)
	})
}

func TestReadPemsData_RealData(t *testing.T) {
	pemsReader, _ := newPemsDataReader(&pemDataProvider{}, common.NewOSFileHandler())
	pemsData, err := pemsReader.readPemsData("testDataPem")
	require.Nil(t, err)

	addrShard0, err := data.NewAddressFromBech32String("erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th")
	require.Nil(t, err)
	skShard0, err := hex.DecodeString("413f42575f7f26fad3317a778771212fdb80245850981e48b58a4f25e344e8f9")
	require.Nil(t, err)

	addrShard1, err := data.NewAddressFromBech32String("erd1k2s324ww2g0yj38qn2ch2jwctdy8mnfxep94q9arncc6xecg3xaq6mjse8")
	require.Nil(t, err)
	skShard1, err := hex.DecodeString("e253a571ca153dc2aee845819f74bcc9773b0586edead15a94cb7235a5027436")
	require.Nil(t, err)

	expectedPemsData := map[uint32]*skAddress{
		0: {
			secretKey: skShard0,
			address:   addrShard0,
		},
		1: {
			secretKey: skShard1,
			address:   addrShard1,
		},
	}
	require.Equal(t, pemsData, expectedPemsData)
}
