package report

import (
	"fmt"
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Compact struct{}

func NewCompact() *Compact { return &Compact{} }

func (r *Compact) FlowStart(name string) {
	fmt.Printf("▶ %s\n", name)
}

func (r *Compact) FlowDone(name string) {
	fmt.Printf("✓ %s\n", name)
}

func (r *Compact) ActionStart(_, _ string) {}

func (r *Compact) ActionDone(id string, _ driver.Output, duration time.Duration) {
	fmt.Printf("  ✓ %s (%dms)\n", id, duration.Milliseconds())
}

func (r *Compact) ActionFail(id string, err error, duration time.Duration) {
	fmt.Printf("  ✗ %s (%dms): %s\n", id, duration.Milliseconds(), err)
}
