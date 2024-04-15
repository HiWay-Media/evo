//go:build !go1.13
// +build !go1.13

package internal

import "errors"

// Errorf wraps xerrors.Errorf.
func Errorf(format string, a ...interface{}) error { return errors.Errorf(format, a...) }

// ErrorAs wraps xerrors.As.
func ErrorAs(err error, target interface{}) bool { return errors.As(err, target) }

// ErrorIs wraps xerrors.Is.
func ErrorIs(err, target error) bool { return errors.Is(err, target) }

// NewError wraps xerrors.New.
func NewError(text string) error { return errors.New(text) }
