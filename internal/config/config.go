package config

import (
	"encoding/json"
	"os"

	"github.com/zeroaddresss/golang-unisat-monitor/internal/errors"
)

type Config struct {
	Protocol          string            `json:"protocol"`
	Collections       []string          `json:"collections"`
	ApiKeys           []string          `json:"apiKeys"`
	Timeout           int               `json:"timeout"`
	Delay             int               `json:"delay"`
	MaxRetries        int               `json:"maxRetries"`
	Webhooks          map[string]string `json:"webhooks"`
	MonitoringURL     string
	CollectionBaseURL string
	ListingBaseURL    string
}

type validationFunc func(*Config) error

func LoadConfig(filePath string) (*Config, error) {
	configFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, errors.NewConfigError("Error reading config file: %v", err)
	}

	// parse file content into Config struct
	var config Config
	if err = json.Unmarshal(configFile, &config); err != nil {
		return nil, errors.NewConfigError("Error decoding config file: %v", err)
	}

	return &config, nil
}

func ValidateConfig(config *Config) error {
	// validate the config values provided by the user
	validations := map[string]validationFunc{
		"Collections is empty": func(c *Config) error {
			if len(c.Collections) == 0 {
				return errors.NewConfigError("Collections is empty")
			}
			return nil
		},
		"No API keys provided": func(c *Config) error {
			if len(c.ApiKeys) == 0 {
				return errors.NewConfigError("No API keys provided")
			}
			return nil
		},
		"Invalid protocol provided": func(c *Config) error {
			if c.Protocol != "brc20" && c.Protocol != "runes" {
				return errors.NewConfigError("Invalid protocol provided")
			}
			return nil
		},
		"No webhooks provided": func(c *Config) error {
			if len(c.Webhooks) == 0 {
				return errors.NewConfigError("No webhooks provided")
			}
			return nil
		},
		"Delay must be greater than 0": func(c *Config) error {
			if c.Delay <= 0 {
				return errors.NewConfigError("Delay must be greater than 0")
			}
			return nil
		},
		"MaxRetries must be greater than 0": func(c *Config) error {
			if c.MaxRetries <= 0 {
				return errors.NewConfigError("MaxRetries must be greater than 0")
			}
			return nil
		},
		"Timeout must be greater than 0": func(c *Config) error {
			if c.Timeout <= 0 {
				return errors.NewConfigError("Timeout must be greater than 0")
			}
			return nil
		},
		"Timeout must be greater than Delay": func(c *Config) error {
			if c.Timeout <= c.Delay {
				return errors.NewConfigError("Timeout must be greater than Delay")
			}
			return nil
		},
	}

	for _, validate := range validations {
		if err := validate(config); err != nil {
			return err
		}
	}

	return nil
}
