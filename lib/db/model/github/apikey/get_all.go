package apikey

import (
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// GetAll reads every Github API Key matching the forScheduledJobs parameter.
func GetAll(q db.Queryable, forScheduledJobs bool) ([]string, error) {
	var (
		key         apiKey
		keyStrings  []string
		keyIterator = query.Select(
			columnNameKey,
			columnNameForScheduledJobs).
			From(tableName).
			Create(q).
			Iter()
	)

	// Put all of the iterator results into a channel.
	for keyIterator.Scan(
		&key.key,
		&key.forScheduledJobs,
	) {
		// Filter out keys according to the forScheduledJobs parameter.
		if forScheduledJobs == key.forScheduledJobs {
			keyStrings = append(keyStrings, key.key)
		}
	}

	// Attempt to close the iterator.
	if err := keyIterator.Close(); err != nil {
		return nil, err
	}

	return keyStrings, nil
}
