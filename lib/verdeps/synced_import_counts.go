package verdeps

import "sync"

type syncedImportCounts struct {
	counts map[string]int
	total  int
	lock   *sync.RWMutex
}

func newSyncedImportCounts() *syncedImportCounts {
	return &syncedImportCounts{
		counts: make(map[string]int),
		lock:   &sync.RWMutex{},
	}
}

func (sic *syncedImportCounts) totalCount() int {
	sic.lock.RLock()
	total := sic.total
	sic.lock.RUnlock()

	return total
}

func (sic *syncedImportCounts) importCountOf(path string) int {
	sic.lock.RLock()
	count := sic.counts[path]
	sic.lock.RUnlock()

	return count
}

func (sic *syncedImportCounts) setImportCount(path string, count int) {
	sic.lock.Lock()
	existingCount, _ := sic.counts[path]
	sic.counts[path] = count
	sic.total = sic.total + (count - existingCount)
	sic.lock.Unlock()
}
