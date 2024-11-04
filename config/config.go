package config

import (
	"fmt"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/joho/godotenv"
)

type (
	Config struct {
		//Generics
		Logger `env-required:"true" yaml:"logger"`

		//Championships
		Championships `env-required:"true" yaml:"championships"`

		//Sets
		Sets `env-required:"true" yaml:"sets"`

		//Effects
		RandomEffects    `env-required:"true" yaml:"random_effects"`
		InventoryEffects `env-required:"true" yaml:"inventory_effects"`
		ConditionEffects `env-required:"true" yaml:"condition_effects"`

		//Messages
		Summaries `env-required:"true" yaml:"summaries"`
	}

	// Generics
	Logger struct {
		Type     string `env-required:"true" yaml:"type"`
		Format   string `env-required:"true" yaml:"format"`
		Level    string `env-required:"true" yaml:"level"`
		Rotation string `env-required:"true" yaml:"rotation"`
	}

	// Championships
	Championships struct {
		Joining struct {
			Availability string `env-required:"true" yaml:"availability"`
			Window       string `env-required:"true" yaml:"window"`
			Duration     string `yaml:"duration"`
		} `env-required:"true" yaml:"joining"`
		Ending struct {
			Typology string `env-required:"true" yaml:"typology"`
			Points   int    `yaml:"points"`
			Gap      int    `yaml:"gap"`
			Duration string `yaml:"duration"`
			Events   int    `yaml:"events"`
		} `env-required:"true" yaml:"ending"`
	}

	// Sets
	Sets struct {
		Available []string `env-required:"true" yaml:"available"`
		Rotation  struct {
			Enabled   *bool  `env-required:"true" yaml:"enabled"`
			AmountMin int    `yaml:"amount_min"`
			AmountMax int    `yaml:"amount_max"`
			Condition string `yaml:"condition"`
			Events    int    `yaml:"events"`
			Duration  string `yaml:"duration"`
			Hints     struct {
				Enabled     *bool `env-required:"true" yaml:"enabled"`
				Amount      int   `yaml:"amount"`
				Constraints []struct {
					Property string `env-required:"true" yaml:"property"`
					Value    int    `env-required:"true" yaml:"value"`
					Type     string `env-required:"true" yaml:"type"`
				} `yaml:"constraints"`
			} `yaml:"hints"`
		} `env-required:"true" yaml:"rotation"`
	}

	// Effects
	RandomEffects struct {
		Enabled   *bool `env-required:"true" yaml:"enabled"`
		AmountMin int   `yaml:"amount_min"`
		AmountMax int   `yaml:"amount_max"`
	}
	InventoryEffects struct {
		Storage struct {
			Enabled       *bool `env-required:"true" yaml:"enabled"`
			InventorySize int   `yaml:"inventory_size"`
			StackSize     []struct {
				Effect string
				Amount int
			} `yaml:"stack_size"`
		} `env-required:"true" yaml:"storage"`
		Looting struct {
			Enabled     *bool `env-required:"true" yaml:"enabled"`
			Percentage  int   `yaml:"percentage"`
			Replacement *bool `yaml:"replacement"`
		} `env-required:"true" yaml:"looting"`
		Shop struct {
			//TODO: perché proprio questo enabled mi da errore?
			Enabled     *bool  `env-required:"true" yaml:"enabled"`
			Amount      int    `yaml:"amount"`
			Stocks      int    `yaml:"stocks"`
			UsersStocks int    `yaml:"users_stocks"`
			Refresh     *bool  `yaml:"refresh"`
			RefreshTime string `yaml:"refresh_time"`
		} `env-required:"true" yaml:"shop"`
	}
	ConditionEffects []struct {
		Name     string `yaml:"name"`
		MaxLevel int    `yaml:"max_level"`
	}

	// Messages
	Summaries struct {
		Enabled   []string `env-required:"true" yaml:"enabled"`
		Frequency string   `env-required:"true" yaml:"frequency"`
		Location  string   `env-required:"true" yaml:"location"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}

	if err := godotenv.Load("./config/.env"); err != nil {
		return nil, err
	}
	if err := cfg.ReadConfig("./config/config.yml"); err != nil {
		return nil, err
	}

	mustExists := []string{"TELEGRAM_API_TOKEN", "TELEGRAM_ADMINS_ID"}
	if err := cfg.ReadEnv(mustExists); err != nil {
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

func (cfg *Config) ReadEnv(mustExists []string) error {
	if err := cleanenv.ReadEnv(cfg); err != nil {
		return err
	}

	for _, v := range mustExists {
		_, exist := os.LookupEnv(v)
		if !exist {
			return fmt.Errorf("env variable %s must exist", v)
		}
	}

	return nil
}
