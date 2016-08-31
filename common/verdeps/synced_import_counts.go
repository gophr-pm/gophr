package verdeps

import "sync"

type syncedImportCounts struct {
	counts map[string]int
	lock   *sync.RWMutex
}

func newSyncedImportCounts() *syncedImportCounts {
	return &syncedImportCounts{
		counts: make(map[string]int),
		lock:   &sync.RWMutex{},
	}
}

func (sic *syncedImportCounts) importCountOf(path string) int {
	sic.lock.RLock()
	count := sic.counts[path]
	sic.lock.RUnlock()

	return count
}

func (sic *syncedImportCounts) setImportCount(path string, count int) {
	sic.lock.Lock()
	sic.counts[path] = count
	sic.lock.Unlock()
}
