package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		Logger `env-required:"true" yaml:"logger"`
	}

	Logger struct {
		Console struct {
			Writer     string `env-required:"true" yaml:"writer" env:"LOG_WRITER"`
			Type       string `env-required:"true" yaml:"type" env:"LOG_TYPE"`
			TimeFormat string `env-required:"true" yaml:"time-format" env:"LOG_TIME_FORMAT"`
			Level      string `env-required:"true" yaml:"level" env:"LOG_LEVEL"`
		} `env-required:"true" yaml:"console"`
		File struct {
			Location   string `env-required:"true" yaml:"location"`
			MaxSize    int    `yaml:"size-rotation"`
			Type       string `env-required:"true" yaml:"type"`
			TimeFormat string `env-required:"true" yaml:"time-format"`
			Level      string `env-required:"true" yaml:"level"`
		} `yaml:"file"`
	}
)

const (
	DefaultEnvMode = "dev"
)

func ResolveEnvMode(appEnvFlag string) string {
	if appEnvFlag != "" {
		return appEnvFlag
	}
	if envMode, exist := os.LookupEnv("APP_ENV"); exist {
		return envMode
	}
	return DefaultEnvMode
}

func NewConfig(envMode string) (*Config, error) {
	cfg := &Config{}

	if _, err := os.Stat(fmt.Sprintf("./config/.env.%s", envMode)); err == nil {
		if err := godotenv.Load(fmt.Sprintf("./config/.env.%s", envMode)); err != nil {
			return nil, err
		}
	} else if _, err := os.Stat("./config/.env"); err == nil {
		if err := godotenv.Load("./config/.env"); err != nil {
			return nil, err
		}
	}
	if err := cfg.ReadConfig("./config/config.yml"); err != nil {
		return nil, err
	}
	if err := cfg.ReadEnv([]string{"TELEGRAM_API_TOKEN", "TELEGRAM_ADMIN_ID"}); err != nil {
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

func (cfg *Config) ReadEnv(mustExist []string) error {
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
