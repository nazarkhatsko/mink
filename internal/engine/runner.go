package engine

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/internal/report"
	"github.com/nazarkhatsko/mink/internal/vars"
	"github.com/nazarkhatsko/mink/pkg/driver"
)

type engine struct {
	cfg      *config.Config
	reporter report.Reporter
}

func New(cfg *config.Config, reporter report.Reporter) *engine {
	return &engine{cfg: cfg, reporter: reporter}
}

func (e *engine) Run(ctx context.Context, flowName string) error {
	for _, flow := range e.cfg.Flows {
		if flowName != "" && flow.Name != flowName {
			continue
		}
		if err := e.runFlow(ctx, flow); err != nil {
			return err
		}
	}
	return nil
}

func (e *engine) runFlow(ctx context.Context, flow config.Flow) error {
	e.reporter.FlowStart(flow.Name)

	envVars := make(map[string]string)
	for _, entry := range os.Environ() {
		parts := strings.SplitN(entry, "=", 2)
		if len(parts) == 2 {
			envVars[parts[0]] = parts[1]
		}
	}

	resolvedVars, err := vars.ResolveVars(e.cfg.Vars, envVars)
	if err != nil {
		return fmt.Errorf("vars: %w", err)
	}

	actions := make(map[string]map[string]any)
	state := make(map[string]any)
	execCtx := newContext(resolvedVars, actions, state)

	for _, action := range flow.Actions {
		instance, ok := e.cfg.Instances[action.Use]
		if !ok {
			return fmt.Errorf("instance %q not found", action.Use)
		}

		e.reporter.ActionStart(action.ID, action.Description)

		output, duration, err := e.executeAction(ctx, execCtx, action, instance)
		if err != nil {
			if action.MutateOn.Fail != "" {
				_ = execCtx.resolver.ExecMutateOn(action.MutateOn.Fail, buildEvent(action.ID, output, err))
			}
			e.reporter.ActionFail(action.ID, output, err, duration, execCtx.resolver.State)
			return fmt.Errorf("action %q: %w", action.ID, err)
		}

		actions[action.ID] = map[string]any(output)

		if action.MutateOn.Done != "" {
			if err := execCtx.resolver.ExecMutateOn(action.MutateOn.Done, buildEvent(action.ID, output, nil)); err != nil {
				mutateErr := fmt.Errorf("mutate_on.done: %w", err)
				e.reporter.ActionFail(action.ID, output, mutateErr, duration, execCtx.resolver.State)
				return fmt.Errorf("action %q mutate_on.done: %w", action.ID, err)
			}
		}

		e.reporter.ActionDone(action.ID, output, duration, execCtx.resolver.State)
	}

	e.reporter.FlowDone(flow.Name, execCtx.resolver.State)
	return nil
}

// buildEvent builds the object passed to mutate_on.done/mutate_on.fail
// scripts as `event`. Its shape is identical on both paths — output is always
// the driver's (possibly partial) output, error is nil unless err != nil —
// so a fail script can inspect what the driver actually returned. Check
// `event['error'] == None` to tell done from fail.
func buildEvent(actionID string, output driver.Output, err error) map[string]any {
	event := map[string]any{
		"action_id": actionID,
		"output":    map[string]any(output),
		"error":     nil,
	}
	if err != nil {
		event["error"] = driver.Classify(err)
	}
	return event
}
