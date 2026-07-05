package driver

import "context"

type Output map[string]any

type Driver interface {
	Name() string
	Execute(ctx context.Context, options map[string]any) (Output, error)
}
