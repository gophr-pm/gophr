package main

import (
	"bytes"
	"fmt"
	"net/http"
)

const internalDepotAPIDomain = "depot-int-svc"

// packageExistsInDepot will return true if a package matching author, repo and
// sha exists in depot.
func packageExistsInDepot(author, repo, sha string) (bool, error) {
	req, err := http.NewRequest(
		"GET",
		fmt.Sprintf(
			"http://%s/api/repos/%s/%s/%s",
			internalDepotAPIDomain,
			author,
			repo,
			sha),
		nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotFound:
		return false, nil
	default:
		errBuffer := bytes.Buffer{}
		errBuffer.ReadFrom(res.Body)
		return false, fmt.Errorf("Could not create repo in depot: %s.", errBuffer.String())
	}
}

// createRepoInDepot creates a package repo in depot matching the author, repo &
// sha specified. Returns true if the repo was created by this func., or returns
// false is the the directory already existed.
func createRepoInDepot(author, repo, sha string) (bool, error) {
	req, err := http.NewRequest(
		"POST",
		fmt.Sprintf(
			"http://%s/api/repos/%s/%s/%s",
			internalDepotAPIDomain,
			author,
			repo,
			sha),
		nil)
	if err != nil {
		return false, err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return false, err
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return true, nil
	case http.StatusNotModified:
		return false, nil
	default:
		errBuffer := bytes.Buffer{}
		errBuffer.ReadFrom(res.Body)
		return false, fmt.Errorf("Could not create repo in depot: %s.", errBuffer.String())
	}
}

// deleteRepoInDepot deletes a package repo in depot matching the author, repo &
// sha specified.
func deleteRepoInDepot(author, repo, sha string) error {
	req, err := http.NewRequest(
		"DELETE",
		fmt.Sprintf(
			"http://%s/api/repos/%s/%s/%s",
			internalDepotAPIDomain,
			author,
			repo,
			sha),
		nil)
	if err != nil {
		return err
	}

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	switch res.StatusCode {
	case http.StatusOK:
		return nil
	default:
		errBuffer := bytes.Buffer{}
		errBuffer.ReadFrom(res.Body)
		return fmt.Errorf("Could not create repo in depot: %s.", errBuffer.String())
	}
}
