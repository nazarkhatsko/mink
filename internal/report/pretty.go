package report

import (
	"time"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Pretty struct{}

func NewPretty() *Pretty { return &Pretty{} }

func (r *Pretty) FlowStart(name string) {
	console.Printlnf("\n▶ flow: %s", name)
}

func (r *Pretty) FlowDone(name string, _ map[string]any) {
	console.Printlnf("✓ flow done: %s", name)
}

func (r *Pretty) ActionStart(id, description string) {
	if description != "" {
		console.Printlnf("  · %s — %s", id, description)
	} else {
		console.Printlnf("  · %s", id)
	}
}

func (r *Pretty) ActionDone(id string, _ driver.Output, duration time.Duration, _ map[string]any) {
	console.Printlnf("  ✓ %s (%dms)", id, duration.Milliseconds())
}

func (r *Pretty) ActionFail(id string, _ driver.Output, err error, duration time.Duration, _ map[string]any) {
	console.Printlnf("  ✗ %s (%dms): %s", id, duration.Milliseconds(), err)
}
