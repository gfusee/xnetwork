package main

import (
	"github.com/multiversx/mx-chain-tools-go/trieTools/accountStorageExporter/config"
	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/urfave/cli"
)

var (
	// address defines a flag that specifies the bech32 address of the account to fetch the storage for
	address = cli.StringFlag{
		Name:  "address",
		Usage: "This flag specifies the bech32 address to fetch the storage for",
		Value: "",
	}
)

func getFlags() []cli.Flag {
	return []cli.Flag{
		trieToolsCommon.WorkingDirectory,
		trieToolsCommon.DbDirectory,
		trieToolsCommon.LogLevel,
		trieToolsCommon.DisableAnsiColor,
		trieToolsCommon.LogSaveFile,
		trieToolsCommon.LogWithLoggerName,
		trieToolsCommon.ProfileMode,
		trieToolsCommon.HexRootHash,
		address,
	}
}

func getFlagsConfig(ctx *cli.Context) config.ContextFlagsConfigAddr {
	flagsConfig := config.ContextFlagsConfigAddr{}

	flagsConfig.WorkingDir = ctx.GlobalString(trieToolsCommon.WorkingDirectory.Name)
	flagsConfig.DbDir = ctx.GlobalString(trieToolsCommon.DbDirectory.Name)
	flagsConfig.LogLevel = ctx.GlobalString(trieToolsCommon.LogLevel.Name)
	flagsConfig.DisableAnsiColor = ctx.GlobalBool(trieToolsCommon.DisableAnsiColor.Name)
	flagsConfig.SaveLogFile = ctx.GlobalBool(trieToolsCommon.LogSaveFile.Name)
	flagsConfig.EnableLogName = ctx.GlobalBool(trieToolsCommon.LogWithLoggerName.Name)
	flagsConfig.EnablePprof = ctx.GlobalBool(trieToolsCommon.ProfileMode.Name)
	flagsConfig.HexRootHash = ctx.GlobalString(trieToolsCommon.HexRootHash.Name)
	flagsConfig.Address = ctx.GlobalString(address.Name)

	return flagsConfig
}
