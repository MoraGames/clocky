package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		App `yaml:"application"`
		Log `yaml:"logger"`
		Env `yaml:"required_envs"`
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

	Env []string
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load("./config/.env"); err != nil {
		return nil, err
	}
	if err := cfg.ReadConfig("./config/config.yml"); err != nil {
		return nil, err
	}
	if err := cfg.ReadEnv(cfg.Env); err != nil {
		return nil, err
	}

	return cfg, nil
}

func (cfg *Config) ReadConfig(path string) error {
	if err := cleanenv.ReadConfig(path, cfg); err != nil {
		return err
	}
	return nil
}

func (cfg *Config) ReadEnv(mustExist Env) error {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return err
	}

	for _, v := range mustExist {
		_, exist := os.LookupEnv(v)
		if !exist {
			return fmt.Errorf("env variable %s must exist", v)
		}
	}

	return nil
}
