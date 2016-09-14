package main

import "github.com/skeswa/gophr/common"

// refsDownloader is responsible for downloading the git refs for a package.
type refsDownloader interface {
	// downloadRefs downloads the refs for a package matching the specified author
	// and repo.
	downloadRefs(author, repo string) (common.Refs, error)
}

// defaultRefsDownloader is a refs downloader that actually goes out to the
// network in order to download refs.
type defaultRefsDownloader struct{}

// downloadRefs downloads the refs for a package matching the specified author
// and repo.
func (d defaultRefsDownloader) downloadRefs(
	author string,
	repo string) (common.Refs, error) {
	// This type simply serves as a proxy for the actual refs package.
	return common.FetchRefs(author, repo)
}
