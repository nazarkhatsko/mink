package config

import (
	"fmt"
	"os"
	"regexp"

	"gopkg.in/yaml.v3"
)

var idPattern = regexp.MustCompile(`^[a-z][a-z0-9_]*$`)

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
	seenFlowIDs := make(map[string]bool)
	for _, flow := range cfg.Flows {
		if flow.ID == "" {
			return fmt.Errorf("flow id is required")
		}
		if !idPattern.MatchString(flow.ID) {
			return fmt.Errorf("flow id %q must be snake_case (lowercase letters, digits, underscores, starting with a letter)", flow.ID)
		}
		if seenFlowIDs[flow.ID] {
			return fmt.Errorf("duplicate flow id %q", flow.ID)
		}
		seenFlowIDs[flow.ID] = true
		for _, action := range flow.Actions {
			if action.ID == "" {
				return fmt.Errorf("action id is required in flow %q", flow.ID)
			}
			if !idPattern.MatchString(action.ID) {
				return fmt.Errorf("action id %q in flow %q must be snake_case (lowercase letters, digits, underscores, starting with a letter)", action.ID, flow.ID)
			}
			if action.Instance == "" {
				return fmt.Errorf("action %q: instance is required", action.ID)
			}
			if _, ok := cfg.Instances[action.Instance]; !ok {
				return fmt.Errorf("action %q: instance %q not found", action.ID, action.Instance)
			}
		}
	}
	return nil
}
