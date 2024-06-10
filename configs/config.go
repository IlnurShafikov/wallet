package configs

import (
	"github.com/caarlos0/env/v10"
	"strconv"
)

type Config struct {
	Port int `env:"PORT" envDefault:"8080"`
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
