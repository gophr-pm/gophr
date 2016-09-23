package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/skeswa/gophr/common/depot"
)

// repoDestructionAttemptsLimit sets the cap on how many times repo destruction
// is attempted.
const repoDestructionAttemptsLimit = 3

// destroyRepo creates a new repo in the depot *record scratch*. If the repo
// was created by this function invocation, repoCreated will return true. If the
// the repo was created by something else, repoCreated will be false.
func destroyRepo(depotReposPath, author, repo, sha string) error {
	log.Printf("Deleting repo on depot %s/%s@%s\n", author, repo, sha)

	// Create the repo dir out here so that in can be used after the for loop.
	repoDir := fmt.Sprintf("%s.git", depot.BuildHashedRepoName(author, repo, sha))

	// Try to create the repo a few times.
	for attempts := 0; attempts < repoDestructionAttemptsLimit; attempts = attempts + 1 {
		// First, check if repo dir exists on the depot volume.
		if exists, err := repoExists(depotReposPath, author, repo, sha); err != nil {
			return fmt.Errorf("Failed to check if repo directory \"%s\" exists: %v.", repoDir, err)
		} else if !exists {
			// If the repo directory doesn't exist already, exit - the job is already
			// done.
			return nil
		}

		// The repo directory doesn't exist, so creating it should be ok.
		if err := os.RemoveAll(filepath.Join(depotReposPath, repoDir)); err == nil {
			// Ding dong the witch is dead!
			return nil
		}

		// Merely log mkdir failure so that it can be re-attempted.
		log.Printf("Failed to delete repo directory \"%s\" in the depot.\n", repoDir)
	}

	return fmt.Errorf(
		"After %d attempts, failed to destroy repo \"%s\".",
		repoDestructionAttemptsLimit,
		repoDir)
}
