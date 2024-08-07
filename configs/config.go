package configs

import (
	"errors"
	"github.com/caarlos0/env/v10"
	"strconv"
	"time"
)

type Config struct {
	Redis       Redis         `envPrefix:"REDIS_"`
	Port        int           `env:"PORT" envDefault:"8080"`
	Secret      string        `env:"SECRET" envDefault:"runli"`
	LogLevel    string        `env:"LOG_LEVEL"`
	Local       bool          `env:"LOCAL"`
	StorageType string        `env:"STORAGE_TYPE" envDefault:"redis"`
	ExpiredAt   time.Duration `env:"LIFE_TIME" envDefault:"24h"`
}

type Redis struct {
	Address string `env:"ADDRESS" envDefault:"localhost:6379" `
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
