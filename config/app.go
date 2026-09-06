package config

import (
	"strings"

	"github.com/go-think/think/support/env"
)

// AppConfig represents application-level configuration
type AppConfig struct {
	Name         string   `config:"name"`
	Env          string   `config:"env"`
	Debug        bool     `config:"debug"`
	Key          string   `config:"key"`
	PreviousKeys []string `config:"previous_keys"`
	Url          string   `config:"url"`
	Port         string   `config:"port"`
	Timezone     string   `config:"timezone"`
}

func loadAppConfig() *AppConfig {
	return &AppConfig{
		Name:         env.Get("APP_NAME", "Think"),
		Env:          env.Get("APP_ENV", "production"),
		Debug:        env.GetBool("APP_DEBUG", false),
		Key:          env.Get("APP_KEY", ""),
		PreviousKeys: parseKeyList(env.Get("APP_PREVIOUS_KEYS", "")),
		Url:          env.Get("APP_URL", "http://localhost:8080"),
		Port:         env.Get("APP_PORT", "8080"),
		Timezone:     env.Get("APP_TIMEZONE", "UTC"),
	}
}

// parseKeyList splits a comma-separated key list (for APP_PREVIOUS_KEYS).
func parseKeyList(list string) []string {
	var keys []string
	for _, k := range strings.Split(list, ",") {
		if k = strings.TrimSpace(k); k != "" {
			keys = append(keys, k)
		}
	}
	return keys
}
