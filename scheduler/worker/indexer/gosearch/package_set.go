package gosearch

import (
	"strings"
	"sync"
)

// packageSetEntry is an author-repo tuple in the packageSet.
type packageSetEntry struct {
	author string
	repo   string
}

// packageSet is a set of all author-repo combinations that map to packages.
type packageSet struct {
	lock     sync.RWMutex
	packages map[string]bool
}

// newPackageSet creates a new packageSet.
func newPackageSet() *packageSet {
	return &packageSet{packages: make(map[string]bool)}
}

// len returns the number of packages in the set.
func (ps *packageSet) len() int {
	ps.lock.RLock()
	size := len(ps.packages)
	ps.lock.RUnlock()

	return size
}

// add puts the specified author-repo combination into the set.
func (ps *packageSet) add(author, repo string) {
	ps.lock.Lock()
	ps.packages[author+"/"+repo] = true
	ps.lock.Unlock()
}

// remove deletes the specified author-repo combination from the set.
func (ps *packageSet) remove(author, repo string) {
	ps.lock.Lock()
	delete(ps.packages, author+"/"+repo)
	ps.lock.Unlock()
}

// stream pipes all of the entries in the set into the provided channel then
// closes it.
func (ps *packageSet) stream(packageSetEntries chan packageSetEntry) {
	ps.lock.RLock()
	for packageStr := range ps.packages {
		if i := strings.Index(packageStr, "/"); i != -1 && i+1 < len(packageStr) {
			packageSetEntries <- packageSetEntry{
				repo:   packageStr[i+1:],
				author: packageStr[:i],
			}
		}
	}
	ps.lock.RUnlock()

	close(packageSetEntries)
}

// contains returns true if the set contains the specified author-repo
// combination.
func (ps *packageSet) contains(author, repo string) bool {
	ps.lock.RLock()
	_, exists := ps.packages[author+"/"+repo]
	ps.lock.RUnlock()

	return exists
}
