package doc

import (
	"embed"
	"path/filepath"
)

//go:embed *.md drivers examples
var FS embed.FS

// Page reads an embedded doc page, e.g. Page("configuration") or Page("drivers", "http").
func Page(parts ...string) ([]byte, error) {
	return FS.ReadFile(filepath.Join(parts...) + ".md")
}
