package errors

import "errors"

// New acts as a proxy to the stdlib's errors.New.
func New(text string) error {
	return errors.New(text)
}
