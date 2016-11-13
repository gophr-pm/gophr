package github

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	etagHeader             = "Etag"
	minSHALengthWithQuotes = 42
	baseGithubArchiveURL   = "https://github.com/%s/%s/archive/%s.zip"
)

// HTTPHeadReq executes an HTTP `HEAD` to the specified URL and returns the
// corresponding response.
type HTTPHeadReq func(url string) (*http.Header, error)

// ExpandPartialSHAArgs is the arguments struct for ExpandPartialSHA.
type ExpandPartialSHAArgs struct {
	Author     string
	Repo       string
	ShortSHA   string
	DoHTTPHead HTTPHeadReq
}

// ExpandPartialSHA is responsible for fetching a full commit SHA from a short
// SHA. This works by sending a HEAD request to the git archive endpoint with a
// short SHA. The request returns a full SHA of the archive in the `Etag`
// of the request header that is sent back.
func (svc *requestServiceImpl) ExpandPartialSHA(
	args ExpandPartialSHAArgs,
) (string, error) {
	archiveURL := fmt.Sprintf(
		baseGithubArchiveURL,
		args.Author,
		args.Repo,
		args.ShortSHA,
	)

	gitHubRespHeader, err := args.DoHTTPHead(archiveURL)
	if err != nil {
		return "", err
	}

	eTagHeader := gitHubRespHeader.Get(etagHeader)
	if len(eTagHeader) != minSHALengthWithQuotes {
		return "", errors.New(
			"Unable to retrieve full commit SHA, " +
				"Etag header was incomplete or empty.")
	}

	// TODO(skeswa): the quotes can be a messy assumption. Should probably just be
	// using strings.Trim(). @Shikkic.
	// The Etag in the header contains the full SHA wrapped in quotes.
	// We need to remove the quotes.
	fullSHA := eTagHeader[1 : len(eTagHeader)-1]

	return fullSHA, nil
}
