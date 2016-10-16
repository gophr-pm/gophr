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

func (si *syncedInt) increment() {
	si.lock.Lock()
	si.val = si.val + 1
	si.lock.Unlock()
}

func (si *syncedInt) decrement() {
	si.lock.Lock()
	si.val = si.val - 1
	si.lock.Unlock()
}
