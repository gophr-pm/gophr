package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/pquerna/ffjson/ffjson"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/dtos"
)

// FetchCommitSHA Fetches a commitSHA closest to a given timestamp
func (gitHubRequestService *GitHubRequestService) FetchCommitSHA(
	author string,
	repo string,
	timestamp time.Time,
) (string, error) {
	commitSHA, err := gitHubRequestService.fetchCommitSHAByTimeSelector(author, repo, timestamp, commitsUntilParameter)
	if err == nil {
		return commitSHA, nil
	}

	log.Printf("%s \n", err)
	commitSHA, err = gitHubRequestService.fetchCommitSHAByTimeSelector(author, repo, timestamp, commitsAfterParameter)
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

func (gitHubRequestService *GitHubRequestService) fetchCommitSHAByTimeSelector(
	author string,
	repo string,
	timestamp time.Time,
	timeSelector string,
) (string, error) {
	APIKeyModel := gitHubRequestService.APIKeyChain.getAPIKeyModel()
	log.Println(APIKeyModel)
	log.Printf("%+v \n", APIKeyModel)
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
	APIKeyModel.print()

	commitSHA, err := parseGitHubCommitTimestamp(resp)
	if err != nil {
		return "", err
	}

	return commitSHA, nil
}

func buildGitHubRepoCommitsFromTimestampAPIURL(
	author string,
	repo string,
	APIKeyModel GitHubAPIKeyModel,
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

	var commitSHAArray []dtos.GitCommitDTO
	err = ffjson.Unmarshal(body, &commitSHAArray)
	if err != nil {
		return "", errors.New("Failed to unmarshal response body")
	}

	if len(commitSHAArray) >= 1 {
		return commitSHAArray[0].SHA, nil
	}

	return "", errors.New("No commit SHAs available for timestamp given")
}
