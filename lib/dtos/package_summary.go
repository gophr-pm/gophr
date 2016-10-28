package dtos

//go:generate ffjson $GOFILE

// PackageSummary is the DTO for most plural package requests.
type PackageSummary struct {
	Repo        string           `json:"repo"`
	Stars       int              `json:"stars"`
	Author      string           `json:"author"`
	Awesome     bool             `json:"awesome"`
	Downloads   PackageDownloads `json:"downloads"`
	Description string           `json:"description"`
}
