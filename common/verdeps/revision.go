package verdeps

import "sync"

// revision represents an a replacement to be made in a source file on disk.
type revision struct {
	path           string
	toIndex        int
	gophrURL       []byte
	fromIndex      int
	revisesImport  bool
	revisesPackage bool
}

func newImportRevision(spec *importSpec, newImportPath []byte) *revision {
	return &revision{
		path:          spec.filePath,
		toIndex:       int(spec.imports.Path.End()),
		gophrURL:      newImportPath,
		fromIndex:     int(spec.imports.Path.Pos()),
		revisesImport: true,
	}
}

func newPackageRevision(spec *packageSpec) *revision {
	return &revision{
		path:           spec.filePath,
		fromIndex:      spec.startIndex,
		revisesPackage: true,
	}
}

type revisionList struct {
	revs            []*revision
	importRevCount  int
	packageRevCount int
}

func newRevisionList() *revisionList {
	return &revisionList{}
}

func (r *revisionList) add(rev *revision) {
	if rev.revisesImport {
		r.importRevCount = r.importRevCount + 1
	} else if rev.revisesPackage {
		r.packageRevCount = r.packageRevCount + 1
	}

	r.revs = append(r.revs, rev)
}

func (r *revisionList) getRevs() []*revision {
	return r.revs
}

type syncedRevisionListMap struct {
	revLists map[string]*revisionList
	lock     *sync.RWMutex
}

func newSyncedRevisionListMap() *syncedRevisionListMap {
	return &syncedRevisionListMap{
		revLists: make(map[string]*revisionList),
		lock:     &sync.RWMutex{},
	}
}

func (m *syncedRevisionListMap) add(key string, rev *revision) {
	/// Get the rev list if it exists.
	m.lock.RLock()
	revList, exists := m.revLists[key]
	m.lock.RUnlock()

	if !exists {
		// It doesn't exist - lock again to create it.
		m.lock.Lock()
		revList, exists = m.revLists[key]
		// Make sure the status did not change while we were locking.
		if !exists {
			// The status did not change. So, create the rev list.
			revList = newRevisionList()
			m.revLists[key] = revList
		}
		m.lock.Unlock()
	}

	revList.add(rev)
}

func (m *syncedRevisionListMap) ready(key string, expectedImports int) bool {
	// Get the rev list if it exists.
	m.lock.RLock()
	revList, exists := m.revLists[key]
	m.lock.RUnlock()

	if !exists {
		return false
	}

	return revList.packageRevCount > 0 &&
		revList.importRevCount >= expectedImports
}

func (m *syncedRevisionListMap) getRevs(key string) []*revision {
	// Get the rev list if it exists.
	m.lock.RLock()
	revList, exists := m.revLists[key]
	m.lock.RUnlock()

	if !exists {
		return nil
	}

	return revList.revs
}

func (m *syncedRevisionListMap) clear() {
	m.lock.Lock()
	for k := range m.revLists {
		delete(m.revLists, k)
	}
	m.lock.Unlock()
}

func (m *syncedRevisionListMap) delete(key string) {
	m.lock.Lock()
	if _, exists := m.revLists[key]; exists {
		delete(m.revLists, key)
	}
	m.lock.Unlock()
}

func (m *syncedRevisionListMap) count() int {
	m.lock.RLock()
	count := len(m.revLists)
	m.lock.RUnlock()
	return count
}

func (m *syncedRevisionListMap) each(fn func(string, *revisionList)) {
	m.lock.RLock()
	for key, val := range m.revLists {
		fn(key, val)
	}
	m.lock.RUnlock()
}
