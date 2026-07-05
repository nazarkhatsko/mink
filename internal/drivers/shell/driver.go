package shell

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "shell" }

func (d *Driver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
	command, ok := options["command"]
	if !ok {
		return nil, fmt.Errorf("shell: command is required")
	}
	cmd, ok := command.(string)
	if !ok {
		return nil, fmt.Errorf("shell: command must be a string, got %T", command)
	}

	c := exec.CommandContext(ctx, "sh", "-c", cmd)

	c.Env = os.Environ()
	if raw, ok := options["env"]; ok {
		switch v := raw.(type) {
		case map[string]any:
			for key, val := range v {
				c.Env = append(c.Env, fmt.Sprintf("%s=%v", key, val))
			}
		case map[string]string:
			for key, val := range v {
				c.Env = append(c.Env, fmt.Sprintf("%s=%s", key, val))
			}
		}
	}

	if dir, ok := options["dir"]; ok {
		if dirStr, ok := dir.(string); ok && dirStr != "" {
			c.Dir = dirStr
		}
	}

	var stdoutBuf, stderrBuf bytes.Buffer
	c.Stdout = &stdoutBuf
	c.Stderr = &stderrBuf

	err := c.Run()

	exitCode := 0
	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			exitCode = exitErr.ExitCode()
		} else {
			return nil, fmt.Errorf("shell: %w", err)
		}
	}

	return driver.Output{
		"exit_code": exitCode,
		"stdout":    stdoutBuf.String(),
		"stderr":    stderrBuf.String(),
		"success":   exitCode == 0,
	}, nil
}
