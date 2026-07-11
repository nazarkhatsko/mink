package report

import (
	"time"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Silent struct{}

func NewSilent() *Silent { return &Silent{} }

func (r *Silent) FlowStart(_, _ string)                                                            {}
func (r *Silent) FlowDone(_ string, _ map[string]any)                                              {}
func (r *Silent) ActionStart(_, _ string)                                                          {}
func (r *Silent) ActionDone(_ string, _ driver.Output, _ time.Duration, _ map[string]any)          {}
func (r *Silent) ActionFail(_ string, _ driver.Output, _ error, _ time.Duration, _ map[string]any) {}
