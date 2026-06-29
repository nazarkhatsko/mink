package report

import (
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Reporter interface {
	FlowStart(name string)
	FlowDone(name string)
	ActionStart(id, description string)
	ActionDone(id string, output driver.Output, duration time.Duration)
	ActionFail(id string, err error, duration time.Duration)
}
