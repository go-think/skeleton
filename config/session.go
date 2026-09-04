package config

import (
	"time"

	"github.com/go-think/think/support/env"
)

// SessionConfig represents session storage configuration
type SessionConfig struct {
	Driver   string        `config:"driver"`
	Cookie   string        `config:"cookie"`
	Lifetime time.Duration `config:"lifetime"`
	Path     string        `config:"path"`
	Domain   string        `config:"domain"`
	Secure   bool          `config:"secure"`
	HttpOnly bool          `config:"httponly"`
}

func loadSessionConfig() *SessionConfig {
	return &SessionConfig{
		Driver:   env.Get("SESSION_DRIVER", "file"),
		Cookie:   env.Get("SESSION_COOKIE", "think_session"),
		Lifetime: time.Duration(env.GetInt("SESSION_LIFETIME", 120)) * time.Minute,
		Path:     env.Get("SESSION_PATH", "/"),
		Domain:   env.Get("SESSION_DOMAIN", ""),
		Secure:   env.GetBool("SESSION_SECURE_COOKIE", false),
		HttpOnly: true,
	}
}
