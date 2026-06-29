package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

func Load(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config: %w", err)
	}

	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config: %w", err)
	}

	if err := validate(&cfg); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	return &cfg, nil
}

func validate(cfg *Config) error {
	if cfg.Version == "" {
		return fmt.Errorf("version is required")
	}
	for _, flow := range cfg.Flows {
		if flow.Name == "" {
			return fmt.Errorf("flow name is required")
		}
		for _, action := range flow.Actions {
			if action.ID == "" {
				return fmt.Errorf("action id is required in flow %q", flow.Name)
			}
			if action.Use == "" {
				return fmt.Errorf("action %q: use is required", action.ID)
			}
			if _, ok := cfg.Instances[action.Use]; !ok {
				return fmt.Errorf("action %q: instance %q not found", action.ID, action.Use)
			}
		}
	}
	return nil
}
