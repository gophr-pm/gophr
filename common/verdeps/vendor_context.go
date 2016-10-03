package verdeps

import (
	"fmt"
	"sync"
)

type vendorContext struct {
	finalized bool
	packages  map[string]bool
	parent    *vendorContext
	depth     int
	lock      *sync.RWMutex
}

func newVendorContext() *vendorContext {
	return &vendorContext{
		packages: make(map[string]bool),
		lock:     &sync.RWMutex{},
	}
}

func (p *vendorContext) add(pkg string) {
	if p.finalized {
		panic("add cannot be called on a finalized vendor context")
	}

	p.lock.Lock()
	p.packages[pkg] = true
	p.lock.Unlock()
}

func (p *vendorContext) String() string {
	p.lock.RLock()
	str := fmt.Sprintf(
		"vendorContext{ finalized: %v, numberOfPackages: %d, depth: %d }",
		p.finalized,
		len(p.packages),
		p.depth)
	p.lock.RUnlock()

	return str
}

func (p *vendorContext) contains(pkg string) bool {
	p.lock.RLock()
	_, exists := p.packages[pkg]
	p.lock.RUnlock()

	if !exists && p.parent != nil {
		exists = p.parent.contains(pkg)
	}

	return exists
}

func (p *vendorContext) finalize() {
	p.finalized = true
}

func (p *vendorContext) spawnChildContext() *vendorContext {
	p.lock.RLock()
	childContext := &vendorContext{
		packages: make(map[string]bool),
		parent:   p,
		depth:    p.depth + 1,
		lock:     &sync.RWMutex{},
	}
	p.lock.RUnlock()

	return childContext
}
