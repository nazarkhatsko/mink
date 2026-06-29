package engine

import (
	"context"
	"fmt"
	"time"

	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

func (e *engine) executeAction(ctx context.Context, execCtx *Context, action config.Action, instance config.Instance) (driver.Output, time.Duration, error) {
	resolved, err := execCtx.resolver.ResolveMap(action.Options)
	if err != nil {
		return nil, 0, fmt.Errorf("resolve options: %w", err)
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
		return nil, 0, err
	}

	start := time.Now()
	output, err := d.Execute(ctx, merged)
	duration := time.Since(start)

	return output, duration, err
}
