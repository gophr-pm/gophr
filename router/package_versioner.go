package main

import (
	"fmt"
	"log"

	"github.com/skeswa/gophr/common/depot"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/verdeps"
)

// versionAndArchivePackage creates a github repo for the packageModel on
// github.com/gophr/gophr-packages versioned a the specified args.sha.
func versionAndArchivePackage(args packageVersionerArgs) error {
	log.Printf("Preparing to sub-version %s/%s@%s \n", args.author, args.repo, args.sha)

	if err := depot.CreateNewRepo(
		args.author,
		args.repo,
		args.sha,
	); err != nil {
		return err
	}

	downloadPaths, err := downloadPackage(packageDownloaderArgs{
		author:               args.author,
		repo:                 args.repo,
		sha:                  args.sha,
		constructionZonePath: args.constructionZonePath,
	})
	if err != nil {
		return err
	}

	// Perform clean-up after function exits.
	defer deleteFolder(downloadPaths.workDirPath)

	// Instantiate New Github Request Service.
	// TODO intialize in handler.
	gitHubRequestService := github.NewRequestService(
		github.RequestServiceParams{
			ForIndexer: false,
			Conf:       args.conf,
			Session:    args.db,
		},
	)

	// Version lock all of the Github dependencies in the packageModel.
	if err = verdeps.VersionDeps(
		verdeps.VersionDepsArgs{
			SHA:           args.sha,
			Repo:          args.repo,
			Path:          downloadPaths.archiveDirPath,
			Author:        args.author,
			GithubService: gitHubRequestService,
		}); err != nil {
		return fmt.Errorf("Could not version deps properly: %v.", err)
	}

	// Push versioned package to depot, then delete the package directory from
	// the construction zone.
	if err = args.pushToDepot(packagePusherArgs{
		author:       args.author,
		repo:         args.repo,
		sha:          args.sha,
		creds:        args.creds,
		packagePaths: downloadPaths,
	}); err != nil {
		return fmt.Errorf("Could not push versioned package to depot: %v.", err)
	}

	// Record that this package has been archived.
	go args.recordPackageArchival(packageArchivalArgs{
		db:     args.db,
		sha:    args.sha,
		repo:   args.repo,
		author: args.author,
	})

	return nil
}
