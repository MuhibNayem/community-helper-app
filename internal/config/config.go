package config

import (
	"fmt"
	"os"
)

type Config struct {
	HTTPPort string

	Env string
}

func Load() (*Config, error) {
	port := os.Getenv("HTTP_PORT")
	if port == "" {
		port = "8080"
	}

	env := os.Getenv("APP_ENV")
	if env == "" {
		env = "development"
	}

	return &Config{
		HTTPPort: port,
		Env:      env,
	}, nil
}

func (c *Config) Address() string {
	return fmt.Sprintf(":%s", c.HTTPPort)
}
