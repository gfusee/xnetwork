package main

import (
	"context"

	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

type proxyProvider interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	GetDefaultTransactionArguments(
		ctx context.Context,
		address core.AddressHandler,
		networkConfigs *data.NetworkConfig,
	) (data.ArgCreateTransaction, error)
}

type transactionInteractor interface {
	ApplySignatureAndGenerateTx(cryptoHolder core.CryptoComponentsHolder, arg data.ArgCreateTransaction) (*data.Transaction, error)
}

type pemProvider interface {
	getPrivateKeyAndAddress(pemFile string) (*skAddress, error)
}
