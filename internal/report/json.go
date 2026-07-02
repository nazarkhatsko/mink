package report

import (
	"encoding/json"
	"time"

	"github.com/nazarkhatsko/mink/internal/console"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type JSON struct {
	flow string
}

func NewJSON() *JSON { return &JSON{} }

func (r *JSON) FlowStart(name string) {
	r.flow = name
	r.write(map[string]any{"type": "flow_start", "flow": name})
}

func (r *JSON) FlowDone(name string, state map[string]any) {
	r.write(map[string]any{"type": "flow_done", "flow": name, "state": state})
}

func (r *JSON) ActionStart(_, _ string) {}

func (r *JSON) ActionDone(id string, output driver.Output, duration time.Duration, state map[string]any) {
	r.write(map[string]any{
		"type":        "action_done",
		"flow":        r.flow,
		"action":      id,
		"duration_ms": duration.Milliseconds(),
		"result":      output,
		"state":       state,
	})
}

func (r *JSON) ActionFail(id string, err error, duration time.Duration, state map[string]any) {
	r.write(map[string]any{
		"type":        "action_fail",
		"flow":        r.flow,
		"action":      id,
		"duration_ms": duration.Milliseconds(),
		"error":       err.Error(),
		"state":       state,
	})
}

func (r *JSON) write(v any) {
	b, _ := json.Marshal(v)
	console.Println(string(b))
}
