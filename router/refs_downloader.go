package main

import "github.com/skeswa/gophr/common"

// refsDownloader is responsible for downloading the git refs for a package.
type refsDownloader interface {
	// downloadRefs downloads the refs for a package matching the specified author
	// and repo.
	downloadRefs(author, repo string) (common.Refs, error)
}
