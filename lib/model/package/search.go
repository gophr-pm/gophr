package pkg

import (
	"errors"
	"fmt"

	"github.com/gophr-pm/gophr/lib/db/query"
)

var (
	// Fill in all of the column name slots.
	refinedSearchExprTemplate = fmt.Sprintf(
		searchExprTemplate,
		packagesColumnNameSearchBlob,
		"%s",
		packagesColumnNameStars,
		packagesColumnNameAwesome,
		packagesColumnNameAllTimeDownloads)
)

// Search (as in "get top ten") gets the top packages sorted descendingly
// within the specified time split.
func Search(
	q query.Queryable,
	searchQuery string,
	limit int,
) (Summaries, error) {
	if len(searchQuery) < 1 {
		return nil, errors.New("Search query cannot be blank")
	}
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
		Where(query.Index(packagesIndexName).Matches(fmt.Sprintf(
			refinedSearchExprTemplate,
			searchQuery))).
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
		return nil, fmt.Errorf(`Failed to search for packages in the db: %v`, err)
	}

	return summaries, nil
}
