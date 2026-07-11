package report

import (
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Reporter interface {
	FlowStart(id, description string)
	FlowDone(id string, state map[string]any)
	ActionStart(id, description string)
	ActionDone(id string, output driver.Output, duration time.Duration, state map[string]any)
	ActionFail(id string, output driver.Output, err error, duration time.Duration, state map[string]any)
}
