package main

import (
	"fmt"
	"log"

	"github.com/skeswa/gophr/common/depot"
	"github.com/skeswa/gophr/common/verdeps"
)

// versionAndArchivePackage creates a github repo for the packageModel on
// github.com/gophr/gophr-packages versioned a the specified args.sha.
func versionAndArchivePackage(args packageVersionerArgs) error {
	log.Printf("Preparing to sub-version %s/%s@%s \n", args.author, args.repo, args.sha)

	// Create a new repository in the depot before pushing to it.
	if err := depot.CreateNewRepo(
		args.author,
		args.repo,
		args.sha,
	); err != nil {
		return err
	}

	// Download the package in the construction zone.
	downloadPaths, err := downloadPackage(packageDownloaderArgs{
		sha:                  args.sha,
		repo:                 args.repo,
		author:               args.author,
		constructionZonePath: args.constructionZonePath,
	})
	if err != nil {
		return err
	}

	// Perform clean-up after function exits.
	defer deleteFolder(downloadPaths.workDirPath)

	// Version lock all of the Github dependencies in the packageModel.
	if err = args.versionDeps(verdeps.VersionDepsArgs{
		SHA:           args.sha,
		Repo:          args.repo,
		Path:          downloadPaths.archiveDirPath,
		Author:        args.author,
		GithubService: args.githubRequestService,
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
