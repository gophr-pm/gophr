package pkg

const (
	baseSubScore = 100

	searchScoreWeightTotal = starsSearchScoreWeight +
		awesomeSearchScoreWeight +
		trendScoreSearchScoreWeight +
		allTimeDownloadsSearchScoreWeight
	starsSearchScoreWeight            = 5
	awesomeSearchScoreWeight          = 2
	trendScoreSearchScoreWeight       = 2
	allTimeDownloadsSearchScoreWeight = 3

	approximateMaxStars            = float32(50000)
	approximateMaxTrendScore       = float32(10)
	approximateMaxAllTimeDownloads = float32(10000)
)

// CalcSearchScore calculates the search score of a package given its stars,
// downloads, awesome and trend score.
func CalcSearchScore(
	stars int,
	allTimeDownloads int,
	awesome bool,
	trendScore float32,
) float32 {
	var (
		starsSubScore = baseSubScore *
			(float32(stars) / approximateMaxStars)
		trendSubScore = baseSubScore *
			(trendScore / approximateMaxTrendScore)
		awesomeSubScore          float32
		allTimeDownloadsSubScore = baseSubScore *
			(float32(allTimeDownloads) / approximateMaxAllTimeDownloads)
	)

	if awesome {
		awesomeSubScore = baseSubScore
	}

	return ((starsSubScore * starsSearchScoreWeight) +
		(trendSubScore * trendScoreSearchScoreWeight) +
		(awesomeSubScore * awesomeSearchScoreWeight) +
		(allTimeDownloadsSubScore * allTimeDownloadsSearchScoreWeight)) /
		searchScoreWeightTotal
}
