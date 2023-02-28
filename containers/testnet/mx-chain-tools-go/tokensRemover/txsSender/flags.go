package main

import (
	"github.com/multiversx/mx-chain-tools-go/tokensRemover/txsSender/config"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/urfave/cli"
)

var (
	input = cli.StringFlag{
		Name:  "input",
		Usage: "This flag specifies the input file; it expects the input to be an array of signed txs",
		Value: "input.json",
	}
	startIndex = cli.Uint64Flag{
		Name:  "start-index",
		Usage: "This flag specifies the starting index from txs input array. This tool will start to send txs starting from this index",
		Value: 0,
	}
)

func getFlags() []cli.Flag {
	return []cli.Flag{
		trieToolsCommon.LogLevel,
		trieToolsCommon.DisableAnsiColor,
		trieToolsCommon.LogSaveFile,
		trieToolsCommon.LogWithLoggerName,
		trieToolsCommon.ProfileMode,
		input,
		startIndex,
	}
}

func getFlagsConfig(ctx *cli.Context) config.ContextFlagsTxsSender {
	flagsConfig := config.ContextFlagsTxsSender{}

	flagsConfig.LogLevel = ctx.GlobalString(trieToolsCommon.LogLevel.Name)
	flagsConfig.SaveLogFile = ctx.GlobalBool(trieToolsCommon.LogSaveFile.Name)
	flagsConfig.EnableLogName = ctx.GlobalBool(trieToolsCommon.LogWithLoggerName.Name)
	flagsConfig.EnablePprof = ctx.GlobalBool(trieToolsCommon.ProfileMode.Name)
	flagsConfig.TxsInput = ctx.GlobalString(input.Name)
	flagsConfig.StartIndex = ctx.GlobalUint64(startIndex.Name)

	return flagsConfig
}
