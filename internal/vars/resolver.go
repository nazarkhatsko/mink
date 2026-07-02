package vars

import (
	"fmt"
	"regexp"

	"go.starlark.net/starlark"
)

var exprRe = regexp.MustCompile(`\$\{([^}]+)\}`)

type Resolver struct {
	Vars    map[string]string
	Actions map[string]map[string]any
	State   map[string]any
	Env     map[string]string
}

func (r *Resolver) makeGlobals() starlark.StringDict {
	return starlark.StringDict{
		"vars":    goToStarlark(r.Vars),
		"actions": goToStarlark(r.Actions),
		"state":   goToStarlark(r.State),
		"env":     goToStarlark(r.Env),
	}
}

func (r *Resolver) evalExpr(expr string) (starlark.Value, error) {
	thread := &starlark.Thread{Name: "expr"}
	return starlark.Eval(thread, "<expr>", expr, r.makeGlobals())
}

// ResolveString замінює всі ${...} у рядку Starlark виразами.
func (r *Resolver) ResolveString(s string) (string, error) {
	var resolveErr error
	result := exprRe.ReplaceAllStringFunc(s, func(match string) string {
		if resolveErr != nil {
			return ""
		}
		expr := match[2 : len(match)-1]
		val, err := r.evalExpr(expr)
		if err != nil {
			resolveErr = fmt.Errorf("expr %q: %w", expr, err)
			return ""
		}
		return fmt.Sprintf("%v", starlarkToGo(val))
	})
	return result, resolveErr
}

// ResolveValue резолвить значення — якщо це рядок з одним виразом, повертає типізоване значення.
func (r *Resolver) ResolveValue(v any) (any, error) {
	s, ok := v.(string)
	if !ok {
		return v, nil
	}

	// якщо весь рядок — один вираз, повертаємо типізований результат
	if exprRe.MatchString(s) && exprRe.FindString(s) == s {
		expr := s[2 : len(s)-1]
		val, err := r.evalExpr(expr)
		if err != nil {
			return nil, fmt.Errorf("expr %q: %w", expr, err)
		}
		return starlarkToGo(val), nil
	}

	return r.ResolveString(s)
}

// ResolveMap рекурсивно резолвить всі рядки в map.
func (r *Resolver) ResolveMap(m map[string]any) (map[string]any, error) {
	result := make(map[string]any, len(m))
	for k, v := range m {
		resolved, err := r.resolveAny(v)
		if err != nil {
			return nil, fmt.Errorf("key %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}

func (r *Resolver) resolveAny(v any) (any, error) {
	switch val := v.(type) {
	case string:
		return r.ResolveValue(val)
	case map[string]any:
		return r.ResolveMap(val)
	case []any:
		result := make([]any, len(val))
		for i, item := range val {
			resolved, err := r.resolveAny(item)
			if err != nil {
				return nil, fmt.Errorf("[%d]: %w", i, err)
			}
			result[i] = resolved
		}
		return result, nil
	default:
		return v, nil
	}
}

// ExecMutateOn виконує Starlark скрипт з доступом до state (мутабельний) та event.
func (r *Resolver) ExecMutateOn(script string, event map[string]any) error {
	stateDict := goToStarlark(r.State).(*starlark.Dict)

	globals := starlark.StringDict{
		"vars":    goToStarlark(r.Vars),
		"actions": goToStarlark(r.Actions),
		"state":   stateDict,
		"env":     goToStarlark(r.Env),
		"event":   goToStarlark(event),
	}

	thread := &starlark.Thread{Name: "mutate_on"}
	if _, err := starlark.ExecFile(thread, "<mutate_on>", script, globals); err != nil {
		return err
	}

	r.State = starlarkToGo(stateDict).(map[string]any)
	return nil
}

// ResolveVars резолвить ${env['KEY']} вирази у значеннях vars.
func ResolveVars(v map[string]string, env map[string]string) (map[string]string, error) {
	r := &Resolver{Env: env}
	result := make(map[string]string, len(v))
	for k, val := range v {
		resolved, err := r.ResolveString(val)
		if err != nil {
			return nil, fmt.Errorf("var %q: %w", k, err)
		}
		result[k] = resolved
	}
	return result, nil
}
