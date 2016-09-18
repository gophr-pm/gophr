package depot

import (
	"fmt"
	"log"
	"os"

	git "github.com/libgit2/git2go"
)

const (
	filePerm = 0644
)

// CreateNewRepo if repo doesn't already exist will create a new
// repo in depot
func CreateNewRepo(author string, repo string, sha string) error {
	log.Printf("Creating New Repo on depot %s/%s@%s \n", author, repo, sha)
	folderName := fmt.Sprintf("%s.git", BuildHashedRepoName(author, repo, sha))

	// First check if repo folder exists on depot volume
	if err := checkIfRepoFolderExists(folderName); err == nil {
		return nil
	}

	// Create repo folder on depot volume
	err := os.Mkdir(
		fmt.Sprintf("%s/%s", depotReposPath, folderName),
		filePerm,
	)

	// Check to make sure we properly created folder
	if err != nil {
		// If we couldn't make the folder, and the foler exists now return error. Someone else is
		// Already versioning this package
		if checkFolderErr := checkIfRepoFolderExists(folderName); checkFolderErr == nil {
			return fmt.Errorf("Error, folder already exists, package is in process of being versioned already. %v %v", checkFolderErr, err)
		}
	}

	// Initialize the new bare git repository
	_, err = git.InitRepository(fmt.Sprintf("%s/%s", depotReposPath, folderName), true)
	if err != nil {
		return fmt.Errorf("Error, could not initialize new repository. %v", err)
	}

	return err
}

func checkIfRepoFolderExists(folderName string) error {
	_, err := os.Stat(fmt.Sprintf("%s/%s", depotReposPath, folderName))
	return err
}
