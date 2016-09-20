package main

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestVersionAndArchivePackage(t *testing.T) {
	args := packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			// TODO(skeswa): @shikkic, write something that tests if "args" is correct
			// here. :D kthxbai
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{}, errors.New("this is an error")
		},
		constructionZonePath: "/my/cons/path",
	}
	err := versionAndArchivePackage(args)
	assert.NotNil(t, err)
}
