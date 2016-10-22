package verdeps

// fetchSHAResult is the result of a call to fetchSHA. When successful, it is a
// mapping between an import path and the sha it is paired with. When
// unsuccessful, it carries its error.
type fetchSHAResult struct {
	sha        string
	err        error
	successful bool
	importPath string
}

// newFetchSHASuccess creates a new fetchSHAResult, but specifies that fetchSHA
// completed successfully.
func newFetchSHASuccess(importPath, sha string) *fetchSHAResult {
	return &fetchSHAResult{
		sha:        sha,
		importPath: importPath,
		successful: true,
	}
}

// newFetchSHAFailure creates a new fetchSHAResult, but specifies that fetchSHA
// completed unsuccessfully.
func newFetchSHAFailure(err error) *fetchSHAResult {
	return &fetchSHAResult{
		err:        err,
		successful: false,
	}
}
