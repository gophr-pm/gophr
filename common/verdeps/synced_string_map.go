package verdeps

import "sync"

type syncedStringMap struct {
	values map[string]string
	lock   *sync.RWMutex
}

func newSyncedStringMap() *syncedStringMap {
	return &syncedStringMap{
		values: make(map[string]string),
		lock:   &sync.RWMutex{},
	}
}

func (m *syncedStringMap) get(key string) (string, bool) {
	m.lock.RLock()
	value, exists := m.values[key]
	m.lock.RUnlock()

	return value, exists
}

func (m *syncedStringMap) set(key, value string) {
	m.lock.Lock()
	m.values[key] = value
	m.lock.Unlock()
}

func (m *syncedStringMap) setIfAbsent(key, value string) {
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

func (m *syncedStringMap) clear() {
	m.lock.Lock()
	for k := range m.values {
		delete(m.values, k)
	}
	m.lock.Unlock()
}

func (m *syncedStringMap) delete(key string) {
	m.lock.Lock()
	if _, exists := m.values[key]; exists {
		delete(m.values, key)
	}
	m.lock.Unlock()
}

func (m *syncedStringMap) count() int {
	m.lock.RLock()
	count := len(m.values)
	m.lock.RUnlock()
	return count
}

func (m *syncedStringMap) each(fn func(string, string)) *syncedStringMap {
	m.lock.RLock()
	for key, val := range m.values {
		fn(key, val)
	}
	m.lock.RUnlock()

	return m
}
