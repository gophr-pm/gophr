package pkg

import (
	"fmt"
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// InsertArgs is the arguments struct for Insert.
type InsertArgs struct {
	Repo        string
	Stars       int
	Author      string
	Awesome     bool
	Queryable   db.Queryable
	Description string
}

// Insert puts a package into the database if it doesn't exist.
func Insert(args InsertArgs) error {
	// Now that we have all the requisite data, insert the new package.
	if err := query.InsertInto(packagesTableName).
		Value(packagesColumnNameRepo, args.Repo).
		Value(packagesColumnNameStars, args.Stars).
		Value(packagesColumnNameAuthor, args.Author).
		Value(packagesColumnNameAwesome, args.Awesome).
		Value(packagesColumnNameTrendScore, 0).
		Value(
			packagesColumnNameSearchScore,
			CalcSearchScore(args.Stars, 0, args.Awesome, 0)).
		Value(packagesColumnNameDescription, args.Description).
		Value(packagesColumnNameDateDiscovered, time.Now()).
		Value(
			packagesColumnNameSearchBlob,
			composeSearchBlob(args.Author, args.Repo, args.Description)).
		IfNotExists().
		Create(args.Queryable).
		Exec(); err != nil {
		return fmt.Errorf(
			"Failed to insert package %s/%s: %v",
			args.Author,
			args.Repo,
			err)
	}

	return nil
}
