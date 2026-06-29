package engine

import (
	"os"

	"github.com/nazarkhatsko/mink/internal/vars"
)

type Context struct {
	resolver *vars.Resolver
}

func newContext(v map[string]string, actions map[string]map[string]any) *Context {
	return &Context{
		resolver: &vars.Resolver{
			Vars:    v,
			Actions: actions,
			Env:     os.Getenv,
		},
	}
}
