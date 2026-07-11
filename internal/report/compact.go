package report

import (
	"time"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Compact struct{}

func NewCompact() *Compact { return &Compact{} }

func (r *Compact) FlowStart(id, _ string) {
	console.Printlnf("▶ %s", id)
}

func (r *Compact) FlowDone(id string, _ map[string]any) {
	console.Printlnf("✓ %s", id)
}

func (r *Compact) ActionStart(_, _ string) {}

func (r *Compact) ActionDone(id string, _ driver.Output, duration time.Duration, _ map[string]any) {
	console.Printlnf("  ✓ %s (%dms)", id, duration.Milliseconds())
}

func (r *Compact) ActionFail(id string, _ driver.Output, err error, duration time.Duration, _ map[string]any) {
	console.Printlnf("  ✗ %s (%dms): %s", id, duration.Milliseconds(), err)
}
