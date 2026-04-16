package errors

import (
	stderrors "errors"
	"fmt"
)

type FixableBy string

const (
	FixableByAgent FixableBy = "agent"
	FixableByHuman FixableBy = "human"
	FixableByRetry FixableBy = "retry"
)

type APIError struct {
	Message   string    `json:"error"`
	Hint      string    `json:"hint,omitempty"`
	FixableBy FixableBy `json:"fixable_by"`
	Cause     error     `json:"-"`
}

func (e *APIError) Error() string { return e.Message }
func (e *APIError) Unwrap() error { return e.Cause }

func New(message string, fixableBy FixableBy) *APIError {
	return &APIError{Message: message, FixableBy: fixableBy}
}

func Newf(fixableBy FixableBy, format string, args ...any) *APIError {
	return &APIError{Message: fmt.Sprintf(format, args...), FixableBy: fixableBy}
}

func Wrap(err error, fixableBy FixableBy) *APIError {
	return &APIError{Message: err.Error(), FixableBy: fixableBy, Cause: err}
}

func (e *APIError) WithHint(hint string) *APIError {
	e.Hint = hint
	return e
}

// As is a convenience wrapper around errors.As.
func As(err error, target **APIError) bool {
	return stderrors.As(err, target)
}
