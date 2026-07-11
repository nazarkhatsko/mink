package config

import (
	"fmt"
	"os"
	"regexp"

	"github.com/nazarkhatsko/mink/pkg/driver"
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

	instanceDrivers := make(map[string]driver.Driver, len(cfg.Instances))
	for name, instance := range cfg.Instances {
		d, err := driver.Get(instance.Driver)
		if err != nil {
			return fmt.Errorf("instance %q: %w", name, err)
		}
		if len(instance.Methods) == 0 {
			return fmt.Errorf("instance %q: methods is required", name)
		}
		driverMethods := stringSet(d.Methods())
		for _, m := range instance.Methods {
			if !driverMethods[m] {
				return fmt.Errorf("instance %q: method %q is not supported by driver %q (available: %v)", name, m, instance.Driver, d.Methods())
			}
		}
		instanceDrivers[name] = d
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
			instance, ok := cfg.Instances[action.Instance]
			if !ok {
				return fmt.Errorf("action %q: instance %q not found", action.ID, action.Instance)
			}
			if action.Method == "" {
				return fmt.Errorf("action %q: method is required", action.ID)
			}

			instanceMethods := stringSet(instance.Methods)
			if !instanceMethods[action.Method] {
				return fmt.Errorf("action %q: method %q not declared in instance %q methods %v", action.ID, action.Method, action.Instance, instance.Methods)
			}

			d := instanceDrivers[action.Instance]
			allowed := d.Options(action.Method)
			allowedNames := make(map[string]bool, len(allowed))
			for _, opt := range allowed {
				allowedNames[opt.Name] = true
			}
			for k := range action.ExecuteWith {
				if !allowedNames[k] {
					return fmt.Errorf("action %q: execute_with key %q is not valid for method %q (allowed: %v)", action.ID, k, action.Method, optionNames(allowed))
				}
			}
			for _, opt := range allowed {
				if !opt.Required {
					continue
				}
				_, inConfig := instance.Config[opt.Name]
				_, inExecuteWith := action.ExecuteWith[opt.Name]
				if !inConfig && !inExecuteWith {
					return fmt.Errorf("action %q: execute_with is missing required key %q for method %q", action.ID, opt.Name, action.Method)
				}
			}
		}
	}
	return nil
}

func stringSet(values []string) map[string]bool {
	set := make(map[string]bool, len(values))
	for _, v := range values {
		set[v] = true
	}
	return set
}

func optionNames(options []driver.Option) []string {
	names := make([]string, len(options))
	for i, opt := range options {
		names[i] = opt.Name
	}
	return names
}
