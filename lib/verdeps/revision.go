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
