package github

import (
	"fmt"
	"log"
	"net/http"
)

// CheckGitHubRepoExists returns whether a repo exists
// TODO(Shikkic): Instead of pinging try downloading refs, might be more sustainable?
func (gitHubRequestService *RequestService) CheckGitHubRepoExists(author, repo string) error {
	repoName := BuildNewGitHubRepoName(author, repo)
	// TODO change this to fetch ref
	url := fmt.Sprintf("https://github.com/%s/%s", GitHubGophrPackageOrgName, repoName)
	resp, err := http.Get(url)

	if err != nil {
		log.Println("Error occurred during request")
		return err
	}

	if resp.StatusCode == 404 {
		log.Printf("No Github repo exists in %s org with the name %s \n", GitHubGophrPackageOrgName, repoName)
		return nil
	}

	return fmt.Errorf("Error status code %d, a repo with that name already exists.", resp.StatusCode)
}
