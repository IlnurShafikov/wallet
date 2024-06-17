package configs

import (
	"errors"
	"github.com/caarlos0/env/v10"
	"strconv"
)

type Config struct {
	Port     int    `env:"PORT" envDefault:"8080"`
	Secret   string `env:"SECRET"`
	LogLevel string `env:"LOG_LEVEL"`
}

func (c Config) Validate() error {
	if c.Secret == "" {
		return errors.New("secret is empty")
	}

	return nil
}

func (c Config) GetServerPort() string {
	return ":" + strconv.Itoa(c.Port)
}

func Parse() (*Config, error) {
	config := &Config{}

	if err := env.Parse(config); err != nil {
		return nil, err
	}

	return config, nil
}
