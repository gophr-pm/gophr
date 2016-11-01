package errors

import "sync"

// SyncedErrors is a thread safe wrapper around an error
// slice.s
type SyncedErrors struct {
	errors []error
	lock   *sync.RWMutex
}

// NewSyncedErrors creates new SyncedErrors.
func NewSyncedErrors() *SyncedErrors {
	return &SyncedErrors{
		errors: nil,
		lock:   &sync.RWMutex{},
	}
}

// Get returns all the errors in the error slice.
func (si *SyncedErrors) Get() []error {
	si.lock.RLock()
	errors := si.errors
	si.lock.RUnlock()

	return errors
}

// Len returns the length of the error slice.
func (si *SyncedErrors) Len() int {
	si.lock.RLock()
	length := len(si.errors)
	si.lock.RUnlock()

	return length
}

// Add adds an error to the acculated errors.
func (si *SyncedErrors) Add(errors ...error) {
	si.lock.Lock()
	si.errors = append(si.errors, errors...)
	si.lock.Unlock()
}

// Compose joins all the accumulated individual errors into one combined
// error.
func (si *SyncedErrors) Compose(msg string) error {
	si.lock.RLock()
	err := ComposeErrors(msg, si.errors)
	si.lock.RUnlock()

	return err
}
