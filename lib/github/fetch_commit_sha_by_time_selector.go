package github

import (
	"fmt"
	"time"
)

// fetchCommitSHAByTimeSelector uses the provided time to find the closest
// commit SHA using the Github API.
func (svc *requestServiceImpl) fetchCommitSHAByTimeSelector(
	author string,
	repo string,
	timeStamp time.Time,
	timeSelector string,
) (string, error) {
	for attempts := 0; attempts < githubAPIAttemptsLimit; attempts++ {
		resp, err := svc.keyChain.acquireKey().getFromGithub(
			buildGitHubRepoCommitsFromTimestampAPIURL(
				author,
				repo,
				timeStamp,
				timeSelector))

		// Make sure that the response body gets closed eventually.
		defer resp.Body.Close()

		if err != nil {
			return "", fmt.Errorf(
				`Failed to get commit SHA for "%s/%s" by time selector: %v.`,
				author,
				repo,
				err)
		}

		// Handle all kinds of failures.
		if resp.StatusCode == 404 {
			return "", fmt.Errorf(
				`Failed to get commit SHA for "%s/%s" by time selector: `+
					`package not found.`,
				author,
				repo)
		} else if resp.StatusCode == 403 {
			// If there was a forbidden status code, try the request again.
			continue
		} else if resp.StatusCode != 200 && resp.StatusCode != 304 {
			return "", fmt.Errorf(
				`Failed to get commit SHA for "%s/%s" by time selector: `+
					`bumped into a status code %d.`,
				author,
				repo,
				resp.StatusCode)
		}

		commitSHA, err := parseGitHubCommitTimestamp(resp)
		if err != nil {
			return "", fmt.Errorf(
				`Failed to parse commit SHA for "%s/%s" by time selector: %v.`,
				author,
				repo,
				err)
		}

		return commitSHA, nil
	}

	return "", fmt.Errorf(
		`Failed to get commit SHA for "%s/%s" by time selector: `+
			`all %d attempts failed.`,
		author,
		repo,
		githubAPIAttemptsLimit)
}
