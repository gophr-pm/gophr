package main

import (
	"errors"
	"testing"

	"github.com/gophr-pm/gophr/common/verdeps"
	"github.com/stretchr/testify/assert"
)

func TestVersionAndArchivePackage(t *testing.T) {
	args := packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{}, errors.New("this is an error")
		},
		constructionZonePath: "/my/cons/path",
	}
	err := versionAndArchivePackage(args)
	assert.NotNil(t, err, "this should return an error")

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return errors.New("this is an error")
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			return
		},
	}
	err = versionAndArchivePackage(args)
	assert.NotNil(t, err, "this should return an error")

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return false, errors.New("this should return an error")
		},
	}
	err = versionAndArchivePackage(args)
	assert.NotNil(t, err)

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return false, nil
		},
		isPackageArchived: func(args packageArchivalCheckerArgs) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return false, errors.New("this should return an error")
		},
	}
	err = versionAndArchivePackage(args)
	assert.NotNil(t, err)

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return false, nil
		},
		isPackageArchived: func(args packageArchivalCheckerArgs) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return true, nil
		},
	}
	err = versionAndArchivePackage(args)
	assert.Nil(t, err)

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return false, nil
		},
		isPackageArchived: func(args packageArchivalCheckerArgs) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return false, nil
		},
	}
	err = versionAndArchivePackage(args)
	assert.NotNil(t, err)

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return true, nil
		},
		pushToDepot: func(args packagePusherArgs) error {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return errors.New("this is an error")
		},
		destroyDepotRepo: func(author, repo, sha string) error {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return errors.New("this is an error")
		},
	}
	err = versionAndArchivePackage(args)
	assert.NotNil(t, err)

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
				workDirPath:    "/work/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			assert.Equal(t, "/work/dir/path", workDirPath)
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return true, nil
		},
		pushToDepot: func(args packagePusherArgs) error {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return errors.New("this is an error")
		},
		destroyDepotRepo: func(author, repo, sha string) error {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return nil
		},
	}
	err = versionAndArchivePackage(args)
	assert.NotNil(t, err)

	args = packageVersionerArgs{
		sha:    "mysha",
		repo:   "myrepo",
		author: "myauthor",
		downloadPackage: func(args packageDownloaderArgs) (packageDownloadPaths, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			assert.Equal(t, "/my/cons/path", args.constructionZonePath)
			return packageDownloadPaths{
				archiveDirPath: "/archive/dir/path",
				workDirPath:    "/work/dir/path",
			}, nil
		},
		constructionZonePath: "/my/cons/path",
		versionDeps: func(args verdeps.VersionDepsArgs) error {
			assert.Equal(t, "myauthor", args.Author)
			assert.Equal(t, "myrepo", args.Repo)
			assert.Equal(t, "mysha", args.SHA)
			assert.Equal(t, "/archive/dir/path", args.Path)
			return nil
		},
		attemptWorkDirDeletion: func(workDirPath string) {
			assert.Equal(t, "/work/dir/path", workDirPath)
			return
		},
		createDepotRepo: func(author, repo, sha string) (bool, error) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return true, nil
		},
		pushToDepot: func(args packagePusherArgs) error {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return nil
		},
		recordPackageArchival: func(args packageArchivalRecorderArgs) {
			assert.Equal(t, "myauthor", args.author)
			assert.Equal(t, "myrepo", args.repo)
			assert.Equal(t, "mysha", args.sha)
			return
		},
	}
	err = versionAndArchivePackage(args)
	assert.Nil(t, err)
}
