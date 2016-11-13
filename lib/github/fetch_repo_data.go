package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib/dtos"
)

// FetchRepoData fetches the Github repository metadata for the specified
// package.
func (svc *requestServiceImpl) FetchRepoData(
	author string,
	repo string,
) (dtos.GithubRepo, error) {
	log.Printf(`Fetching Github repository data for "%s/%s".
`, author, repo)

	resp, err := svc.keyChain.acquireKey().getFromGithub(
		buildGitHubRepoDataAPIURL(
			author,
			repo))

	// Make sure that the response body gets closed eventually.
	defer resp.Body.Close()

	if err != nil {
		return dtos.GithubRepo{}, fmt.Errorf(
			`Failed to get repo data for "%s/%s": %v.`,
			author,
			repo,
			err)
	}

	// Handle all kinds of failures.
	if resp.StatusCode == 404 {
		return dtos.GithubRepo{}, fmt.Errorf(
			`Failed to get repo data for "%s/%s": package not found.`,
			author,
			repo)
	} else if resp.StatusCode != 200 && resp.StatusCode != 304 {
		return dtos.GithubRepo{}, fmt.Errorf(
			`Failed to get repo data for "%s/%s": bumped into a status code %d.`,
			author,
			repo,
			resp.StatusCode)
	}

	repoData, err := parseGitHubRepoDataResponseBody(resp)
	if err != nil {
		return dtos.GithubRepo{}, fmt.Errorf(
			`Failed to parse repo data for "%s/%s": %v.`,
			author,
			repo,
			err)
	}

	return repoData, nil
}

func buildGitHubRepoDataAPIURL(
	author string,
	repo string,
) string {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s",
		author,
		repo)
	return url
}

func parseGitHubRepoDataResponseBody(
	response *http.Response,
) (dtos.GithubRepo, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return dtos.GithubRepo{}, errors.New("Failed to parse response body")
	}

	var repoData dtos.GithubRepo
	if err = json.Unmarshal(body, &repoData); err != nil {
		return dtos.GithubRepo{}, errors.New("Failed to unmarshal response body")
	}

	return repoData, nil
}
