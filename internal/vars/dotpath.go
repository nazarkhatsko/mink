package vars

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

var arrIndex = regexp.MustCompile(`^(.+)\[(\d+)\]$`)

func dotGet(v any, path string) (any, error) {
	if path == "" {
		return v, nil
	}

	parts := strings.SplitN(path, ".", 2)
	key := parts[0]
	rest := ""
	if len(parts) == 2 {
		rest = parts[1]
	}

	if m := arrIndex.FindStringSubmatch(key); m != nil {
		key = m[1]
		idx, _ := strconv.Atoi(m[2])

		child, err := mapGet(v, key)
		if err != nil {
			return nil, err
		}

		arr, ok := child.([]any)
		if !ok {
			return nil, fmt.Errorf("expected array at %q", key)
		}
		if idx >= len(arr) {
			return nil, fmt.Errorf("index %d out of bounds at %q", idx, key)
		}
		return dotGet(arr[idx], rest)
	}

	child, err := mapGet(v, key)
	if err != nil {
		return nil, err
	}
	return dotGet(child, rest)
}

func mapGet(v any, key string) (any, error) {
	if key == "" {
		return v, nil
	}
	m, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("expected object, got %T", v)
	}
	val, ok := m[key]
	if !ok {
		return nil, fmt.Errorf("key %q not found", key)
	}
	return val, nil
}
