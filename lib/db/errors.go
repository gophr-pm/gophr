package db

import "github.com/gocql/gocql"

// IsErrNotFound returns true err resulted from a query that had no results.
func IsErrNotFound(err error) bool {
	return err != nil && err.Error() == gocql.ErrNotFound.Error()
}
