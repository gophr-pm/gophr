package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/pquerna/ffjson/ffjson"
)

// ddEventFetchCommitSHA is the name of the custom datadog event for this
// function.
const ddEventFetchCommitSHA = "github.fetch-commit-sha"

// FetchCommitSHA fetches the commit SHA that is chronologically closest to a
// given timestamp.
func (svc *requestServiceImpl) FetchCommitSHA(
	author string,
	repo string,
	timestamp time.Time,
) (string, error) {
	// Specify monitoring parameters.
	trackingArgs := datadog.TrackTransactionArgs{
		Tags:      []string{"github", datadog.TagInternal},
		Client:    svc.ddClient,
		AlertType: datadog.Success,
		StartTime: time.Now(),
		EventInfo: []string{fmt.Sprintf(
			`{ author: "%s", repo: "%s", timestamp: "%s" }`,
			author,
			repo,
			timestamp.Format(time.RFC3339),
		)},
		MetricName:      datadog.MetricJobDuration,
		CreateEvent:     statsd.NewEvent,
		CustomEventName: ddEventFetchCommitSHA,
	}

	// Ensure that the transaction is tracked after the job finishes.
	defer datadog.TrackTransaction(trackingArgs)

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

	// Make sure that the anomaly is recorded in the datadog transaction.
	trackingArgs.EventInfo = append(
		trackingArgs.EventInfo,
		`Failed to get the commit SHA both before and after the timestamp.`)

	// Before and after failed somehow. Just take the latest.
	refs, err := lib.FetchRefs(author, repo)
	if err != nil {
		// Make sure that the error is recorded in the datadog transaction.
		trackingArgs.AlertType = datadog.Error
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

		return refs.MasterRefHash, nil
	}

	// Make sure that the error is recorded in the datadog transaction.
	if err != nil {
		trackingArgs.AlertType = datadog.Error
		trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())
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
