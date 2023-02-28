package trieToolsCommon

import (
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path"

	"github.com/multiversx/mx-chain-core-go/core/check"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	nodeFactory "github.com/multiversx/mx-chain-go/cmd/node/factory"
	"github.com/multiversx/mx-chain-go/common"
	commonDisabled "github.com/multiversx/mx-chain-go/common/disabled"
	nodeConfig "github.com/multiversx/mx-chain-go/config"
	"github.com/multiversx/mx-chain-go/epochStart/notifier"
	"github.com/multiversx/mx-chain-go/state"
	stateFactory "github.com/multiversx/mx-chain-go/state/factory"
	disabled2 "github.com/multiversx/mx-chain-go/state/storagePruningManager/disabled"
	"github.com/multiversx/mx-chain-go/storage"
	"github.com/multiversx/mx-chain-go/storage/databaseremover/disabled"
	"github.com/multiversx/mx-chain-go/storage/factory"
	"github.com/multiversx/mx-chain-go/storage/pruning"
	"github.com/multiversx/mx-chain-go/testscommon"
	"github.com/multiversx/mx-chain-go/trie"
	hashesHolder "github.com/multiversx/mx-chain-go/trie/hashesHolder/disabled"
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
	"github.com/multiversx/mx-chain-storage-go/memorydb"
	"github.com/multiversx/mx-chain-storage-go/storageUnit"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon/components"
	"github.com/urfave/cli"
)

var (
	log = logger.GetOrCreate("trie")
)

const (
	defaultLogsPath      = "logs"
	maxTrieLevelInMemory = 5
	maxDirs              = 100
	addressLength        = 32
)

// AttachFileLogger will attach the file logger, using provided flags
func AttachFileLogger(log logger.Logger, logFilePrefix string, flagsConfig ContextFlagsConfig) (nodeFactory.FileLoggingHandler, error) {
	var fileLogging nodeFactory.FileLoggingHandler
	var err error
	if flagsConfig.SaveLogFile {
		fileLogging, err = file.NewFileLogging(file.ArgsFileLogging{
			WorkingDir:      flagsConfig.WorkingDir,
			DefaultLogsPath: defaultLogsPath,
			LogFilePrefix:   logFilePrefix,
		})
		if err != nil {
			return nil, fmt.Errorf("%w creating a log file", err)
		}
	}

	err = logger.SetDisplayByteSlice(logger.ToHex)
	log.LogIfError(err)
	logger.ToggleLoggerName(flagsConfig.EnableLogName)
	logLevelFlagValue := flagsConfig.LogLevel
	err = logger.SetLogLevel(logLevelFlagValue)
	if err != nil {
		return nil, err
	}

	if flagsConfig.DisableAnsiColor {
		err = logger.RemoveLogObserver(os.Stdout)
		if err != nil {
			return nil, err
		}

		err = logger.AddLogObserver(os.Stdout, &logger.PlainFormatter{})
		if err != nil {
			return nil, err
		}
	}
	log.Trace("logger updated", "level", logLevelFlagValue, "disable ANSI color", flagsConfig.DisableAnsiColor)

	return fileLogging, nil
}

// GetMaxDBValue will search in parentDir for all dbs directories and return max db value
func GetMaxDBValue(parentDir string, log logger.Logger) (int, error) {
	contents, err := ioutil.ReadDir(parentDir)
	if err != nil {
		return 0, err
	}

	directories := make([]string, 0)
	for _, c := range contents {
		if !c.IsDir() {
			continue
		}

		_, ok := big.NewInt(0).SetString(c.Name(), 10)
		if !ok {
			log.Debug("DB directory found that will not be taken into account", "name", c.Name())
			continue
		}

		directories = append(directories, c.Name())
	}

	numDirs := 0
	for i := 0; i < maxDirs; i++ {
		expectedDir := fmt.Sprintf("%d", i)
		if !contains(directories, expectedDir) {
			break
		}

		numDirs++
	}

	if numDirs == 0 {
		return 0, fmt.Errorf("missing ordered directories in %s, like 0, 1 and so on", parentDir)
	}
	if numDirs != len(directories) {
		return 0, fmt.Errorf("unordered directories in %s, like 0, 1 and so on", parentDir)
	}

	return numDirs - 1, nil
}

func contains(haystack []string, needle string) bool {
	for _, h := range haystack {
		if h == needle {
			return true
		}
	}

	return false
}

// CreatePruningStorer will create and return a pruning storer using the provided flags
func CreatePruningStorer(flags ContextFlagsConfig, maxDBValue int) (storage.Storer, error) {
	localDbConfig := dbConfig // copy
	localDbConfig.FilePath = path.Join(flags.WorkingDir, flags.DbDir)

	epochsData := pruning.EpochArgs{
		NumOfEpochsToKeep:     uint32(maxDBValue) + 1,
		NumOfActivePersisters: uint32(maxDBValue) + 1,
		StartingEpoch:         uint32(maxDBValue),
	}

	dbPath := path.Join(flags.WorkingDir, flags.DbDir)
	args := pruning.StorerArgs{
		Identifier:                "",
		ShardCoordinator:          testscommon.NewMultiShardsCoordinatorMock(1),
		CacheConf:                 cacheConfig,
		PathManager:               components.NewSimplePathManager(dbPath),
		DbPath:                    "",
		PersisterFactory:          factory.NewPersisterFactory(localDbConfig),
		Notifier:                  notifier.NewManualEpochStartNotifier(),
		OldDataCleanerProvider:    &testscommon.OldDataCleanerProviderStub{},
		CustomDatabaseRemover:     disabled.NewDisabledCustomDatabaseRemover(),
		MaxBatchSize:              45000,
		EpochsData:                epochsData,
		PruningEnabled:            true,
		EnabledDbLookupExtensions: false,
		PersistersTracker:         pruning.NewPersistersTracker(epochsData),
	}

	return pruning.NewTriePruningStorer(args)
}

