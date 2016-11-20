package github

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/DataDog/datadog-go/statsd"
	"github.com/gophr-pm/gophr/lib/datadog"
	"github.com/gophr-pm/gophr/lib/dtos"
)

// ddEventFetchRepoData is the name of the custom datadog event for this
// function.
const ddEventFetchRepoData = "github.fetch-repo-data"

// FetchRepoData fetches the Github repository metadata for the specified
// package.
func (svc *requestServiceImpl) FetchRepoData(
	author string,
	repo string,
) (dtos.GithubRepo, error) {
	// Specify monitoring parameters.
	trackingArgs := datadog.TrackTransactionArgs{
		Tags:      []string{"github", datadog.TagInternal},
		Client:    svc.ddClient,
		AlertType: datadog.Success,
		StartTime: time.Now(),
		EventInfo: []string{fmt.Sprintf(
			`{ author: "%s", repo: "%s" }`,
			author,
			repo,
		)},
		MetricName:      datadog.MetricJobDuration,
		CreateEvent:     statsd.NewEvent,
		CustomEventName: ddEventFetchRepoData,
	}

	// Ensure that the transaction is tracked after the job finishes.
	defer datadog.TrackTransaction(trackingArgs)

	log.Printf(`Fetching Github repository data for "%s/%s".
`, author, repo)

	for attempts := 0; attempts < githubAPIAttemptsLimit; attempts++ {
		resp, err := svc.keyChain.acquireKey().getFromGithub(
			buildGitHubRepoDataAPIURL(
				author,
				repo))

		// Make sure that the response body gets closed eventually.
		defer resp.Body.Close()

		if err != nil {
			err = fmt.Errorf(
				`Failed to get repo data for "%s/%s": %v.`,
				author,
				repo,
				err)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return dtos.GithubRepo{}, err
		}

		// Handle all kinds of failures.
		if resp.StatusCode == 404 {
			err = fmt.Errorf(
				`Failed to get repo data for "%s/%s": package not found.`,
				author,
				repo)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return dtos.GithubRepo{}, err
		} else if resp.StatusCode == 403 {
			// If there was a forbidden status code, try the request again.
			continue
		} else if resp.StatusCode != 200 && resp.StatusCode != 304 {
			err = fmt.Errorf(
				`Failed to get repo data for "%s/%s": bumped into a status code %d.`,
				author,
				repo,
				resp.StatusCode)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return dtos.GithubRepo{}, err
		}

		repoData, err := parseGitHubRepoDataResponseBody(resp)
		if err != nil {
			err = fmt.Errorf(
				`Failed to parse repo data for "%s/%s": %v.`,
				author,
				repo,
				err)

			// Make sure that the error is recorded in the datadog transaction.
			trackingArgs.AlertType = datadog.Error
			trackingArgs.EventInfo = append(trackingArgs.EventInfo, err.Error())

			return dtos.GithubRepo{}, err
		}

		return repoData, nil
	}

	return dtos.GithubRepo{}, fmt.Errorf(
		`Failed to get repo data for "%s/%s": `+
			`all %d attempts failed.`,
		author,
		repo,
		githubAPIAttemptsLimit)
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
