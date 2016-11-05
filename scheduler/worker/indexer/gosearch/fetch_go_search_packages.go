package gosearch

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gophr-pm/gophr/lib/verdeps"
)

// httpGetter executes an HTTP get to the specified URL and returns the
// corresponding response.
type httpGetter func(url string) (*http.Response, error)

// go-search.org API endpoint that returns a list of every package that it has
// ever seen.
const goSearchAPIPackagesEndpoint = "http://go-search.org/api?action=packages"

// noPackageLimit is the value of the limit parameter of fetchGoSearchPackages
// that indicates that there is no limit.
const noPackageLimit = -1

// devPackageLimit is the value of the limit parameter of fetchGoSearchPackages
// when in the dev environment.
const devPackageLimit = 1000

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
		// set.
		author, repo, _ := verdeps.ParseImportPath(packageImportPath)
		packages.add(author, repo)
	}

	return packages, nil
}
