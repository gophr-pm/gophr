package verdeps

import "sync"

type syncedStringMap interface {
	get(key string) (string, bool)
	set(key, value string)
}

type syncedStringMapImpl struct {
	values map[string]string
	lock   *sync.RWMutex
}

func newSyncedStringMap() syncedStringMap {
	return &syncedStringMapImpl{
		values: make(map[string]string),
		lock:   &sync.RWMutex{},
	}
}

func (m *syncedStringMapImpl) get(key string) (string, bool) {
	m.lock.RLock()
	value, exists := m.values[key]
	m.lock.RUnlock()

	return value, exists
}

func (m *syncedStringMapImpl) set(key, value string) {
	m.lock.Lock()
	m.values[key] = value
	m.lock.Unlock()
}
