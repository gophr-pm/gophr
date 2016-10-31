package verdeps

import "sync"

type syncedInt struct {
	val  int
	lock *sync.RWMutex
}

func newSyncedInt() *syncedInt {
	return &syncedInt{
		val:  0,
		lock: &sync.RWMutex{},
	}
}

func (si *syncedInt) value() int {
	si.lock.RLock()
	val := si.val
	si.lock.RUnlock()

	return val
}

func (si *syncedInt) increment() int {
	si.lock.Lock()
	newVal := si.val + 1
	si.val = newVal
	si.lock.Unlock()

	return newVal
}

func (si *syncedInt) decrement() int {
	si.lock.Lock()
	newVal := si.val - 1
	si.val = newVal
	si.lock.Unlock()

	return newVal
}
