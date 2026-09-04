package config

import (
	"sync"

	"github.com/go-think/think/support/env"
)

// Global strongly-typed configuration instances
var (
	App      *AppConfig
	Session  *SessionConfig
	Database *DatabaseConfig
	Logging  *LoggingConfig
	loadOnce sync.Once
)

// ensureLoaded ensures that configuration instances are populated
func ensureLoaded() {
	loadOnce.Do(func() {
		App = loadAppConfig()
		Session = loadSessionConfig()
		Database = loadDatabaseConfig()
		Logging = loadLoggingConfig()
	})
}

// Reload reloads environment variables and refreshes all typed configuration instances
func Reload(customEnv ...string) map[string]interface{} {
	_ = env.LoadWithPriority(customEnv...)
	App = loadAppConfig()
	Session = loadSessionConfig()
	Database = loadDatabaseConfig()
	Logging = loadLoggingConfig()
	return All()
}

// All returns all application configurations mapped by prefix
func All() map[string]interface{} {
	ensureLoaded()
	return map[string]interface{}{
		"app":      App,
		"session":  Session,
		"database": Database,
		"logging":  Logging,
	}
}
