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

func (d *Driver) Methods() []string { return []string{"run_code", "run_script"} }

func (d *Driver) Options(method string) []driver.Option {
	opts := []driver.Option{
		{Name: "args"},
		{Name: "env"},
		{Name: "dir"},
	}
	if method == "run_script" {
		return append([]driver.Option{{Name: "script", Required: true}}, opts...)
	}
	return append([]driver.Option{{Name: "code", Required: true}}, opts...)
}

func (d *Driver) Execute(ctx context.Context, method string, options map[string]any) (driver.Output, error) {
	switch method {
	case "run_code":
		code, _ := options["code"].(string)
		if code == "" {
			return nil, driver.NewError(driver.ErrConfig, "python: code is required")
		}
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
		return d.run(ctx, tmp.Name(), options)
	case "run_script":
		script, _ := options["script"].(string)
		if script == "" {
			return nil, driver.NewError(driver.ErrConfig, "python: script is required")
		}
		return d.run(ctx, script, options)
	default:
		return nil, driver.NewError(driver.ErrConfig, "python: unknown method %q", method)
	}
}

func (d *Driver) run(ctx context.Context, scriptPath string, options map[string]any) (driver.Output, error) {
	interpreter, _ := options["interpreter"].(string)
	if interpreter == "" {
		interpreter = "python3"
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
