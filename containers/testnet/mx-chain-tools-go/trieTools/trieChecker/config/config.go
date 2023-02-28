package config

// ContextFlagsConfig the configuration for flags
type ContextFlagsConfig struct {
	WorkingDir       string
	DbDir            string
	LogLevel         string
	DisableAnsiColor bool
	SaveLogFile      bool
	EnableLogName    bool
	EnablePprof      bool
	HexRootHash      string
}
