package main

import (
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/git"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
	"github.com/gophr-pm/gophr/lib/verdeps"
)

// refsDownloader is responsible for downloading the git refs for a package.
type refsDownloader func(author, repo string) (lib.Refs, error)

// fullSHAFetcher is responsible for fetching a full commit SHA from a short SHA.
type fullSHAFetcher func(args github.RequestService) (string, error)

// packageDownloadRecorderArgs is the arguments struct for
// packageDownloadRecorders.
type packageDownloadRecorderArgs struct {
	db     db.BatchingQueryable
	sha    string
	repo   string
	ghSvc  github.RequestService
	author string
}

// packageDownloadRecorder is responsible for recording package downloads. If
// there is a problem while recording, then the error is logged instead of
// bubbled.
type packageDownloadRecorder func(args packageDownloadRecorderArgs)

// packageArchivalArgs is the arguments struct for packageArchivalRecorders and
// packageArchivalCheckers.
type packageArchivalRecorderArgs struct {
	db     db.Queryable
	sha    string
	repo   string
	author string
}

// packageArchivalRecorder is responsible for recording package archival. If
// there is a problem while recording, then the error is logged instead of
// bubbled.
type packageArchivalRecorder func(args packageArchivalRecorderArgs)

// packageArchivalArgs is the arguments struct for packageArchivalRecorders and
// packageArchivalCheckers.
type packageArchivalCheckerArgs struct {
	db                    db.Queryable
	sha                   string
	repo                  string
	author                string
	packageExistsInDepot  depotExistenceChecker
	recordPackageArchival packageArchivalRecorder
	isPackageArchivedInDB dbPackageArchivalChecker
}

// packageArchivalChecker is responsible for checking whether a package has
// been archived or not. Returns true if the package has been archived, and
// false otherwise.
type packageArchivalChecker func(args packageArchivalCheckerArgs) (bool, error)

// packageVersionerArgs is the arguments struct for packageVersioners.
type packageVersionerArgs struct {
	io                     io.IO
	db                     db.Queryable
	sha                    string
	repo                   string
	conf                   *config.Config
	creds                  *config.Credentials
	ghSvc                  github.RequestService
	author                 string
	pushToDepot            packagePusher
	versionDeps            depsVersioner
	createDepotRepo        depotRepoCreator
	downloadPackage        packageDownloader
	destroyDepotRepo       depotRepoDestroyer
	isPackageArchived      packageArchivalChecker
	constructionZonePath   string
	recordPackageArchival  packageArchivalRecorder
	attemptWorkDirDeletion workDirDeletionAttempter
}

// packageVersioner is responsible for versioning a downloaded package.
type packageVersioner func(args packageVersionerArgs) error

// packageDownloaderArgs is the arguments struct for packageDownloader.
type packageDownloaderArgs struct {
	io                   io.IO
	author               string
	repo                 string
	sha                  string
	doHTTPGet            httpGetter
	unzipArchive         archiveUnzipper
	deleteWorkDir        workDirDeletionAttempter
	constructionZonePath string
}

// packageDownloadPaths is a tuple of downloaded package paths.
type packageDownloadPaths struct {
	workDirPath    string
	archiveDirPath string
}

// packageDownloader is responsible for downloading, unzipping, and writing
// package to constructionZonePath. Returns downloaded package directory path.
type packageDownloader func(args packageDownloaderArgs) (packageDownloadPaths, error)

// packagePusherArgs is the arguments struct for packagePusher.
type packagePusherArgs struct {
	author       string
	repo         string
	sha          string
	creds        *config.Credentials
	packagePaths packageDownloadPaths
	gitClient    git.Client
}

// dbPackageArchivalChecker returns true if a package version matching the
// parameters exists in the database.
type dbPackageArchivalChecker func(
	q db.Queryable,
	author string,
	repo string,
	sha string) (bool, error)

// httpGetter executes an HTTP get to the specified URL and returns the
// corresponding response.
type httpGetter func(url string) (*http.Response, error)

// packagePusher is responbile for pushing package to depot.
type packagePusher func(args packagePusherArgs) error

// depsVersioner is responsible for versioning the dependencies in a package.
type depsVersioner func(args verdeps.VersionDepsArgs) error

// archiveUnzipper unzips a zip archive.
type archiveUnzipper func(archive, target string) error

// depotRepoCreator creates a repository in depot in accordance to the author,
// repo and sha specified. Returns true if the repo was created by this func.,
// or returns false is the the directory already existed.
type depotRepoCreator func(author, repo, sha string) (bool, error)

// depotRepoDestroyer destroys a repository in depot according to the author,
// repo and sha.
type depotRepoDestroyer func(author, repo, sha string) error

// depotExistenceChecker checks if a package matching author, repo and sha
// exists in depot.
type depotExistenceChecker func(author, repo, sha string) (bool, error)

// workDirDeletionAttempter attempts to delete a working directory. If it fails,
// instead of returning the error, it logs the problem and moves on. Functions
// implementing this type are designed to run in go-routines and defers.
type workDirDeletionAttempter func(workDirPath string)
