package config

import (
	"fmt"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type (
	Config struct {
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load("./config/env/.env"); err != nil {
		if err.Error() != "open ./config/env/.env: no such file or directory" {
			return nil, err
		} else {
			logrus.Warn(".env file not found, using env variables")
		}
	}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
