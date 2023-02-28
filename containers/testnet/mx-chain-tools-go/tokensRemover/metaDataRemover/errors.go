package main

import "errors"

var errInvalidTokenFormat = errors.New("invalid token format")

var errCouldNotConvertNonceToBigInt = errors.New("could not convert nonce to big int")

var errNilPemProvider = errors.New("received nil pem provider")

var errNilFileHandler = errors.New("received nil file handler")
