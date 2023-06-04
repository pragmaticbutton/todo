package config

import (
	"fmt"
	"os"
)

type Config struct {
	DbConfig
}

type DbConfig struct {
	Dsn string
}

func (c *DbConfig) Validate() error {
	if c.Dsn == "" {
		return fmt.Errorf("environment variable dsn is missing")
	}
	return nil
}

func ReadConfig() (*Config, error) {

	dbConfig, err := readDbConfig()
	if err != nil {
		return nil, err
	}
	if err := dbConfig.Validate(); err != nil {
		return nil, err
	}

	c := Config{DbConfig: *dbConfig}
	return &c, nil
}

func readDbConfig() (*DbConfig, error) {

	dsn, err := os.LookupEnv("TODO_DB_DSN")
	if !err {
		return nil, fmt.Errorf("environment variable TODO_DB_DSN not found")
	}

	return &DbConfig{Dsn: dsn}, nil
}
