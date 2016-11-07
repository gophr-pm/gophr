package download

import (
	"time"

	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

const (
	// 32 days is long enough for an hourly download count to be useful.
	hourlyDownloadLifespan = (time.Hour * 24) * 32
)

// DeleteOld deletes all hourly download counts older than the hourly download
// count lifespan.
func DeleteOld(q db.Queryable, author string, repo string) error {
	// Delete everything older than one "lifespan" ago.
	downloadAgeBoundary := time.Now().Add(-1 * hourlyDownloadLifespan)

	return query.
		DeleteRows().
		From(hourlyTableName).
		Where(query.Column(hourlyColumnNameAuthor).Equals(author)).
		And(query.Column(hourlyColumnNameRepo).Equals(repo)).
		And(query.Column(hourlyColumnNameHour).
			IsLessThanOrEqualTo(downloadAgeBoundary)).
		Create(q).
		Exec()
}
