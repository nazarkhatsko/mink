package render

import "github.com/charmbracelet/glamour"

// Markdown renders markdown source for terminal display.
func Markdown(src []byte) (string, error) {
	r, err := glamour.NewTermRenderer(
		glamour.WithAutoStyle(),
		glamour.WithWordWrap(100),
	)
	if err != nil {
		return "", err
	}
	return r.Render(string(src))
}
