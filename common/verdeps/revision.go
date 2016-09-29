package verdeps

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
