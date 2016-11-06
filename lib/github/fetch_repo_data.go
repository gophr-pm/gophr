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

// FetchGitHubDataForPackageModel fetches current repo data of a given
// package.
// TODO optimize this with FFJSON models
func (svc *requestService) FetchGitHubDataForPackageModel(
	author string,
	repo string,
) (dtos.GithubRepo, error) {
	APIKeyModel := svc.APIKeyChain.getAPIKeyModel()
	githubURL := buildGitHubRepoDataAPIURL(author, repo, *APIKeyModel)
	log.Printf("Fetching GitHub data for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()

	if err != nil {
		return dtos.GithubRepo{}, errors.New("Request error.")
	}

	if resp.StatusCode == 404 {
		log.Println("PackageModel was not found on Github")
		return dtos.GithubRepo{}, nil
	}

	APIKeyModel.incrementUsageFromResponseHeader(resp.Header)

	responseBodyMap, err := parseGitHubRepoDataResponseBody(resp)
	if err != nil {
		return dtos.GithubRepo{}, err
	}

	return responseBodyMap, nil
}

func buildGitHubRepoDataAPIURL(
	author string,
	repo string,
	keyModel APIKeyModel,
) string {
	url := fmt.Sprintf(
		"%s/repos/%s/%s?access_token=%s",
		GitHubBaseAPIURL,
		author,
		repo,
		keyModel.Key)
	return url
}

// TODO Optimize this with ffjson struct!
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
