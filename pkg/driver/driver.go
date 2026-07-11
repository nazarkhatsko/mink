package driver

import "context"

type Output map[string]any

// Option describes one top-level execute_with key a driver method accepts.
// Required options must appear in instance.Config or the action's
// execute_with (or both) — the engine checks presence, not value validity.
type Option struct {
	Name     string
	Required bool
}

type Driver interface {
	Name() string
	Methods() []string
	Options(method string) []Option
	Execute(ctx context.Context, method string, options map[string]any) (Output, error)
}
