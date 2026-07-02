package vars

import (
	"fmt"

	"go.starlark.net/starlark"
)

func goToStarlark(v any) starlark.Value {
	switch val := v.(type) {
	case nil:
		return starlark.None
	case bool:
		return starlark.Bool(val)
	case int:
		return starlark.MakeInt(val)
	case int64:
		return starlark.MakeInt64(val)
	case float64:
		return starlark.Float(val)
	case string:
		return starlark.String(val)
	case map[string]string:
		d := new(starlark.Dict)
		for k, v := range val {
			d.SetKey(starlark.String(k), starlark.String(v))
		}
		return d
	case map[string]any:
		d := new(starlark.Dict)
		for k, v := range val {
			d.SetKey(starlark.String(k), goToStarlark(v))
		}
		return d
	case map[string]map[string]any:
		d := new(starlark.Dict)
		for k, v := range val {
			d.SetKey(starlark.String(k), goToStarlark(v))
		}
		return d
	case []any:
		elems := make([]starlark.Value, len(val))
		for i, item := range val {
			elems[i] = goToStarlark(item)
		}
		return starlark.NewList(elems)
	default:
		return starlark.String(fmt.Sprintf("%v", val))
	}
}

func starlarkToGo(v starlark.Value) any {
	switch val := v.(type) {
	case starlark.NoneType:
		return nil
	case starlark.Bool:
		return bool(val)
	case starlark.Int:
		if i, ok := val.Int64(); ok {
			return i
		}
		return val.String()
	case starlark.Float:
		return float64(val)
	case starlark.String:
		return string(val)
	case *starlark.Dict:
		m := make(map[string]any, val.Len())
		for _, item := range val.Items() {
			if k, ok := item[0].(starlark.String); ok {
				m[string(k)] = starlarkToGo(item[1])
			}
		}
		return m
	case *starlark.List:
		result := make([]any, val.Len())
		for i := range val.Len() {
			result[i] = starlarkToGo(val.Index(i))
		}
		return result
	default:
		return val.String()
	}
}
