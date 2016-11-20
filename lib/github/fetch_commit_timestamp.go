package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/pquerna/ffjson/ffjson"
)

// ddEventFetchCommitTimestamp is the name of the custom datadog event for this
// function.
const ddEventFetchCommitTimestamp = "github.fetch-commit-timestamp"

// FetchCommitTimestamp fetches the timestamp of a commit from Github API.
func (svc *requestServiceImpl) FetchCommitTimestamp(
	author string,
	repo string,
	sha string,
) (time.Time, error) {
	// Specify monitoring parameters.
	trackingArgs := datadog.TrackTransactionArgs{
		Tags:      []string{"github", datadog.TagInternal},
		Client:    svc.ddClient,
		AlertType: datadog.Success,
		StartTime: time.Now(),
		EventInfo: []string{fmt.Sprintf(
			`{ author: "%s", repo: "%s", sha: "%s" }`,
			author,
			repo,
			sha,
		)},
		MetricName:      datadog.MetricJobDuration,
		CreateEvent:     statsd.NewEvent,
		CustomEventName: ddEventFetchCommitTimestamp,
	}

	// Ensure that the transaction is tracked after the job finishes.
	defer datadog.TrackTransaction(trackingArgs)

	log.Printf(`Fetching Github commit timestamp for "%s/%s@%s".
`, author, repo, sha)

	for attempts := 0; attempts < githubAPIAttemptsLimit; attempts++ {
		resp, err := svc.keyChain.acquireKey().getFromGithub(
			buildGitHubCommitTimestampAPIURL(
				author,
				repo,
				sha))

		// Make sure that the response body gets closed eventually.
		defer resp.Body.Close()

		if err != nil {
			err = fmt.Errorf(
				`Failed to get timestamp for commit "%s/%s%s": %v.`,
				author,
				repo,
				sha,
				err)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return time.Time{}, err
		}

		// Handle all kinds of failures.
		if resp.StatusCode == 404 {
			err = fmt.Errorf(
				`Failed to get timestamp for commit "%s/%s%s": package not found.`,
				author,
				repo,
				sha)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return time.Time{}, err
		} else if resp.StatusCode == 403 {
			// If there was a forbidden status code, try the request again.
			continue
		} else if resp.StatusCode != 200 && resp.StatusCode != 304 {
			err = fmt.Errorf(
				`Failed to get timestamp for commit "%s/%s%s": `+
					`bumped into a status code %d.`,
				author,
				repo,
				sha,
				resp.StatusCode)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return time.Time{}, err
		}

		timeStamp, err := parseGitHubCommitLookUpResponseBody(resp)
		if err != nil {
			err = fmt.Errorf(
				`Failed to parse timestamp for commit "%s/%s%s": %v.`,
				author,
				repo,
				sha,
				err)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return time.Time{}, err
		}

		return timeStamp, nil
	}

	return time.Time{}, fmt.Errorf(
		`Failed to fetch timestamp for commit "%s/%s%s": `+
			`all %d attempts failed.`,
		author,
		repo,
		sha,
		githubAPIAttemptsLimit)
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
