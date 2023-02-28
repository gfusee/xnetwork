package integrationTests

import (
	"testing"

	"github.com/multiversx/mx-chain-tools-go/dbmerger/path"
	"github.com/multiversx/mx-chain-tools-go/dbmerger/storer"
	"github.com/stretchr/testify/assert"
)

func TestFullDBMergerWith3Persisters(t *testing.T) {
	persisterCreator := storer.NewPersisterCreator()
	writeChecker := NewDBDataWriteChecker()

	dbPath1 := createDBAndAddData(t, persisterCreator, writeChecker, 10)
	dbPath2 := createDBAndAddData(t, persisterCreator, writeChecker, 20)
	dbPath3 := createDBAndAddData(t, persisterCreator, writeChecker, 30)
	dbPathDest := t.TempDir()

	args := storer.ArgsFullDBMerger{
		DataMergerInstance:  storer.NewDataMerger(),
		PersisterCreator:    persisterCreator,
		OsOperationsHandler: path.NewOsOperationsHandler(),
	}
	fullDataMerger, err := storer.NewFullDBMerger(args)
	assert.Nil(t, err)

	dest, err := fullDataMerger.MergeDBs(dbPathDest, dbPath1, dbPath2, dbPath3)
	assert.Nil(t, err)

	writeChecker.CheckDB(t, dest)
}

func createDBAndAddData(tb testing.TB, persisterCreator storer.PersisterCreator, writeChecker *dbDataWriteChecker, numData int) string {
	dbPath := tb.TempDir()
	db, err := persisterCreator.CreatePersister(dbPath)
	assert.Nil(tb, err)
	writeChecker.AddDataToDB(db, numData)
	_ = db.Close()

	return dbPath
}
