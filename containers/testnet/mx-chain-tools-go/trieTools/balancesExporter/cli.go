package main

import (
	"fmt"

	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/multiversx/mx-chain-tools-go/trieTools/balancesExporter/common"
	"github.com/multiversx/mx-chain-tools-go/trieTools/balancesExporter/export"
	"github.com/urfave/cli"
)

var (
	helpTemplate = `NAME:
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

	cliFlagDbPath = cli.StringFlag{
		Name:     "db-path",
		Usage:    "The path to a node's database.",
		Required: true,
	}

	cliFlagShard = cli.Uint64Flag{
		Name:     "shard",
		Usage:    "The shard to use for export.",
		Required: true,
	}

	cliFlagNumShards = cli.UintFlag{
		Name:  "num-shards",
		Usage: "Specifies the total number of actual network shards (with the exception of the metachain). Must be 3 for mainnet.",
		Value: 3,
	}

	cliFlagEpoch = cli.Uint64Flag{
		Name:     "epoch",
		Usage:    "The epoch to use for export.",
		Required: true,
	}

	cliFlagLogLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogDebug.String(),
	}

	cliFlagLogSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}

	cliFlagCurrency = cli.StringFlag{
		Name:  "currency",
		Usage: "What balances to export.",
		Value: "EGLD",
	}

	cliFlagCurrencyDecimals = cli.UintFlag{
		Name:  "currency-decimals",
		Usage: "Number of decimals for chosen currency.",
		Value: 18,
	}

	cliFlagExportFormat = cli.StringFlag{
		Name:  "format",
		Usage: fmt.Sprintf("Export format. One of the following: %s", export.AllFormattersNames),
		Value: export.FormatterNamePlainText,
	}

	cliFlagWithContracts = cli.BoolFlag{
		Name:  "with-contracts",
		Usage: "Whether to include contracts in the export.",
	}

	cliFlagWithZero = cli.BoolFlag{
		Name:  "with-zero",
		Usage: "Whether to include accounts with zero balance in the export.",
	}

	cliFlagByProjectedShard = cli.Uint64Flag{
		Name:     "by-projected-shard",
		Usage:    "The projected shard to use for export.",
		Required: false,
	}
)

func getAllCliFlags() []cli.Flag {
	return []cli.Flag{
		cliFlagDbPath,
		cliFlagShard,
		cliFlagNumShards,
		cliFlagEpoch,
		cliFlagLogLevel,
		cliFlagLogSaveFile,
		cliFlagCurrency,
		cliFlagCurrencyDecimals,
		cliFlagExportFormat,
		cliFlagWithContracts,
		cliFlagWithZero,
		cliFlagByProjectedShard,
	}
}

type parsedCliFlags struct {
	dbPath           string
	shard            uint32
	numShards        uint32
	epoch            uint32
	logLevel         string
	saveLogFile      bool
	currency         string
	currencyDecimals uint
	exportFormat     string
	withContracts    bool
	withZero         bool
	byProjectedShard common.OptionalUint32
}

func getParsedCliFlags(ctx *cli.Context) parsedCliFlags {
	return parsedCliFlags{
		dbPath:           ctx.GlobalString(cliFlagDbPath.Name),
		shard:            uint32(ctx.GlobalUint64(cliFlagShard.Name)),
		numShards:        uint32(ctx.GlobalUint(cliFlagNumShards.Name)),
		epoch:            uint32(ctx.GlobalUint64(cliFlagEpoch.Name)),
		logLevel:         ctx.GlobalString(cliFlagLogLevel.Name),
		saveLogFile:      ctx.GlobalBool(cliFlagLogSaveFile.Name),
		currency:         ctx.GlobalString(cliFlagCurrency.Name),
		currencyDecimals: uint(ctx.GlobalUint(cliFlagCurrencyDecimals.Name)),
		exportFormat:     ctx.GlobalString(cliFlagExportFormat.Name),
		withContracts:    ctx.GlobalBool(cliFlagWithContracts.Name),
		withZero:         ctx.GlobalBool(cliFlagWithZero.Name),
		byProjectedShard: common.OptionalUint32{
			Value:    uint32(ctx.GlobalUint64(cliFlagByProjectedShard.Name)),
			HasValue: ctx.GlobalIsSet(cliFlagByProjectedShard.Name),
		},
	}
}
