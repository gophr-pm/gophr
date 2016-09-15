package depot

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const (
	refsFetchURLTemplate          = "http://%s/%s.git/info/refs?service=git-upload-pack"
	errorRefsFetchNetworkFailure  = "Could not reach depot at the moment; Please try again later"
	errorRefsFetchNoSuchRepo      = "Could not find a depot repository at %s"
	errorRefsFetchDepotError      = "Depot responded with an error: %v"
	errorRefsFetchDepotParseError = "Cannot read refs from depot: %v"
)

// CheckIfRefExists checks whether a given ref exists in the remote refs list.
func CheckIfRefExists(author, repo string, ref string) (bool, error) {
	repoName := BuildHashedRepoName(author, repo, ref)

	httpClient := &http.Client{Timeout: 10 * time.Second}
	res, err := httpClient.Get(fmt.Sprintf(refsFetchURLTemplate, DepotInternalServiceAddress, repoName))
	if err != nil {
		return false, errors.New(errorRefsFetchNetworkFailure)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return false, nil
	} else if res.StatusCode >= 500 {
		// FYI no reliable way to get test coverage here; this never happens
		return false, fmt.Errorf(errorRefsFetchDepotError, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// FYI no reliable way to get test coverage here; this never happens
		return false, fmt.Errorf(errorRefsFetchDepotParseError, err)
	}

	refsString := string(data)
	refExists := strings.Contains(refsString, ref)

	return refExists, nil
}
