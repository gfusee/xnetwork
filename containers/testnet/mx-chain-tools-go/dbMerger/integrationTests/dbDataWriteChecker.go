package integrationTests

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/multiversx/mx-chain-go/storage"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/stretchr/testify/assert"
)

var log = logger.GetOrCreate("integrationTests")

const keyFormat = "key_%d"
const valueFormat = "value_%d"

type dbDataWriteChecker struct {
	counter uint32
}

// NewDBDataWriteChecker creates a new db write checker instance
func NewDBDataWriteChecker() *dbDataWriteChecker {
	return &dbDataWriteChecker{}
}

// AddDataToDB will add the unique data to the db
func (dataWriteChecker *dbDataWriteChecker) AddDataToDB(db storage.Persister, numData int) {
	for i := 0; i < numData; i++ {
		dataWriteChecker.addKeyValue(db)
	}
}

func (dataWriteChecker *dbDataWriteChecker) addKeyValue(db storage.Persister) {
	counterValue := atomic.AddUint32(&dataWriteChecker.counter, 1)
	key := []byte(fmt.Sprintf(keyFormat, counterValue))
	value := []byte(fmt.Sprintf(valueFormat, counterValue))

	log.Debug("writing (key, value) pair", "key", string(key), "value", string(value))

	err := db.Put(key, value)
	log.LogIfError(err)
}

// CheckDB will check the provided DB that all keys & values previously stored exists
func (dataWriteChecker *dbDataWriteChecker) CheckDB(tb testing.TB, db storage.Persister) {
	counterValue := int(atomic.LoadUint32(&dataWriteChecker.counter))
	for i := 1; i <= counterValue; i++ {
		key := []byte(fmt.Sprintf(keyFormat, i))
		expectedValue := []byte(fmt.Sprintf(valueFormat, i))
		log.Debug("checking key", "key", string(key))
		recoveredValue, err := db.Get(key)

		assert.Nil(tb, err)
		assert.Equal(tb, string(expectedValue), string(recoveredValue))
	}
}
