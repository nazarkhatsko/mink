package engine

import (
	"os"
	"strings"

	"github.com/nazarkhatsko/mink/internal/vars"
)

type Context struct {
	resolver *vars.Resolver
}

func newContext(v map[string]string, actions map[string]map[string]any, state map[string]any) *Context {
	envVars := make(map[string]string)
	for _, e := range os.Environ() {
		parts := strings.SplitN(e, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}
	return &Context{
		resolver: &vars.Resolver{
			Vars:    v,
			Actions: actions,
			State:   state,
			Env:     envVars,
		},
	}
}