// CreateStorer will create and return a storer using the provided flags
func CreateStorer(flags ContextFlagsConfig) (storage.Storer, error) {
	localDbConfig := dbConfig // copy
	localDbConfig.FilePath = path.Join(flags.WorkingDir, flags.DbDir)
	dbPath := path.Join(flags.WorkingDir, flags.DbDir)

	dbConf := storageUnit.DBConfig{
		FilePath:          dbPath,
		Type:              storageUnit.DBType(dbConfig.Type),
		BatchDelaySeconds: dbConfig.BatchDelaySeconds,
		MaxBatchSize:      dbConfig.MaxBatchSize,
		MaxOpenFiles:      dbConfig.MaxOpenFiles,
	}

	return storageUnit.NewStorageUnitFromConf(cacheConfig, dbConf)
}

// CreateTrie will create and return a trie using the provided flags
func CreateTrie(storer storage.Storer) (common.Trie, error) {
	if check.IfNil(storer) {
		return nil, fmt.Errorf("nil storer provided")
	}
	tsm, err := CreateStorageManager(storer)
	if err != nil {
		return nil, err
	}

	return trie.NewTrie(tsm, Marshaller, Hasher, maxTrieLevelInMemory)
}

// CreateStorageManager creates a new trie storage manager using the given storer
func CreateStorageManager(storer storage.Storer) (common.StorageManager, error) {
	tsmArgs := trie.NewTrieStorageManagerArgs{
		MainStorer:        storer,
		CheckpointsStorer: memorydb.New(),
		Marshalizer:       Marshaller,
		Hasher:            Hasher,
		GeneralConfig: nodeConfig.TrieStorageManagerConfig{
			SnapshotsBufferLen:    10,
			SnapshotsGoroutineNum: 100,
		},
		CheckpointHashesHolder: hashesHolder.NewDisabledCheckpointHashesHolder(),
		IdleProvider:           commonDisabled.NewProcessStatusHandler(),
	}

	options := trie.StorageManagerOptions{
		PruningEnabled:     false,
		SnapshotsEnabled:   false,
		CheckpointsEnabled: false,
	}

	return trie.CreateTrieStorageManager(tsmArgs, options)
}

// NewAccountsAdapter will create a new accounts adapter using provided trie
func NewAccountsAdapter(trie common.Trie) (state.AccountsAdapter, error) {
	accCreator := stateFactory.NewAccountCreator()
	storagePruningManager := disabled2.NewDisabledStoragePruningManager()

	addressConverter, err := pubkeyConverter.NewBech32PubkeyConverter(addressLength, log)
	if err != nil {
		return nil, err
	}

	accountsAdapter, err := state.NewAccountsDB(state.ArgsAccountsDB{
		Trie:                  trie,
		Hasher:                Hasher,
		Marshaller:            Marshaller,
		AccountFactory:        accCreator,
		StoragePruningManager: storagePruningManager,
		ProcessingMode:        common.Normal,
		ProcessStatusHandler:  commonDisabled.NewProcessStatusHandler(),
		AppStatusHandler:      commonDisabled.NewAppStatusHandler(),
		AddressConverter:      addressConverter,
	})

	return accountsAdapter, err
}

// GetNumTokens will return the number of tokens in the map
func GetNumTokens(addressTokensMap map[string]map[string]struct{}) int {
	numTokensInShard := 0
	for _, tokens := range addressTokensMap {
		for range tokens {
			numTokensInShard++
		}
	}

	return numTokensInShard
}

// GetFlags will return an array of the defined flags
func GetFlags() []cli.Flag {
	return []cli.Flag{
		WorkingDirectory,
		DbDirectory,
		LogLevel,
		DisableAnsiColor,
		LogSaveFile,
		LogWithLoggerName,
		ProfileMode,
		HexRootHash,
	}
}

// GetFlagsConfig returns a populated ContextFlagsConfig
func GetFlagsConfig(ctx *cli.Context) ContextFlagsConfig {
	flagsConfig := ContextFlagsConfig{}

	flagsConfig.WorkingDir = ctx.GlobalString(WorkingDirectory.Name)
	flagsConfig.DbDir = ctx.GlobalString(DbDirectory.Name)
	flagsConfig.LogLevel = ctx.GlobalString(LogLevel.Name)
	flagsConfig.DisableAnsiColor = ctx.GlobalBool(DisableAnsiColor.Name)
	flagsConfig.SaveLogFile = ctx.GlobalBool(LogSaveFile.Name)
	flagsConfig.EnableLogName = ctx.GlobalBool(LogWithLoggerName.Name)
	flagsConfig.EnablePprof = ctx.GlobalBool(ProfileMode.Name)
	flagsConfig.HexRootHash = ctx.GlobalString(HexRootHash.Name)

	return flagsConfig
}
