package generate

import (
	"context"
	"fmt"

	"github.com/brianvoe/gofakeit/v7"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "generate" }

func (d *Driver) Describe() driver.Doc {
	return driver.Doc{
		Description: "Generate fake data from a schema",
		Options: []driver.FieldDoc{
			{Name: "schema", Type: "object", Required: true, Description: "Map of field names to faker type strings (e.g. faker.email, faker.uuid) or static values"},
		},
		Output: []driver.FieldDoc{
			{Name: "<field>", Type: "any", Description: "Generated value for each schema field"},
		},
	}
}

func (d *Driver) Execute(_ context.Context, options map[string]any) (driver.Output, error) {
	raw, ok := options["schema"]
	if !ok {
		return nil, fmt.Errorf("generate: schema is required")
	}

	schema, ok := raw.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("generate: schema must be an object")
	}

	result, err := buildObject(schema)
	if err != nil {
		return nil, err
	}

	return driver.Output(result), nil
}

func buildObject(schema map[string]any) (map[string]any, error) {
	out := make(map[string]any, len(schema))
	for key, raw := range schema {
		field, ok := raw.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("generate: field %q must be an object with type and value", key)
		}

		// вкладений обʼєкт без type/value — рекурсія
		_, hasType := field["type"]
		_, hasValue := field["value"]
		if !hasType && !hasValue {
			nested, err := buildObject(field)
			if err != nil {
				return nil, fmt.Errorf("generate: field %q: %w", key, err)
			}
			out[key] = nested
			continue
		}

		val, err := generateField(key, field)
		if err != nil {
			return nil, err
		}
		out[key] = val
	}
	return out, nil
}

func generateField(key string, field map[string]any) (any, error) {
	typ, _ := field["type"].(string)
	value, _ := field["value"].(string)

	raw := fakerValue(value)

	switch typ {
	case "string":
		return fmt.Sprintf("%v", raw), nil
	case "int":
		switch v := raw.(type) {
		case int:
			return int64(v), nil
		case int64:
			return v, nil
		case float64:
			return int64(v), nil
		default:
			return nil, fmt.Errorf("generate: field %q: cannot cast %T to int", key, raw)
		}
	case "float":
		switch v := raw.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		default:
			return nil, fmt.Errorf("generate: field %q: cannot cast %T to float", key, raw)
		}
	case "bool":
		b, ok := raw.(bool)
		if !ok {
			return nil, fmt.Errorf("generate: field %q: cannot cast %T to bool", key, raw)
		}
		return b, nil
	default:
		return nil, fmt.Errorf("generate: field %q: unknown type %q", key, typ)
	}
}

func fakerValue(value string) any {
	switch value {
	case "faker.email":
		return gofakeit.Email()
	case "faker.name":
		return gofakeit.Name()
	case "faker.password":
		return gofakeit.Password(true, true, true, true, false, 12)
	case "faker.uuid":
		return gofakeit.UUID()
	case "faker.phone":
		return gofakeit.Phone()
	case "faker.age":
		return gofakeit.Number(18, 80)
	default:
		return value // статичне значення
	}
}
