package dtos

import "time"

//go:generate ffjson $GOFILE

// PackageDetails is the DTO for most singular package requests.
type PackageDetails struct {
	Repo            string           `json:"repo"`
	Stars           int              `json:"stars"`
	Author          string           `json:"author"`
	Awesome         bool             `json:"awesome"`
	Versions        []PackageVersion `json:"versions"`
	Downloads       PackageDownloads `json:"downloads"`
	TrendScore      float32          `json:"trendScore"`
	Description     string           `json:"description"`
	DateDiscovered  time.Time        `json:"dateDiscovered"`
	DateLastIndexed time.Time        `json:"dateLastIndexed"`
}
