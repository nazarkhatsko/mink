package docutil

import (
	"fmt"
	"strings"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/internal/manual"
	"github.com/nazarkhatsko/mink/internal/render"
)

// Show reads an embedded doc page, renders it, and prints it to stdout.
func Show(parts ...string) error {
	content, err := manual.Page(parts...)
	if err != nil {
		return fmt.Errorf("documentation not found: %s", strings.Join(parts, "/"))
	}
	rendered, err := render.Markdown(content)
	if err != nil {
		return err
	}
	console.Print("%s", rendered)
	return nil
}
