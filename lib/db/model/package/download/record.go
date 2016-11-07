package download

import (
	"bytes"
	"errors"
	"strconv"
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/github"
)

// Record records a single download of specific package version.
func Record(
	q db.BatchingQueryable,
	author string,
	repo string,
	sha string,
	ghSvc github.RequestService,
) error {
	var (
		// Normalize the hourly date by setting all lesser time fields to zero.
		now      = time.Now()
		thisHour = time.Date(
			now.Year(),
			now.Month(),
			now.Day(),
			now.Hour(),
			0, // minutes
			0, // seconds
			0, // nanoseconds
			time.Local)
		resultsChan = make(chan error)
	)

	// Execute the first update query. Exit if it fails.
	go bumpDownloads(q, thisHour, author, repo, sha, resultsChan)
	go assertPackageExistence(q, author, repo, ghSvc, resultsChan)

	var (
		i    = 0
		errs []error
	)

	// Read all (2) of the responses.
	for err := range resultsChan {
		// Record all of the errors.
		if err != nil {
			errs = append(errs, err)
		}

		i++
		// Exit if every db call has returned.
		if i >= 2 {
			close(resultsChan)
		}
	}

	if len(errs) > 0 {
		// TODO(skeswa): this style of error composition is also done in verdeps.
		// This should be extracted to a helper function.
		var buffer bytes.Buffer
		buffer.WriteString("Failed to record a daily download due to ")
		buffer.WriteString(strconv.Itoa(len(errs)))
		buffer.WriteString(" error(s): [ ")
		for i, err := range errs {
			if i > 0 {
				buffer.WriteString(", ")
			}

			buffer.WriteString(err.Error())
		}
		buffer.WriteString(" ]")

		return errors.New(buffer.String())
	}

	return nil
}
