package pkg

import (
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/db/query"
)

// ReadAll reads every package in the database into a channel.
func ReadAll(q db.Queryable, summaries chan Summary, errs chan error) {
	var (
		summary         Summary
		summaryIterator = query.Select(
			packagesColumnNameRepo,
			packagesColumnNameStars,
			packagesColumnNameAuthor,
			packagesColumnNameAwesome,
			packagesColumnNameDescription,
			packagesColumnNameDailyDownloads,
			packagesColumnNameWeeklyDownloads,
			packagesColumnNameMonthlyDownloads,
			packagesColumnNameAllTimeDownloads).
			From(packagesTableName).
			Create(q).
			Iter()
	)

	// Signal that there are no more summaries coming through the pipe on exit.
	defer close(summaries)

	// Put all of the iterator results into a channel.
	for summaryIterator.Scan(
		&summary.Repo,
		&summary.Stars,
		&summary.Author,
		&summary.Awesome,
		&summary.Description,
		&summary.DailyDownloads,
		&summary.WeeklyDownloads,
		&summary.MonthlyDownloads,
		&summary.AllTimeDownloads) {
		summaries <- summary
	}

	// Attempt to close the iterator.
	if err := summaryIterator.Close(); err != nil {
		errs <- err
	}
}
