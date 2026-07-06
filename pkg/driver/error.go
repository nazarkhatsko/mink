package driver

import (
	"context"
	"errors"
	"fmt"
)

// ErrorCode classifies why a driver's Execute call failed, so callers
// (engine, mutate_on scripts, reporters) can branch on it deterministically
// instead of parsing the error message.
type ErrorCode string

const (
	// ErrConfig means the flow author passed a missing or invalid option.
	ErrConfig ErrorCode = "config"
	// ErrTransport means an external call or process failed (network, exec, RPC).
	ErrTransport ErrorCode = "transport"
	// ErrTimeout means the action's context deadline was exceeded.
	ErrTimeout ErrorCode = "timeout"
	// ErrInternal covers everything else (marshaling, parsing, unexpected state).
	ErrInternal ErrorCode = "internal"
)

// Error is a classified driver error. Drivers should return this (via
// NewError/Wrap) instead of a bare fmt.Errorf so the engine can attach a
// machine-readable code to the mutate_on event.
type Error struct {
	Code    ErrorCode
	Message string
	Cause   error
}

func (e *Error) Error() string { return e.Message }
func (e *Error) Unwrap() error { return e.Cause }

// NewError builds a classified error with no wrapped cause.
func NewError(code ErrorCode, format string, a ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, a...)}
}

// Wrap builds a classified error that keeps cause reachable via errors.Is/As.
func Wrap(code ErrorCode, cause error, format string, a ...any) *Error {
	return &Error{Code: code, Message: fmt.Sprintf(format, a...), Cause: cause}
}

// Classify turns any error returned from a driver (or the engine) into a
// machine-readable {code, message} pair for mutate_on events and reporters.
// A timeout always wins over whatever code a driver attached, since it
// reflects what actually happened to the call.
func Classify(err error) map[string]any {
	code := ErrInternal

	var de *Error
	if errors.As(err, &de) {
		code = de.Code
	}
	if errors.Is(err, context.DeadlineExceeded) {
		code = ErrTimeout
	}

	return map[string]any{
		"code":    string(code),
		"message": err.Error(),
	}
}
