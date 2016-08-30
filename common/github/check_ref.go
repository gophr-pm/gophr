package github

// CheckIfRefExists downloads and processes refs data from Github and checks
import (
	"errors"
	"fmt"
	"io/ioutil"
	"strings"
)

// whether a given ref exists in the remote refs list.
func CheckIfRefExists(author, repo string, ref string) (bool, error) {
	ref = BuildGitHubBranch(ref)
	repo = BuildNewGitHubRepoName(author, repo)
	author = GitHubGophrPackageOrgName
	githubRoot := fmt.Sprintf(
		githubRootTemplate,
		author,
		repo,
	)

	res, err := httpClient.Get(fmt.Sprintf(refsFetchURLTemplate, githubRoot))
	if err != nil {
		return false, errors.New(errorRefsFetchNetworkFailure)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return false, fmt.Errorf(errorRefsFetchNoSuchRepo, githubRoot)
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
