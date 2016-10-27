package pkg

import (
	"errors"
	"fmt"

	"github.com/gophr-pm/gophr/lib/db/query"
)

var (
	// descSortExprTemplate, but scoped for the date discovered column.
	descSortByDateDiscExpr = fmt.Sprintf(
		descSortExprTemplate,
		packagesColumnNameDateDiscovered)
)

// GetNew gets up to "limit" of the most recently discovered packages.
func GetNew(q query.Queryable, limit int) (Summaries, error) {
	if limit < 1 {
		return nil, errors.New("Limit must be greater than zero")
	}

	// Create and execute the query, then create in iterator for the results.
	iter := query.
		Select(
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
		Where(query.Index(packagesIndexName).Matches(descSortByDateDiscExpr)).
		Limit(limit).
		Create(q).
		Iter()

	var (
		summaries   []Summary
		nextSummary Summary
	)

	// Scan into a summary struct. Add it to the list if successful.
	for iter.Scan(
		&nextSummary.Repo,
		&nextSummary.Stars,
		&nextSummary.Author,
		&nextSummary.Awesome,
		&nextSummary.Description,
		&nextSummary.DailyDownloads,
		&nextSummary.WeeklyDownloads,
		&nextSummary.MonthlyDownloads,
		&nextSummary.AllTimeDownloads) {
		summaries = append(summaries, nextSummary)
	}

	if err := iter.Close(); err != nil {
		return nil, fmt.Errorf(
			`Failed to get new packages from the db: %v`,
			err)
	}

	return summaries, nil
}
