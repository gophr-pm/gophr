package pkg

import "time"

// Summary summarizes a package. Usually used in plural situations.
type Summary struct {
	Repo             string
	Stars            int
	Author           string
	Awesome          bool
	Description      string
	DailyDownloads   int64
	WeeklyDownloads  int64
	MonthlyDownloads int64
	AllTimeDownloads int64
}

// Details holds package details. Usually used in singular situations.
type Details struct {
	Summary

	TrendScore              float64
	DateLastIndexed         time.Time
	AllTimeVersionDownloads map[string]int64
}
