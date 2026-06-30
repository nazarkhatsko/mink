package sleep

import (
	"context"
	"fmt"
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "sleep" }

func (d *Driver) Describe() driver.Doc {
	return driver.Doc{
		Description: "Pause flow execution for N milliseconds",
		Options: []driver.FieldDoc{
			{Name: "ms", Type: "int", Required: true, Description: "Duration in milliseconds"},
		},
		Output: []driver.FieldDoc{
			{Name: "slept_ms", Type: "int", Description: "Actual duration slept in milliseconds"},
		},
	}
}

func (d *Driver) Execute(_ context.Context, options map[string]any) (driver.Output, error) {
	ms, ok := options["ms"]
	if !ok {
		return nil, fmt.Errorf("sleep: ms is required")
	}

	var duration time.Duration
	switch v := ms.(type) {
	case int:
		duration = time.Duration(v) * time.Millisecond
	case int64:
		duration = time.Duration(v) * time.Millisecond
	case float64:
		duration = time.Duration(v) * time.Millisecond
	default:
		return nil, fmt.Errorf("sleep: ms must be a number, got %T", ms)
	}

	time.Sleep(duration)
	return driver.Output{"slept_ms": duration.Milliseconds()}, nil
}
