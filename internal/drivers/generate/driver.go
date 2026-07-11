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

func (d *Driver) Methods() []string { return []string{"object"} }

func (d *Driver) Options(method string) []driver.Option {
	return []driver.Option{{Name: "schema", Required: true}}
}

func (d *Driver) Execute(_ context.Context, method string, options map[string]any) (driver.Output, error) {
	if method != "object" {
		return nil, driver.NewError(driver.ErrConfig, "generate: unknown method %q", method)
	}

	raw, ok := options["schema"]
	if !ok {
		return nil, driver.NewError(driver.ErrConfig, "generate: schema is required")
	}

	schema, ok := raw.(map[string]any)
	if !ok {
		return nil, driver.NewError(driver.ErrConfig, "generate: schema must be an object")
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
			return nil, driver.NewError(driver.ErrConfig, "generate: field %q must be an object with type and value", key)
		}

		// вкладений обʼєкт без type/value — рекурсія
		_, hasType := field["type"]
		_, hasValue := field["value"]
		if !hasType && !hasValue {
			nested, err := buildObject(field)
			if err != nil {
				return nil, driver.Wrap(driver.ErrConfig, err, "generate: field %q: %v", key, err)
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
			return nil, driver.NewError(driver.ErrConfig, "generate: field %q: cannot cast %T to int", key, raw)
		}
	case "float":
		switch v := raw.(type) {
		case float64:
			return v, nil
		case int:
			return float64(v), nil
		default:
			return nil, driver.NewError(driver.ErrConfig, "generate: field %q: cannot cast %T to float", key, raw)
		}
	case "bool":
		b, ok := raw.(bool)
		if !ok {
			return nil, driver.NewError(driver.ErrConfig, "generate: field %q: cannot cast %T to bool", key, raw)
		}
		return b, nil
	default:
		return nil, driver.NewError(driver.ErrConfig, "generate: field %q: unknown type %q", key, typ)
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
