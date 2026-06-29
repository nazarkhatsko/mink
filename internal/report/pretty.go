package report

import (
	"fmt"
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Pretty struct{}

func NewPretty() *Pretty { return &Pretty{} }

func (r *Pretty) FlowStart(name string) {
	fmt.Printf("\n▶ flow: %s\n", name)
}

func (r *Pretty) FlowDone(name string) {
	fmt.Printf("✓ flow done: %s\n", name)
}

func (r *Pretty) ActionStart(id, description string) {
	if description != "" {
		fmt.Printf("  · %s — %s\n", id, description)
	} else {
		fmt.Printf("  · %s\n", id)
	}
}

func (r *Pretty) ActionDone(id string, _ driver.Output, duration time.Duration) {
	fmt.Printf("  ✓ %s (%dms)\n", id, duration.Milliseconds())
}

func (r *Pretty) ActionFail(id string, err error, duration time.Duration) {
	fmt.Printf("  ✗ %s (%dms): %s\n", id, duration.Milliseconds(), err)
}
