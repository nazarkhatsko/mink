package report

import (
	"time"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Compact struct{}

func NewCompact() *Compact { return &Compact{} }

func (r *Compact) FlowStart(name string) {
	console.Printlnf("▶ %s", name)
}

func (r *Compact) FlowDone(name string) {
	console.Printlnf("✓ %s", name)
}

func (r *Compact) ActionStart(_, _ string) {}

func (r *Compact) ActionDone(id string, _ driver.Output, duration time.Duration) {
	console.Printlnf("  ✓ %s (%dms)", id, duration.Milliseconds())
}

func (r *Compact) ActionFail(id string, err error, duration time.Duration) {
	console.Printlnf("  ✗ %s (%dms): %s", id, duration.Milliseconds(), err)
}
