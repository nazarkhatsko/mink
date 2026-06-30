package engine

import (
	"context"
	"fmt"
	"os"

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

	resolvedVars, err := vars.ResolveVars(e.cfg.Vars, os.Getenv)
	if err != nil {
		return fmt.Errorf("vars: %w", err)
	}

	actions := make(map[string]map[string]any)
	execCtx := newContext(resolvedVars, actions)

	for _, action := range flow.Actions {
		instance, ok := e.cfg.Instances[action.Use]
		if !ok {
			return fmt.Errorf("instance %q not found", action.Use)
		}

		e.reporter.ActionStart(action.ID, action.Description)

		output, duration, err := e.executeAction(ctx, execCtx, action, instance)
		if err != nil {
			e.reporter.ActionFail(action.ID, err, duration)
			return fmt.Errorf("action %q: %w", action.ID, err)
		}

		actions[action.ID] = map[string]any(output)
		e.reporter.ActionDone(action.ID, output, duration)
	}

	e.reporter.FlowDone(flow.Name)
	return nil
}
