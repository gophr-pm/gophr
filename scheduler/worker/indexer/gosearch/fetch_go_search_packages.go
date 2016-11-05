package gosearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
)

// httpGetter executes an HTTP get to the specified URL and returns the
// corresponding response.
type httpGetter func(url string) (*http.Response, error)

const (
	// go-search.org API endpoint that returns a list of every package that it has
	// ever seen.
	goSearchAPIPackagesEndpoint = "http://go-search.org/api?action=packages"
	// noPackageLimit is the value of the limit parameter of fetchGoSearchPackages
	// that indicates that there is no limit.
	noPackageLimit = -1
	// devPackageLimit is the value of the limit parameter of fetchGoSearchPackages
	// when in the dev environment.
	devPackageLimit = 250
)

var (
	// githubImportPathRegex matches the import paths from go-search that
	// gophr supports. It has two capture groups: (1) author, (2) repo.
	githubImportPathRegex = regexp.MustCompile(`^github\.com/([^/]+)/([^/]+)`)
)

// fetchGoSearchPackages fetches the set of all packages of which go-search.org
// is aware.
func fetchGoSearchPackages(
	goHTTPGet httpGetter,
	limit int,
) (*packageSet, error) {
	// Hit the packages endpoint to get that sweet, sweet package data.
	resp, err := goHTTPGet(goSearchAPIPackagesEndpoint)
	if err != nil {
		return nil, fmt.Errorf("Failed to make request to go-search: %v.", err)
	}

	// Read everything (about 19MB's worth 😲) into memory.
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Failed to read go-search response: %v.", err)
	}

	// Make the json a little easier to work with.
	var packageImportPaths []string
	if err = json.Unmarshal(body, &packageImportPaths); err != nil {
		return nil, fmt.Errorf("Failed to unmarshal go-search response: %v.", err)
	}

	// Add every author-repo combination in the import paths into the package set.
	packages := newPackageSet()
	for _, packageImportPath := range packageImportPaths {
		// Check if we've hit the limit (if one such limit exists).
		if limit != noPackageLimit && packages.len() >= limit {
			break
		}

		// Get the important bits out of the import path and throw them into the
		// set. Only continue, however, if the import path is one supported by
		// gophr.
		if parts := githubImportPathRegex.FindStringSubmatch(
			packageImportPath,
		); parts != nil {
			packages.add(parts[1], parts[2])
		}
	}

	return packages, nil
}
