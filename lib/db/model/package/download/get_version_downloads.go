package download

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db"
)

// GetForVersions takes the author and repo of a package along with a pairing of
// SHAs to the versions thereof, and returns a map of versions to all-time
// download totals.
func GetForVersions(
	q db.Queryable,
	author string,
	repo string,
	shaVersions map[string]string,
) (map[string]int, error) {
	var (
		errs             []error
		resultsChan      = make(chan countResult)
		resultsCount     = 0
		resultsTotal     = len(shaVersions)
		versionDownloads = make(map[string]int)
	)

	// Run an all-time downloads query for each version SHA.
	for sha := range shaVersions {
		go countAllTimeDownloads(q, author, repo, sha, resultsChan)
	}

	// Read all of the results, then exit when we run out.
	for result := range resultsChan {
		if result.err != nil {
			errs = append(errs, result.err)
		} else if version, exists := shaVersions[result.sha]; !exists {
			errs = append(
				errs,
				fmt.Errorf(
					`Failed to match result sha "%s" to a version`,
					result.sha))
		} else {
			// Associate the version with the corresponding downloads total.
			versionDownloads[version] = result.count
		}

		if resultsCount++; resultsCount == resultsTotal {
			close(resultsChan)
		}
	}

	// If there were any errors, return them composed together.
	if len(errs) > 0 {
		return nil, concatErrors(
			"Failed to read version download counts from the database.",
			errs)
	}

	return versionDownloads, nil
}
