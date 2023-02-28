package main

import (
	"context"

	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/data"
)

type proxyProvider interface {
	GetNetworkConfig(ctx context.Context) (*data.NetworkConfig, error)
	SendTransaction(ctx context.Context, tx *data.Transaction) (string, error)
	GetAccount(ctx context.Context, address core.AddressHandler) (*data.Account, error)
}
