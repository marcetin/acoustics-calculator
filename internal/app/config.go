package app

import (
	"database/sql"
	"fmt"
	"os"

	_ "modernc.org/sqlite"
)

type Config struct {
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Path string
}

func NewConfig() *Config {
	return &Config{
		Database: DatabaseConfig{
			Path: getEnv("DB_PATH", "acoustics.db"),
		},
	}
}

func (d DatabaseConfig) Open() (*sql.DB, error) {
	dsn := fmt.Sprintf("file:%s?cache=shared&_fk=1", d.Path)
	return sql.Open("sqlite", dsn)
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
