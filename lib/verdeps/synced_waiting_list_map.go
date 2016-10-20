package verdeps

import "sync"

type syncedWaitingListMap interface {
	get(key string) (specWaitingList, bool)
	setIfAbsent(key string, value specWaitingList)
	clear()
	each(fn func(string, specWaitingList)) syncedWaitingListMap
}

type syncedWaitingListMapImpl struct {
	values map[string]specWaitingList
	lock   *sync.RWMutex
}

func newSyncedWaitingListMap() syncedWaitingListMap {
	return &syncedWaitingListMapImpl{
		values: make(map[string]specWaitingList),
		lock:   &sync.RWMutex{},
	}
}

func (m *syncedWaitingListMapImpl) get(key string) (specWaitingList, bool) {
	m.lock.RLock()
	value, exists := m.values[key]
	m.lock.RUnlock()

	return value, exists
}

func (m *syncedWaitingListMapImpl) setIfAbsent(
	key string,
	value specWaitingList,
) {
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

func (m *syncedWaitingListMapImpl) clear() {
	m.lock.Lock()
	for k := range m.values {
		delete(m.values, k)
	}
	m.lock.Unlock()
}

func (m *syncedWaitingListMapImpl) each(
	fn func(string, specWaitingList),
) syncedWaitingListMap {
	m.lock.RLock()
	for key, val := range m.values {
		fn(key, val)
	}
	m.lock.RUnlock()

	return m
}
