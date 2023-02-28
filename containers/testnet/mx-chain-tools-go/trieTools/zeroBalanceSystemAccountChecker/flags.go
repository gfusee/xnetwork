package main

import (
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/config"
	"github.com/urfave/cli"
)

var (
	tokensDirectory = cli.StringFlag{
		Name:  "tokens-dir",
		Usage: "This flag specifies the `directory` where the application will find all exported tokens from all shards. For each file, it expects the input to be a map<address,tokens>",
		Value: "in",
	}

	outfile = cli.StringFlag{
		Name:  "outfile",
		Usage: "This flag specifies where the output will be stored. It consists of a map<tokens>",
		Value: "output.json",
	}

	crossCheck = cli.BoolFlag{
		Name:  "cross-check",
		Usage: "This flag specifies if a cross check for zero balances result should be done. If set, checks indexer storage using API calls, so it might take a while.",
	}
)

func getFlags() []cli.Flag {
	return []cli.Flag{
		trieToolsCommon.LogLevel,
		trieToolsCommon.DisableAnsiColor,
		trieToolsCommon.LogSaveFile,
		trieToolsCommon.LogWithLoggerName,
		trieToolsCommon.ProfileMode,
		tokensDirectory,
		outfile,
		crossCheck,
	}
}

func getFlagsConfig(ctx *cli.Context) config.ContextFlagsZeroBalanceSysAccChecker {
	flagsConfig := config.ContextFlagsZeroBalanceSysAccChecker{}

	flagsConfig.LogLevel = ctx.GlobalString(trieToolsCommon.LogLevel.Name)
	flagsConfig.SaveLogFile = ctx.GlobalBool(trieToolsCommon.LogSaveFile.Name)
	flagsConfig.EnableLogName = ctx.GlobalBool(trieToolsCommon.LogWithLoggerName.Name)
	flagsConfig.EnablePprof = ctx.GlobalBool(trieToolsCommon.ProfileMode.Name)
	flagsConfig.TokensDirectory = ctx.GlobalString(tokensDirectory.Name)
	flagsConfig.Outfile = ctx.GlobalString(outfile.Name)
	flagsConfig.CrossCheck = ctx.GlobalBool(crossCheck.Name)

	return flagsConfig
}
