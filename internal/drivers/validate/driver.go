package validate

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/nazarkhatsko/mink/pkg/driver"
	"github.com/santhosh-tekuri/jsonschema/v6"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "validate" }

func (d *Driver) Execute(_ context.Context, options map[string]any) (driver.Output, error) {
	value, ok := options["value"]
	if !ok {
		return nil, fmt.Errorf("validate: value is required")
	}

	rawSchema, ok := options["schema"]
	if !ok {
		return nil, fmt.Errorf("validate: schema is required")
	}

	compiler := jsonschema.NewCompiler()
	if err := compiler.AddResource("schema.json", rawSchema); err != nil {
		return nil, fmt.Errorf("validate: add schema resource: %w", err)
	}
	schema, err := compiler.Compile("schema.json")
	if err != nil {
		return nil, fmt.Errorf("validate: compile schema: %w", err)
	}

	// normalize value через JSON round-trip щоб отримати json-сумісні типи
	valueBytes, err := json.Marshal(value)
	if err != nil {
		return nil, fmt.Errorf("validate: marshal value: %w", err)
	}
	var normalized any
	if err := json.Unmarshal(valueBytes, &normalized); err != nil {
		return nil, fmt.Errorf("validate: normalize value: %w", err)
	}

	if err := schema.Validate(normalized); err != nil {
		return driver.Output{"valid": false, "error": err.Error()}, &ValidationError{err.Error()}
	}

	return driver.Output{"valid": true}, nil
}

type ValidationError struct {
	msg string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %s", strings.TrimSpace(e.msg))
}
