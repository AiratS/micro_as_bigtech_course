package config

import (
	"errors"
	"os"
)

const (
	dsnEnvName = "PG_DSN"
)

type PGConfig interface {
	DSN() string
}

type pgConfig struct {
	dsn string
}

func (c *pgConfig) DSN() string {
	return c.dsn
}

func NewPGConfig() (PGConfig, error) {
	dsn := os.Getenv(dsnEnvName)
	if len(dsn) == 0 {
		return nil, errors.New("PG_DSN is not found")
	}

	return &pgConfig{
		dsn: dsn,
	}, nil
}
