package verdeps

// revision represents an import in the filesystem that needs to be versioned with a gophr URL.
type revision struct {
	path               string
	gophrURL           []byte
	toIndex, fromIndex int
}

func newRevision(spec *importSpec, newImportPath []byte) *revision {
	return &revision{
		path:      spec.filePath,
		toIndex:   int(spec.imports.Path.End()),
		gophrURL:  newImportPath,
		fromIndex: int(spec.imports.Path.Pos()),
	}
}
