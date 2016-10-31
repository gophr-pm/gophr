package verdeps

import "sync"

type specWaitingList interface {
	add(spec *importSpec) bool
	clear() []*importSpec
}

type specWaitingListImpl struct {
	lock    *sync.Mutex
	specs   []*importSpec
	cleared bool
}

func newSpecWaitingList(initialSpecs ...*importSpec) specWaitingList {
	return &specWaitingListImpl{
		lock:    &sync.Mutex{},
		specs:   initialSpecs,
		cleared: false,
	}
}

// add adds a spec to the waiting list and returns true if it was successful.
func (swl *specWaitingListImpl) add(spec *importSpec) bool {
	swl.lock.Lock()
	defer swl.lock.Unlock()

	if swl.cleared {
		return false
	}

	swl.specs = append(swl.specs, spec)
	return true
}

// clear returns every spec on the waiting list and prevents all future adds from
// succeeding.
func (swl *specWaitingListImpl) clear() []*importSpec {
	swl.lock.Lock()
	defer swl.lock.Unlock()

	specs := swl.specs
	swl.specs = nil
	swl.cleared = true
	return specs
}
