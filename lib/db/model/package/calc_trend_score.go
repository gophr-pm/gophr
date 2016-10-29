package pkg

// CalcTrendScore calculates the trend score of a package given its daily,
// weekly and monthly download totals.
func CalcTrendScore(
	dailyDownloads int,
	weeklyDownloads int,
	monthlyDownloads int,
) float32 {
	// Divisors cannot be zero.
	if dailyDownloads < 1 {
		dailyDownloads = 1
	}
	if weeklyDownloads < 1 {
		weeklyDownloads = 1
	}
	if monthlyDownloads < 1 {
		monthlyDownloads = 1
	}

	weeklyDifferential := float32(dailyDownloads*7.0) / float32(weeklyDownloads)
	monthlyDifferential := float32(weeklyDownloads*4.0) /
		float32(monthlyDownloads)

	// How fast the downloads are growing on a day-to-day scale is twice is
	// important as how fast as its growing on a week-to-week scale.
	return (2 * weeklyDifferential) * monthlyDifferential
}
