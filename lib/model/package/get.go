package pkg

import (
	"fmt"

	"github.com/gophr-pm/gophr/lib/db/query"
)

// Get fetches a single package, matching the author and repo parameters, from
// the the database.
func Get(q query.Queryable, author, repo string) (Details, error) {
	var result Details

	// Create and execute the query, then create in iterator for the results.
	if err := query.
		Select(
			packagesColumnNameRepo,
			packagesColumnNameStars,
			packagesColumnNameAuthor,
			packagesColumnNameAwesome,
			packagesColumnNameTrendScore,
			packagesColumnNameDescription,
			packagesColumnNameDailyDownloads,
			packagesColumnNameWeeklyDownloads,
			packagesColumnNameDateLastIndexed,
			packagesColumnNameMonthlyDownloads,
			packagesColumnNameAllTimeDownloads,
			packagesColumnNameAllTimeVersionDownloads).
		From(packagesTableName).
		Where(query.Column(packagesColumnNameAuthor).Equals(author)).
		And(query.Column(packagesColumnNameRepo).Equals(repo)).
		Limit(1).
		Create(q).
		Scan(
			result.Repo,
			result.Stars,
			result.Author,
			result.Awesome,
			result.TrendScore,
			result.Description,
			result.DailyDownloads,
			result.WeeklyDownloads,
			result.DateLastIndexed,
			result.MonthlyDownloads,
			result.AllTimeDownloads,
			result.AllTimeVersionDownloads); err != nil {
		return result, fmt.Errorf(
			`Failed to get package %s/%s from the db: %v`,
			author,
			repo,
			err)
	}

	return result, nil
}
