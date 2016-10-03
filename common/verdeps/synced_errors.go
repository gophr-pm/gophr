package verdeps

import "sync"

type syncedErrors struct {
	errors []error
	lock   *sync.RWMutex
}

func newSyncedErrors() *syncedErrors {
	return &syncedErrors{
		errors: nil,
		lock:   &sync.RWMutex{},
	}
}

func (si *syncedErrors) get() []error {
	si.lock.RLock()
	errors := si.errors
	si.lock.RUnlock()

	return errors
}

func (si *syncedErrors) len() int {
	si.lock.RLock()
	length := len(si.errors)
	si.lock.RUnlock()

	return length
}

func (si *syncedErrors) add(errors ...error) {
	si.lock.Lock()
	si.errors = append(si.errors, errors...)
	si.lock.Unlock()
}
