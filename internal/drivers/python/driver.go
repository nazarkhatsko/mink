package python

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"

	"github.com/nazarkhatsko/mink/pkg/driver"
)

type Driver struct{}

func New() *Driver { return &Driver{} }

func (d *Driver) Name() string { return "python" }

func (d *Driver) Execute(ctx context.Context, options map[string]any) (driver.Output, error) {
	interpreter, _ := options["interpreter"].(string)
	if interpreter == "" {
		interpreter = "python3"
	}

	code, hasCode := options["code"].(string)
	script, hasScript := options["script"].(string)
	if hasCode && code == "" {
		hasCode = false
	}
	if hasScript && script == "" {
		hasScript = false
	}
	if hasCode == hasScript {
		return nil, driver.NewError(driver.ErrConfig, "python: exactly one of code or script is required")
	}

	scriptPath := script
	if hasCode {
		tmp, err := os.CreateTemp("", "mink-py-*.py")
		if err != nil {
			return nil, driver.Wrap(driver.ErrInternal, err, "python: create temp script: %v", err)
		}
		defer os.Remove(tmp.Name())

		if _, err := tmp.WriteString(code); err != nil {
			tmp.Close()
			return nil, driver.Wrap(driver.ErrInternal, err, "python: write temp script: %v", err)
		}
		if err := tmp.Close(); err != nil {
			return nil, driver.Wrap(driver.ErrInternal, err, "python: close temp script: %v", err)
		}
		scriptPath = tmp.Name()
	}

	var args []string
	if raw, ok := options["args"]; ok {
		list, ok := raw.([]any)
		if !ok {
			return nil, driver.NewError(driver.ErrConfig, "python: args must be a list")
		}
		for _, v := range list {
			args = append(args, fmt.Sprintf("%v", v))
		}
	}

	cmdArgs := append([]string{scriptPath}, args...)
	c := exec.CommandContext(ctx, interpreter, cmdArgs...)

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
			return nil, driver.Wrap(driver.ErrTransport, err, "python: %v", err)
		}
	}

	var stdout any
	if jsonErr := json.Unmarshal(stdoutBuf.Bytes(), &stdout); jsonErr != nil {
		stdout = stdoutBuf.String()
	}

	return driver.Output{
		"exit_code": exitCode,
		"stdout":    stdout,
		"stderr":    stderrBuf.String(),
		"success":   exitCode == 0,
	}, nil
}
