package verdeps

import "sync"

type syncedStringMap interface {
	get(key string) (string, bool)
	set(key, value string)
	setIfAbsent(key, value string)
	clear()
	delete(key string)
	count() int
	each(fn func(string, string)) syncedStringMap
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

func (m *syncedStringMapImpl) setIfAbsent(key, value string) {
	m.lock.RLock()
	_, exists := m.values[key]
	m.lock.RUnlock()

	if !exists {
		m.lock.Lock()
		_, exists = m.values[key]
		if !exists {
			m.values[key] = value
		}
		m.lock.Unlock()
	}
}

func (m *syncedStringMapImpl) clear() {
	m.lock.Lock()
	for k := range m.values {
		delete(m.values, k)
	}
	m.lock.Unlock()
}

func (m *syncedStringMapImpl) delete(key string) {
	m.lock.Lock()
	if _, exists := m.values[key]; exists {
		delete(m.values, key)
	}
	m.lock.Unlock()
}

func (m *syncedStringMapImpl) count() int {
	m.lock.RLock()
	count := len(m.values)
	m.lock.RUnlock()
	return count
}

func (m *syncedStringMapImpl) each(fn func(string, string)) syncedStringMap {
	m.lock.RLock()
	for key, val := range m.values {
		fn(key, val)
	}
	m.lock.RUnlock()

	return m
}
