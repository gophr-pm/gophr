package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/pquerna/ffjson/ffjson"
)

// FetchCommitTimestamp fetches the timestamp of a commit from Github API.
func (svc *requestService) FetchCommitTimestamp(
	author string,
	repo string,
	sha string,
) (time.Time, error) {
	log.Printf(`Fetching Github commit timestamp for "%s/%s@%s".
`, author, repo, sha)

	resp, err := svc.keyChain.acquireKey().getFromGithub(
		buildGitHubCommitTimestampAPIURL(
			author,
			repo,
			sha))

	// Make sure that the response body gets closed eventually.
	defer resp.Body.Close()

	if err != nil {
		return time.Time{}, fmt.Errorf(
			`Failed to get timestamp for commit "%s/%s%s": %v.`,
			author,
			repo,
			sha,
			err)
	}

	// Handle all kinds of failures.
	if resp.StatusCode == 404 {
		return time.Time{}, fmt.Errorf(
			`Failed to get timestamp for commit "%s/%s%s": package not found.`,
			author,
			repo,
			sha)
	} else if resp.StatusCode != 200 && resp.StatusCode != 304 {
		return time.Time{}, fmt.Errorf(
			`Failed to get timestamp for commit "%s/%s%s": `+
				`bumped into a status code %d.`,
			author,
			repo,
			sha,
			resp.StatusCode)
	}

	timeStamp, err := parseGitHubCommitLookUpResponseBody(resp)
	if err != nil {
		return time.Time{}, fmt.Errorf(
			`Failed to parse timestamp for commit "%s/%s%s": %v.`,
			author,
			repo,
			sha,
			err)
	}

	return timeStamp, nil
}

func buildGitHubCommitTimestampAPIURL(
	author string,
	repo string,
	sha string) string {
	return fmt.Sprintf("https://api.github.com/repos/%s/%s/commits/%s",
		author,
		repo,
		sha,
	)
}

func parseGitHubCommitLookUpResponseBody(
	response *http.Response,
) (time.Time, error) {
	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return time.Time{}, errors.New("Failed to parse response body")
	}

	var commitLookUpDTO dtos.GithubCommitLookUp
	err = ffjson.Unmarshal(body, &commitLookUpDTO)
	if err != nil {
		return time.Time{}, errors.New("Failed to unmarshal response body")
	}

	if commitLookUpDTO.Commit != nil && commitLookUpDTO.Commit.Committer != nil {
		return commitLookUpDTO.Commit.Committer.Date, nil
	}

	return time.Time{}, errors.New("No commit timestamp available for given SHA")
}
