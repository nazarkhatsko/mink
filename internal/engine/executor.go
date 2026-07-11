package engine

import (
	"context"
	"time"

	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

func (e *engine) executeAction(ctx context.Context, execCtx *Context, action config.Action, instance config.Instance) (driver.Output, time.Duration, error) {
	resolved, err := execCtx.resolver.ResolveMap(action.ExecuteWith)
	if err != nil {
		return nil, 0, driver.Wrap(driver.ErrConfig, err, "resolve execute_with: %v", err)
	}

	merged := make(map[string]any, len(instance.Config)+len(resolved))
	for k, v := range instance.Config {
		merged[k] = v
	}
	for k, v := range resolved {
		merged[k] = v
	}

	d, err := driver.Get(instance.Driver)
	if err != nil {
		return nil, 0, driver.Wrap(driver.ErrConfig, err, "%v", err)
	}

	driverCtx := ctx
	if action.Timeout > 0 {
		var cancel context.CancelFunc
		driverCtx, cancel = context.WithTimeout(ctx, time.Duration(action.Timeout)*time.Millisecond)
		defer cancel()
	}

	start := time.Now()
	output, err := d.Execute(driverCtx, action.Method, merged)
	duration := time.Since(start)

	return output, duration, err
}
