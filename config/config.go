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
		Console struct {
			Writer     string `env-required:"true" yaml:"writer" env:"LOG_WRITER"`
			Type       string `env-required:"true" yaml:"type" env:"LOG_TYPE"`
			TimeFormat string `env-required:"true" yaml:"time-format" env:"LOG_TIME_FORMAT"`
			Level      string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		} `env-required:"true" yaml:"console"`
		File struct {
			Location   string `env-required:"true" yaml:"location"`
			MaxSize    int    `env-required:"true" yaml:"size-rotation"`
			Type       string `env-required:"true" yaml:"type"`
			TimeFormat string `env-required:"true" yaml:"time-format"`
			Level      string `env-required:"true" yaml:"level"`
		} `env-required:"true" yaml:"file"`
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
