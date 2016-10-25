package dtos

import "time"

//go:generate ffjson $GOFILE

// PackageDetails is the DTO for most singular package requests.
type PackageDetails struct {
	Repo            string           `json:"repo"`
	Stars           int              `json:"stars"`
	Author          string           `json:"author"`
	Awesome         bool             `json:"awesome"`
	Versions        []PackageVersion `json:"version"`
	Downloads       PackageDownloads `json:"downloads"`
	TrendScore      float64          `json:"trendScore"`
	Description     string           `json:"description"`
	DateLastIndexed time.Time        `json:"dateLastIndexed"`
}
