package main

import (
	"fmt"
	"io"
	"os"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
)

const (
	defaultLogsPath = "logs"
	logFilePrefix   = "accounts-exporter"
)

var log = logger.GetOrCreate("main")

func initializeLogger(logLevel string) (io.Closer, error) {
	currentDirectory, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	err = logger.SetLogLevel(logLevel)
	if err != nil {
		return nil, err
	}

	logger.ToggleLoggerName(true)

	fileLogging, err := file.NewFileLogging(file.ArgsFileLogging{
		WorkingDir:      currentDirectory,
		DefaultLogsPath: defaultLogsPath,
		LogFilePrefix:   logFilePrefix,
	})
	if err != nil {
		return nil, fmt.Errorf("%w creating a log file", err)
	}

	return fileLogging, nil
}
