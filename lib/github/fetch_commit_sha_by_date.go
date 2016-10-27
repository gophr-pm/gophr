package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/pquerna/ffjson/ffjson"
)

// FetchCommitSHA Fetches a commitSHA closest to a given timestamp
func (svc *requestService) FetchCommitSHA(
	author string,
	repo string,
	timestamp time.Time,
) (string, error) {
	commitSHA, err := svc.fetchCommitSHAByTimeSelector(author, repo, timestamp, commitsUntilParameter)
	if err == nil {
		return commitSHA, nil
	}

	log.Printf("%s \n", err)
	commitSHA, err = svc.fetchCommitSHAByTimeSelector(author, repo, timestamp, commitsAfterParameter)
	if err == nil {
		return commitSHA, nil
	}

	log.Printf("%s \n", err)
	refs, err := common.FetchRefs(author, repo)
	if err != nil {
		return refs.MasterRefHash, nil
	}

	return "", err
}

func (svc *requestService) fetchCommitSHAByTimeSelector(
	author string,
	repo string,
	timestamp time.Time,
	timeSelector string,
) (string, error) {
	APIKeyModel := svc.APIKeyChain.getAPIKeyModel()
	log.Printf("Determining APIKey %s \n", APIKeyModel.Key)

	githubURL := buildGitHubRepoCommitsFromTimestampAPIURL(author, repo, *APIKeyModel, timestamp, timeSelector)
	log.Printf("Fetching GitHub data for %s \n", githubURL)

	resp, err := http.Get(githubURL)
	defer resp.Body.Close()

	if err != nil {
		return "", errors.New("Request error.")
	}

	if resp.StatusCode == 404 {
		log.Println("PackageModel was not found on Github")
		return "", nil
	}

	APIKeyModel.incrementUsageFromResponseHeader(resp.Header)

	commitSHA, err := parseGitHubCommitTimestamp(resp)
	if err != nil {
		return "", err
	}

	return commitSHA, nil
}

func buildGitHubRepoCommitsFromTimestampAPIURL(
	author string,
	repo string,
	APIKeyModel APIKeyModel,
	timestamp time.Time,
	timeSelector string,
) string {
	url := fmt.Sprintf("%s/repos/%s/%s/commits?%s=%s&access_token=%s",
		GitHubBaseAPIURL,
		author,
		repo,
		timeSelector,
		strings.Replace(timestamp.String(), " ", "", -1),
		APIKeyModel.Key,
	)
	return url
}

func parseGitHubCommitTimestamp(response *http.Response) (string, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return "", errors.New("Failed to parse response body")
	}

	var commitSHAArray []dtos.GithubCommit
	err = ffjson.Unmarshal(body, &commitSHAArray)
	if err != nil {
		return "", errors.New("Failed to unmarshal response body")
	}

	if len(commitSHAArray) >= 1 {
		return commitSHAArray[0].SHA, nil
	}

	return "", errors.New("No commit SHAs available for timestamp given")
}
