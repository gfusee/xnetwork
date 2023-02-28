package main

import (
	"github.com/multiversx/mx-sdk-go/core"
	"github.com/multiversx/mx-sdk-go/interactors"
)

type skAddress struct {
	secretKey []byte
	address   core.AddressHandler
}

type pemDataProvider struct {
}

func (pdp *pemDataProvider) getPrivateKeyAndAddress(pemFile string) (*skAddress, error) {
	w := interactors.NewWallet()
	privateKey, err := w.LoadPrivateKeyFromPemFile(pemFile)
	if err != nil {
		return nil, err
	}

	address, err := w.GetAddressFromPrivateKey(privateKey)
	if err != nil {
		return nil, err

	}

	return &skAddress{
		secretKey: privateKey,
		address:   address,
	}, nil
}
