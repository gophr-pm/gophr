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

// FetchCommitSHA fetches the commit SHA that is chronologically closest to a
// given timestamp.
func (svc *requestServiceImpl) FetchCommitSHA(
	author string,
	repo string,
	timestamp time.Time,
) (string, error) {
	log.Printf(`Fetching Github commit SHA of "%s/%s" for timestamp %s.
`, author, repo, timestamp.String())

	// Fetch commits chronologically before the timestamp.
	commitSHA, err := svc.fetchCommitSHAByTimeSelector(
		author,
		repo,
		timestamp,
		commitsUntilParameter)
	if err == nil {
		return commitSHA, nil
	}

	// Fetch commits chronologically after the timestamp.
	commitSHA, err = svc.fetchCommitSHAByTimeSelector(
		author,
		repo,
		timestamp,
		commitsAfterParameter)
	if err == nil {
		return commitSHA, nil
	}

	// Before and after failed somehow. Just take the latest.
	refs, err := lib.FetchRefs(author, repo)
	if err != nil {
		return refs.MasterRefHash, nil
	}

	return "", err
}

func buildGitHubRepoCommitsFromTimestampAPIURL(
	author string,
	repo string,
	timestamp time.Time,
	timeSelector string,
) string {
	url := fmt.Sprintf(
		"https://api.github.com/repos/%s/%s/commits?%s=%s",
		author,
		repo,
		timeSelector,
		strings.Replace(timestamp.String(), " ", "", -1))

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
