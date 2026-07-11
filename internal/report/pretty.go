package report

import (
	"time"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Pretty struct{}

func NewPretty() *Pretty { return &Pretty{} }

func (r *Pretty) FlowStart(id, description string) {
	if description != "" {
		console.Printlnf("\n▶ flow: %s — %s", id, description)
	} else {
		console.Printlnf("\n▶ flow: %s", id)
	}
}

func (r *Pretty) FlowDone(id string, _ map[string]any) {
	console.Printlnf("✓ flow done: %s", id)
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
