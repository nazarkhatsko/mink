package skills

import (
	"embed"
	"io/fs"
	"path/filepath"
	"strings"
)

//go:embed claude codex
var FS embed.FS

// Get returns the content of a skill's SKILL.md by platform and name.
func Get(platform, name string) ([]byte, error) {
	return FS.ReadFile(filepath.Join(platform, name, "SKILL.md"))
}

// GetVersion returns the version string from a skill's VERSION file.
func GetVersion(platform, name string) (string, error) {
	data, err := FS.ReadFile(filepath.Join(platform, name, "VERSION"))
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(data)), nil
}

// List returns names of all available skills for the given platform.
func List(platform string) ([]string, error) {
	entries, err := fs.ReadDir(FS, platform)
	if err != nil {
		return nil, err
	}
	var names []string
	for _, e := range entries {
		if e.IsDir() {
			names = append(names, e.Name())
		}
	}
	return names, nil
}
