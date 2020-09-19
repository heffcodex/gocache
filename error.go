package gocache

import "github.com/pkg/errors"

type RememberError struct {
	fullKey string
	err     error
}

func (e *RememberError) Error() string {
	return e.err.Error()
}

func (e *RememberError) FullKey() string {
	return e.fullKey
}

func WrapRememberError(err error, fullKey string) *RememberError {
	return &RememberError{err: errors.Wrap(err, "cannot remember cache value"), fullKey: fullKey}
}
