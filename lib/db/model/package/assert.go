package pkg

import (
	"fmt"
	"sync"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
)

// AssertExistence asserts that a package exists.
func AssertExistence(
	q db.Queryable,
	author string,
	repo string,
	ghSvc github.RequestService,
) error {
	var (
		err   error
		count int
	)

	// Create the query to check if this package version exists.
	if err = query.SelectCount().
		From(packagesTableName).
		Where(query.Column(packagesColumnNameAuthor).Equals(author)).
		And(query.Column(packagesColumnNameRepo).Equals(repo)).
		Limit(1).
		Create(q).
		Scan(&count); err != nil {
		return fmt.Errorf(
			"Failed to check if package %s/%s exists: %v",
			author,
			repo,
			err)
	}

	// If this package version doesn't exist, then make it exist.
	if count < 1 {
		var (
			wg                 sync.WaitGroup
			awesome            bool
			repoData           dtos.GithubRepo
			awesomeCheckError  error
			repoDataFetchError error
		)

		// Start two workers that get package github metadata, and whether the
		// package is awesome.
		wg.Add(2)
		go checkIfAwesomeAsynchronously(
			q,
			author,
			repo,
			&awesome,
			&awesomeCheckError,
			&wg)
		go getGithubRepoDataAsynchronously(
			ghSvc,
			author,
			repo,
			&repoData,
			&repoDataFetchError,
			&wg)

		// Wait, then handle the outputs.
		wg.Wait()
		if awesomeCheckError != nil {
			return fmt.Errorf(
				"Failed to check if package %s/%s is awesome: %v",
				author,
				repo,
				awesomeCheckError)
		}
		if repoDataFetchError != nil {
			return fmt.Errorf(
				"Failed to fetch repo data for package %s/%s: %v",
				author,
				repo,
				repoDataFetchError)
		}

		// Now that we have all the requisite data, insert the new package.
		if err = Insert(InsertArgs{
			Repo:        repo,
			Stars:       repoData.Stars,
			Author:      author,
			Awesome:     awesome,
			Queryable:   q,
			Description: repoData.Description,
		}); err != nil {
			return err
		}
	}

	return nil
}
