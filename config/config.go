package config

import (
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App `yaml:"app"`
		Log `yaml:"logger"`
	}

	App struct {
		Name    string `env-required:"true" yaml:"name"    env:"APP_NAME"`
		Version string `env-required:"true" yaml:"version" env:"APP_VERSION"`
	}

	Log struct {
		Type   string `env-required:"true" yaml:"type"  env:"LOG_TYPE"`
		Format string `env-required:"true" yaml:"format" env:"LOG_FORMAT"`
		Level  string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load("./config/.env"); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadConfig("./config/config.yml", cfg); err != nil {
		return nil, err
	}

	if err := cleanenv.ReadEnv(cfg); err != nil {
		return nil, err
	}

	return cfg, nil
}
