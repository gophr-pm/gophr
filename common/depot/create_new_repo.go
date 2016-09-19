package depot

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	git "github.com/libgit2/git2go"
)

// repoCreationAttemptsLimit sets the cap on how many times repo creation is
// attempted.
const repoCreationAttemptsLimit = 3

// CreateNewRepo creates a new repo in the depot *record scratch*. If the repo
// was created by this function invocation, repoCreated will return true. If the
// the repo was created by something else, repoCreated will be false.
func CreateNewRepo(author string, repo string, sha string) (bool, error) {
	log.Printf("Creating New Repo on depot %s/%s@%s \n", author, repo, sha)

	// Create the repo dir out here so that in can be used after the for loop.
	repoDir := fmt.Sprintf("%s.git", BuildHashedRepoName(author, repo, sha))

	// Try to create the repo a few times.
	for attempts := 0; attempts < repoCreationAttemptsLimit; attempts = attempts + 1 {
		// First, check if repo dir exists on the depot volume.
		if exists, err := repoDirExists(repoDir); err != nil {
			return false, fmt.Errorf("Failed to check if repo directory \"%s\" exists: %v.", repoDir, err)
		} else if exists {
			// If the repo directory exists already, exit - the job is already done.
			return false, nil
		}

		// The repo directory doesn't exist, so creating it should be ok.
		if err := os.Mkdir(filepath.Join(depotReposPath, repoDir), 0644); err == nil {
			// The folder was created just fine. Now create the bare git repo.
			if _, err = git.InitRepository(filepath.Join(depotReposPath, repoDir), true); err != nil {
				return false, fmt.Errorf("Could not initialize new repository: %v.", err)
			}

			// Woop! New repo in the depot!
			return true, nil
		}

		// Merely log mkdir failure so that it can be re-attempted.
		log.Printf("Failed to create repo directory \"%s\" in the depot.\n", repoDir)
	}

	return false, fmt.Errorf(
		"After %d attempts, failed to create new repo \"%s\".",
		repoCreationAttemptsLimit,
		repoDir)
}

// repoDirExists returns true if the repoDir exists, and false otherwise. If
// some other error occurs, returns that error.
func repoDirExists(repoDir string) (bool, error) {
	// TODO(skeswa): make the depotReposPath configurable (could change).
	path := filepath.Join(depotReposPath, repoDir)
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}
