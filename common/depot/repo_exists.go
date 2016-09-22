package depot

import (
	"os"
	"path/filepath"
)

// RepoExists returns true if the repoDir exists, and false otherwise. If
// some other error occurs, returns that error.
func RepoExists(author, repo, sha string) (bool, error) {
	// Compose author, repo and sha together to get the repo name.
	repoDir := BuildHashedRepoName(author, repo, sha)

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
