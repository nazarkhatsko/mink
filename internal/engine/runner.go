package engine

import (
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/nazarkhatsko/mink/internal/config"
	"github.com/nazarkhatsko/mink/internal/report"
	"github.com/nazarkhatsko/mink/internal/vars"
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
				event := map[string]any{
					"action": action.ID,
					"result": nil,
					"error":  map[string]any{"message": err.Error()},
				}
				_ = execCtx.resolver.ExecMutateOn(action.MutateOn.Fail, event)
			}
			e.reporter.ActionFail(action.ID, err, duration)
			return fmt.Errorf("action %q: %w", action.ID, err)
		}

		actions[action.ID] = map[string]any(output)

		if action.MutateOn.Done != "" {
			event := map[string]any{
				"action": action.ID,
				"result": map[string]any(output),
				"error":  nil,
			}
			if err := execCtx.resolver.ExecMutateOn(action.MutateOn.Done, event); err != nil {
				e.reporter.ActionFail(action.ID, fmt.Errorf("mutate_on.done: %w", err), duration)
				return fmt.Errorf("action %q mutate_on.done: %w", action.ID, err)
			}
		}

		e.reporter.ActionDone(action.ID, output, duration)
	}

	e.reporter.FlowDone(flow.Name)
	return nil
}
