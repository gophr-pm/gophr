package verdeps

// revisionList is a wrapper for a revision slice that keeps track of revision
// types as well.
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
