package github

import (
	"errors"
	"fmt"
	"net/http"
)

const (
	etagHeader = "Etag"
	shaLength  = 40
)

// FetchFullSHAFromPartialSHA is responsible for fetching a full commit SHA from a short SHA
func FetchFullSHAFromPartialSHA(author, repo, shortSHA string) (string, error) {
	client := &http.Client{}
	archiveURL := fmt.Sprintf(
		"https://github.com/%s/%s/archive/%s.zip",
		author,
		repo,
		shortSHA,
	)
	req, err := http.NewRequest(
		"HEAD",
		archiveURL,
		nil,
	)
	if err != nil {
		return "", err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode == 404 {
		return "", err
	}

	eTagHeader := resp.Header.Get(etagHeader)
	if len(eTagHeader) < shaLength {
		return "", errors.New("Unable to retrieve full commit SHA, Etag header was incomplete or empty.")
	}

	fullSHA := eTagHeader[1 : len(eTagHeader)-1]

	return fullSHA, nil
}
