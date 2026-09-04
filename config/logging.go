package config

import "github.com/go-think/think/support/env"

// LoggingConfig represents application logging configuration
type LoggingConfig struct {
	Default  string                      `config:"default"`
	Channels map[string]LogChannelConfig `config:"channels"`
}

// LogChannelConfig represents configuration for a specific log channel
type LogChannelConfig struct {
	Driver string `config:"driver"`
	Path   string `config:"path"`
	Level  string `config:"level"`
	Days   int    `config:"days"`
}

func loadLoggingConfig() *LoggingConfig {
	defaultChannel := env.Get("LOG_CHANNEL", "stack")
	defaultLevel := env.Get("LOG_LEVEL", "debug")

	return &LoggingConfig{
		Default: defaultChannel,
		Channels: map[string]LogChannelConfig{
			"stack": {
				Driver: "stack",
				Level:  defaultLevel,
			},
			"single": {
				Driver: "single",
				Path:   "storage/logs/think.log",
				Level:  defaultLevel,
			},
			"daily": {
				Driver: "daily",
				Path:   "storage/logs/think.log",
				Level:  defaultLevel,
				Days:   14,
			},
			"console": {
				Driver: "console",
				Level:  defaultLevel,
			},
		},
	}
}
