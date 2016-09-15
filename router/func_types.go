package main

import (
	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
)

// refsDownloader is responsible for downloading the git refs for a package.
type refsDownloader func(author, repo string) (common.Refs, error)

// packageDownloadRecorderArgs is the arguments struct for
// packageDownloadRecorders.
type packageDownloadRecorderArgs struct {
	db      *gocql.Session
	sha     string
	repo    string
	author  string
	version string
}

// packageDownloadRecorder is responsible for recording package downloads. If
// there is a problem while recording, then the error is logged instead of
// bubbled.
type packageDownloadRecorder func(args packageDownloadRecorderArgs)

// packageArchivalArgs is the arguments struct for packageArchivalRecorders and
// packageArchivalCheckers.
type packageArchivalArgs struct {
	db     *gocql.Session
	sha    string
	repo   string
	author string
}

// packageArchivalRecorder is responsible for recording package archival. If
// there is a problem while recording, then the error is logged instead of
// bubbled.
type packageArchivalRecorder func(args packageArchivalArgs)

// packageArchivalChecker is responsible for checking whether a package has been
// archived or not. Returns true if the package has been archived, and false
// otherwise.
type packageArchivalChecker func(args packageArchivalArgs) (bool, error)

// packageVersionerArgs is the arguments struct for packageVersioners.
type packageVersionerArgs struct {
	db                    *gocql.Session
	sha                   string
	repo                  string
	conf                  *config.Config
	creds                 *config.Credentials
	author                string
	constructionZonePath  string
	recordPackageArchival packageArchivalRecorder
}

type packageVersioner func(args packageVersionerArgs) error
