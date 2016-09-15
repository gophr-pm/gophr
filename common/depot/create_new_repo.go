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
func CreateNewRepo(author string, repo string, ref string) error {
	log.Println("Creating New Repo")
	folderName := fmt.Sprintf(
		"%s.git",
		BuildHashedRepoName(author, repo, ref),
	)
	if err := checkIfRepoFolderExists(folderName); err == nil {
		return nil
	}

	err := os.Mkdir(
		fmt.Sprintf("%s/%s", depotReposPath, folderName),
		filePerm,
	)

	// TODO(Shikkic): If we can't make the repo folder bail, that means someone else is already versioning
	if err != nil {
		if checkIfRepoFolderExists(folderName); err != nil {
			return fmt.Errorf("Error, could not create folder or verify that it currently exists. %v", err)
		}
	}

	// Git init bare
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
