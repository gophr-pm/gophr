package pkg

const (
	trendScoreWeightTotal = weeklyDifferentialWeight +
		monthlyDifferentialWeight
	// TODO(skeswa): look into comparing against the median dynamically.
	// A "reasonable" weekly download total.
	reasonableWeeklyDownloads  = float32(200)
	reasonableMonthlyDownloads = reasonableWeeklyDownloads * 4
	// Weekly growth rate impacts "trending-ness" more than its monthly
	// counterpart.
	weeklyDifferentialWeight  = 2
	monthlyDifferentialWeight = 1
)

// CalcTrendScore calculates the trend score of a package given its daily,
// weekly and monthly download totals.
func CalcTrendScore(
	dailyDownloads int,
	weeklyDownloads int,
	monthlyDownloads int,
) float32 {
	var (
		weeklyDifferential  float32
		monthlyDifferential float32
	)

	// Compare daily growth to weekly growth.
	if weeklyDownloads > 0 {
		weeklyDifferential = float32(dailyDownloads*7) / float32(weeklyDownloads)
	} else {
		// If there is no weekly, compare against a reasonable weekly downloads
		// count.
		weeklyDifferential = float32(dailyDownloads*7) / reasonableWeeklyDownloads
	}

	// Compare weekly growth to monthly growth.
	if monthlyDownloads > 0 {
		monthlyDifferential = float32(weeklyDownloads*4) /
			float32(monthlyDownloads)
	} else {
		// If there is no monthly, compare against a reasonable monthly downloads
		// count.
		weeklyDifferential = float32(weeklyDownloads*4) / reasonableMonthlyDownloads
	}

	// How fast the downloads are growing on a day-to-day scale is twice is
	// important as how fast as its growing on a week-to-week scale.
	return ((weeklyDifferentialWeight * weeklyDifferential) +
		(monthlyDifferentialWeight * monthlyDifferential)) / trendScoreWeightTotal
}
