package sleep

import (
	"context"
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "sleep" }

func (d *Driver) Execute(_ context.Context, options map[string]any) (driver.Output, error) {
	ms, ok := options["ms"]
	if !ok {
		return nil, driver.NewError(driver.ErrConfig, "sleep: ms is required")
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
		return nil, driver.NewError(driver.ErrConfig, "sleep: ms must be a number, got %T", ms)
	}

	time.Sleep(duration)
	return driver.Output{"slept_ms": duration.Milliseconds()}, nil
}
