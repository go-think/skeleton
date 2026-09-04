package config

import "github.com/go-think/think/support/env"

// DatabaseConfig represents database connections configuration
type DatabaseConfig struct {
	Default     string           `config:"default"`
	Connections ConnectionConfig `config:"connections"`
	Redis       RedisConfig      `config:"redis"`
}

type ConnectionConfig struct {
	MySQL MySQLConnection `config:"mysql"`
}

type MySQLConnection struct {
	Driver   string `config:"driver"`
	Host     string `config:"host"`
	Port     int    `config:"port"`
	Database string `config:"database"`
	Username string `config:"username"`
	Password string `config:"password"`
}

type RedisConfig struct {
	Default RedisConnection `config:"default"`
}

type RedisConnection struct {
	Host     string `config:"host"`
	Port     int    `config:"port"`
	Password string `config:"password"`
}

func loadDatabaseConfig() *DatabaseConfig {
	return &DatabaseConfig{
		Default: env.Get("DB_CONNECTION", "mysql"),
		Connections: ConnectionConfig{
			MySQL: MySQLConnection{
				Driver:   "mysql",
				Host:     env.Get("DB_HOST", "127.0.0.1"),
				Port:     env.GetInt("DB_PORT", 3306),
				Database: env.Get("DB_DATABASE", "think"),
				Username: env.Get("DB_USERNAME", "root"),
				Password: env.Get("DB_PASSWORD", ""),
			},
		},
		Redis: RedisConfig{
			Default: RedisConnection{
				Host:     env.Get("REDIS_HOST", "127.0.0.1"),
				Port:     env.GetInt("REDIS_PORT", 6379),
				Password: env.Get("REDIS_PASSWORD", ""),
			},
		},
	}
}
