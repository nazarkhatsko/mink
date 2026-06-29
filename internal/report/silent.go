package report

import (
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Silent struct{}

func NewSilent() *Silent { return &Silent{} }

func (r *Silent) FlowStart(_ string)                          {}
func (r *Silent) FlowDone(_ string)                           {}
func (r *Silent) ActionStart(_, _ string)                     {}
func (r *Silent) ActionDone(_ string, _ driver.Output, _ time.Duration) {}
func (r *Silent) ActionFail(_ string, _ error, _ time.Duration) {}
