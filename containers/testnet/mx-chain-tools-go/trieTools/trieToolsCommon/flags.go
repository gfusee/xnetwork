package trieToolsCommon

import (
	logger "github.com/multiversx/mx-chain-logger-go"
	"github.com/urfave/cli"
)

var (
	// workingDirectory defines a flag for the path for the working directory.
	WorkingDirectory = cli.StringFlag{
		Name:  "working-directory",
		Usage: "This flag specifies the `directory` where the application will use the databases and logs.",
		Value: "",
	}
	// LogLevel defines a flag that specifies the log level to be used
	LogLevel = cli.StringFlag{
		Name: "log-level",
		Usage: "This flag specifies the logger `level(s)`. It can contain multiple comma-separated value. For example" +
			", if set to *:INFO the logs for all packages will have the INFO level. However, if set to *:INFO,api:DEBUG" +
			" the logs for all packages will have the INFO level, excepting the api package which will receive a DEBUG" +
			" log level.",
		Value: "*:" + logger.LogDebug.String(),
	}
	// LogSaveFile is used when the log output needs to be logged in a file
	LogSaveFile = cli.BoolFlag{
		Name:  "log-save",
		Usage: "Boolean option for enabling log saving. If set, it will automatically save all the logs into a file.",
	}
	// LogWithLoggerName is used to enable log correlation elements
	LogWithLoggerName = cli.BoolFlag{
		Name:  "log-logger-name",
		Usage: "Boolean option for logger name in the logs.",
	}
	// DisableAnsiColor defines if the logger subsystem should prevent displaying ANSI colors
	DisableAnsiColor = cli.BoolFlag{
		Name:  "disable-ansi-color",
		Usage: "Boolean option for disabling ANSI colors in the logging system.",
	}
	// ProfileMode defines a flag for profiling the binary
	// If enabled, it will open the pprof routes over the default gin rest webserver.
	// There are several routes that will be available for profiling (profiling can be analyzed with: go tool pprof):
	//  /debug/pprof/ (can be accessed in the browser, will list the available options)
	//  /debug/pprof/goroutine
	//  /debug/pprof/heap
	//  /debug/pprof/threadcreate
	//  /debug/pprof/block
	//  /debug/pprof/mutex
	//  /debug/pprof/profile (CPU profile)
	//  /debug/pprof/trace?seconds=5 (CPU trace) -> being a trace, can be analyzed with: go tool trace
	// Usage: go tool pprof http(s)://ip.of.the.server/debug/pprof/xxxxx
	ProfileMode = cli.BoolFlag{
		Name: "profile-mode",
		Usage: "Boolean option for enabling the profiling mode. If set, the /debug/pprof routes will be available " +
			"on the node for profiling the application.",
	}
	// DbDirectory defines a flag for the db path inside the working directory.
	DbDirectory = cli.StringFlag{
		Name:  "db-directory",
		Usage: "This flag specifies the `directory` where the application will find the trie storage.",
		Value: "db",
	}
	// HexRootHash defines a flag for the trie root hash expressed in hex format
	HexRootHash = cli.StringFlag{
		Name:  "hex-roothash",
		Usage: "This flag specifies the roothash to start the checking from",
		Value: "",
	}
)
