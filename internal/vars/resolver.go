package vars

import (
	"fmt"
	"regexp"
	"strings"
)

var exprRe = regexp.MustCompile(`\$\{([^}]+)\}`)

type Resolver struct {
	Vars    map[string]string
	Actions map[string]map[string]any
	Env     func(string) string
}

// ResolveString замінює всі ${...} у рядку.
func (r *Resolver) ResolveString(s string) (string, error) {
	var resolveErr error
	result := exprRe.ReplaceAllStringFunc(s, func(match string) string {
		if resolveErr != nil {
			return ""
		}
		inner := match[2 : len(match)-1]
		val, err := r.resolveExpr(inner)
		if err != nil {
			resolveErr = err
			return ""
		}
		return fmt.Sprintf("%v", val)
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
		inner := s[2 : len(s)-1]
		return r.resolveExpr(inner)
	}

	return r.ResolveString(s)
}

func (r *Resolver) resolveExpr(expr string) (any, error) {
	expr = strings.TrimSpace(expr)

	if strings.HasPrefix(expr, "vars.") {
		key := strings.TrimPrefix(expr, "vars.")
		val, ok := r.Vars[key]
		if !ok {
			return nil, fmt.Errorf("var %q not found", key)
		}
		return val, nil
	}

	if strings.HasPrefix(expr, "env.") {
		key := strings.TrimPrefix(expr, "env.")
		return r.Env(key), nil
	}

	if strings.HasPrefix(expr, "actions.") {
		rest := strings.TrimPrefix(expr, "actions.")
		parts := strings.SplitN(rest, ".", 2)
		actionID := parts[0]

		output, ok := r.Actions[actionID]
		if !ok {
			return nil, fmt.Errorf("action %q not found", actionID)
		}

		if len(parts) == 1 {
			return output, nil
		}
		return dotGet(any(output), parts[1])
	}

	return nil, fmt.Errorf("unknown expression: %q", expr)
}

// ResolveVars резолвить ${env.*} вирази у значеннях vars.
func ResolveVars(v map[string]string, env func(string) string) (map[string]string, error) {
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
