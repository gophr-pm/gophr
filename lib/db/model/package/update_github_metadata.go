package pkg

import (
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// UpdateGithubMetadataArgs is the arguments struct for UpdateGithubMetadata.
type UpdateGithubMetadataArgs struct {
	Repo        string
	Stars       int
	Author      string
	Queryable   db.Queryable
	Description string
}

// UpdateGithubMetadata updates the Github metadata for a package.
func UpdateGithubMetadata(args UpdateGithubMetadataArgs) error {
	return query.
		Update(packagesTableName).
		Set(packagesColumnNameStars, args.Stars).
		Set(packagesColumnNameDescription, args.Description).
		Where(query.Column(packagesColumnNameRepo).Equals(args.Repo)).
		And(query.Column(packagesColumnNameAuthor).Equals(args.Author)).
		IfExists().
		Create(args.Queryable).
		Exec()
}
