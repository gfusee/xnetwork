package main

import (
	"errors"
	"fmt"
	"os"
	"strings"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-logger-go/file"
	"github.com/multiversx/mx-chain-tools-go/dbmerger/path"
	"github.com/multiversx/mx-chain-tools-go/dbmerger/storer"
	"github.com/urfave/cli"
)

const sourcePathsDelimiter = ","
const defaultLogsPath = "logs"
const logFilePrefix = "log"

var (
	log = logger.GetOrCreate("main")

	dest = cli.StringFlag{
		Name:  "dest",
		Usage: "This flag specifies the destination path",
		Value: "",
	}
	sources = cli.StringFlag{
		Name:  "sources",
		Usage: `This flag specifies the source paths separated by ",". Example "-sources ` + strings.Join([]string{"path/1", "path/2", "path/3"}, sourcePathsDelimiter) + "\"",
		Value: "",
	}
	logLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogDebug.String(),
	}
	logSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}

	errEmptyPathProvided = errors.New("empty path provided")
)

const helpTemplate = `NAME:
   {{.Name}} - {{.Usage}}
USAGE:
   {{.HelpName}} {{if .VisibleFlags}}[global options]{{end}}
   {{if len .Authors}}
AUTHOR:
   {{range .Authors}}{{ . }}{{end}}
   {{end}}{{if .Commands}}
GLOBAL OPTIONS:
   {{range .VisibleFlags}}{{.}}
   {{end}}
VERSION:
   {{.Version}}
   {{end}}
`

type parsedFlags struct {
	destPath    string
	sourcePaths []string
	logLevel    string
	logSave     bool
}

func main() {
	app := cli.NewApp()
	cli.AppHelpTemplate = helpTemplate
	app.Name = "DB merger tool CLI App"
	app.Version = "v1.0.0"
	app.Usage = "This is the entry point for DB merge tool able to merge 2 or more level DB databases"
	app.Flags = []cli.Flag{
		dest,
		sources,
		logLevel,
		logSaveFile,
	}
	app.Authors = []cli.Author{
		{
			Name:  "The MultiversX Team",
			Email: "contact@multiversx.com",
		},
	}

	app.Action = action

	err := app.Run(os.Args)
	if err != nil {
		log.Error(err.Error())
		os.Exit(1)
	}
}

func action(ctx *cli.Context) {
	flags, err := parseFlags(ctx)
	if err != nil {
		log.Error("cannot process input flags", "error", err)
		return
	}

	err = doAction(flags)
	if err != nil {
		log.Error("cannot perform action", "error", err)
		return
	}

	log.Info("action performed")
}

func parseFlags(ctx *cli.Context) (parsedFlags, error) {
	sourcePaths := ctx.GlobalString(sources.Name)

	flags := parsedFlags{
		destPath:    ctx.GlobalString(dest.Name),
		sourcePaths: strings.Split(sourcePaths, sourcePathsDelimiter),
		logLevel:    ctx.GlobalString(logLevel.Name),
		logSave:     ctx.GlobalBool(logSaveFile.Name),
	}

	// TODO add separate check functions
	if len(flags.destPath) == 0 {
		return parsedFlags{}, fmt.Errorf("%w for `dest` flag", errEmptyPathProvided)
	}
	for idx, src := range flags.sourcePaths {
		if len(src) == 0 {
			return parsedFlags{}, fmt.Errorf("%w for source flag with index %d", errEmptyPathProvided, idx)
		}
	}

	return flags, nil
}

func doAction(flags parsedFlags) error {
	err := processFileLogger(log, flags)
	if err != nil {
		return err
	}

	persisterCreator := storer.NewPersisterCreator()
	args := storer.ArgsFullDBMerger{
		DataMergerInstance:  storer.NewDataMerger(),
		PersisterCreator:    persisterCreator,
		OsOperationsHandler: path.NewOsOperationsHandler(),
	}
	fullDataMerger, err := storer.NewFullDBMerger(args)
	if err != nil {
		return err
	}

	destDB, err := fullDataMerger.MergeDBs(flags.destPath, flags.sourcePaths...)
	if err != nil {
		return err
	}

	return destDB.Close()
}

func processFileLogger(log logger.Logger, flags parsedFlags) error {
	var err error
	if flags.logSave {
		_, err = file.NewFileLogging(file.ArgsFileLogging{
			WorkingDir:      "",
			DefaultLogsPath: defaultLogsPath,
			LogFilePrefix:   logFilePrefix,
		})
		if err != nil {
			return fmt.Errorf("%w creating a log file", err)
		}
	}

	err = logger.SetLogLevel(flags.logLevel)
	if err != nil {
		return err
	}

	log.Trace("logger updated", "level", flags.logLevel)

	return nil
}
