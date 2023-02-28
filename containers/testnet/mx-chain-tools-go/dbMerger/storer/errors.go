package storer

import "errors"

var errNilPersister = errors.New("nil persister")
var errInvalidNumberOfPersisters = errors.New("invalid number of persisters")
var errNilComponent = errors.New("nil component")
