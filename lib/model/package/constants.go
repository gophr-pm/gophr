package pkg

const (
	packagesTableName                         = "packages"
	packagesIndexName                         = "packages_index"
	packagesColumnNameRepo                    = "repo"
	packagesColumnNameStars                   = "stars"
	packagesColumnNameAuthor                  = "author"
	packagesColumnNameAwesome                 = "awesome"
	packagesColumnNameTrendScore              = "trend_score"
	packagesColumnNameSearchBlob              = "search_blob"
	packagesColumnNameDescription             = "description"
	packagesColumnNameDateDiscovered          = "date_discovered"
	packagesColumnNameDailyDownloads          = "daily_downloads"
	packagesColumnNameWeeklyDownloads         = "weekly_downloads"
	packagesColumnNameDateLastIndexed         = "date_last_indexed"
	packagesColumnNameMonthlyDownloads        = "monthly_downloads"
	packagesColumnNameAllTimeDownloads        = "all_time_downloads"
	packagesColumnNameAllTimeVersionDownloads = "all_time_downloads"

	awesomeTableName        = "awesome_packages"
	awesomeColumnNameRepo   = "repo"
	awesomeColumnNameAuthor = "author"

	descSortExprTemplate = `{sort:{fields:[{ field: "%s", reverse: true }]}}`
	searchExprTemplate   = `{
		query: { type: "fuzzy", field: "%s", value: "%s" },
		sort: {
			fields: [
				{ field: "%s", reverse: true },
				{ field: "%s", reverse: true },
				{ field: "%s", reverse: true }
			]
		}
	}`
)
