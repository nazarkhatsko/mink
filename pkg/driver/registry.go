package driver

import "fmt"

var registry = map[string]Driver{}

func Register(d Driver) {
	registry[d.Name()] = d
}

func Get(name string) (Driver, error) {
	d, ok := registry[name]
	if !ok {
		return nil, fmt.Errorf("driver %q not found", name)
	}
	return d, nil
}

func List() []string {
	names := make([]string, 0, len(registry))
	for name := range registry {
		names = append(names, name)
	}
	return names
}
