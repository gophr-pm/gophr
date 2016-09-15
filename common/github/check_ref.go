package github

import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

const (
	refsFetchURLTemplate           = "http://depot-svc:3000/%s.git/info/refs?service=git-upload-pack"
	errorRefsFetchNetworkFailure   = "Could not reach Github at the moment; Please try again later"
	errorRefsFetchNoSuchRepo       = "Could not find a Github repository at %s"
	errorRefsFetchGithubError      = "Github responded with an error: %v"
	errorRefsFetchGithubParseError = "Cannot read refs from Github: %v"
)

// CheckIfRefExists checks whether a given ref exists in the remote refs list.
func CheckIfRefExists(author, repo string, ref string) (bool, error) {
	ref = BuildGitHubBranch(ref)
	repoName := BuildNewGitHubRepoName(author, repo)

	res, err := httpClient.Get(fmt.Sprintf(refsFetchURLTemplate, repoName))
	if err != nil {
		return false, errors.New(errorRefsFetchNetworkFailure)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return false, nil
	} else if res.StatusCode >= 500 {
		// FYI no reliable way to get test coverage here; this never happens
		return false, fmt.Errorf(errorRefsFetchGithubError, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// FYI no reliable way to get test coverage here; this never happens
		return false, fmt.Errorf(errorRefsFetchGithubParseError, err)
	}

	refsString := string(data)
	refExists := strings.Contains(refsString, ref)

	return refExists, nil
}
