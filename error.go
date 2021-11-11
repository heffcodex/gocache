package gocache

import (
	"fmt"

	"github.com/pkg/errors"
)

type NilError struct {
	err error
}

func (e *NilError) Error() string {
	return fmt.Sprintf("key does not exist (%s)", e.err)
}

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
