package verdeps

import "sync"

type syncedWaitingListMap struct {
	values map[string]specWaitingList
	lock   *sync.RWMutex
}

func newSyncedWaitingListMap() *syncedWaitingListMap {
	return &syncedWaitingListMap{
		values: make(map[string]specWaitingList),
		lock:   &sync.RWMutex{},
	}
}

func (m *syncedWaitingListMap) get(key string) (specWaitingList, bool) {
	m.lock.RLock()
	value, exists := m.values[key]
	m.lock.RUnlock()

	return value, exists
}

func (m *syncedWaitingListMap) set(key string, value specWaitingList) {
	m.lock.Lock()
	m.values[key] = value
	m.lock.Unlock()
}

func (m *syncedWaitingListMap) setIfAbsent(key string, value specWaitingList) {
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

func (m *syncedWaitingListMap) clear() {
	m.lock.Lock()
	for k := range m.values {
		delete(m.values, k)
	}
	m.lock.Unlock()
}

func (m *syncedWaitingListMap) delete(key string) {
	m.lock.Lock()
	if _, exists := m.values[key]; exists {
		delete(m.values, key)
	}
	m.lock.Unlock()
}

func (m *syncedWaitingListMap) count() int {
	m.lock.RLock()
	count := len(m.values)
	m.lock.RUnlock()
	return count
}

func (m *syncedWaitingListMap) each(fn func(string, specWaitingList)) *syncedWaitingListMap {
	m.lock.RLock()
	for key, val := range m.values {
		fn(key, val)
	}
	m.lock.RUnlock()

	return m
}
