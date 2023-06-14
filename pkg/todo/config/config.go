package config

import (
	"fmt"

	"github.com/caarlos0/env/v8"
)

type Config struct {
	DbConfig
}

type DbConfig struct {
	Dsn string `env:"TODO_DB_DSN"`
}

func (c *DbConfig) Validate() error {
	if c.Dsn == "" {
		return fmt.Errorf("environment variable dsn is missing")
	}
	return nil
}

func ReadConfig() (*Config, error) {

	cfg := Config{}
	if err := env.Parse(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
