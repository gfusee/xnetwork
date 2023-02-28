package main

import (
	logger "github.com/multiversx/mx-chain-logger-go"
)

var (
	log             = logger.GetOrCreate("main")
	outputFileName  = "output.json"
	outputFilePerms = 0644
)
