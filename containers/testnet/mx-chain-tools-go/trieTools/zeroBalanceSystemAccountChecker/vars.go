package main

import (
	"github.com/multiversx/mx-chain-core-go/marshal"
	logger "github.com/multiversx/mx-chain-logger-go"
)

var (
	log            = logger.GetOrCreate("main")
	jsonMarshaller = &marshal.JsonMarshalizer{}
)
