package config

import "github.com/go-think/think/support/env"

// AppConfig represents application-level configuration
type AppConfig struct {
	Name     string `config:"name"`
	Env      string `config:"env"`
	Debug    bool   `config:"debug"`
	Url      string `config:"url"`
	Port     string `config:"port"`
	Timezone string `config:"timezone"`
}

func loadAppConfig() *AppConfig {
	return &AppConfig{
		Name:     env.Get("APP_NAME", "Think"),
		Env:      env.Get("APP_ENV", "production"),
		Debug:    env.GetBool("APP_DEBUG", false),
		Url:      env.Get("APP_URL", "http://localhost:8080"),
		Port:     env.Get("APP_PORT", "8080"),
		Timezone: env.Get("APP_TIMEZONE", "UTC"),
	}
}
