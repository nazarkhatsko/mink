package time

import (
	"context"
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "time" }

func (d *Driver) Methods() []string { return []string{"wait"} }

func (d *Driver) Options(method string) []driver.Option {
	return []driver.Option{{Name: "ms", Required: true}}
}

func (d *Driver) Execute(_ context.Context, method string, options map[string]any) (driver.Output, error) {
	if method != "wait" {
		return nil, driver.NewError(driver.ErrConfig, "time: unknown method %q", method)
	}

	ms, ok := options["ms"]
	if !ok {
		return nil, driver.NewError(driver.ErrConfig, "time: ms is required")
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
		return nil, driver.NewError(driver.ErrConfig, "time: ms must be a number, got %T", ms)
	}

	time.Sleep(duration)
	return driver.Output{"slept_ms": duration.Milliseconds()}, nil
}
