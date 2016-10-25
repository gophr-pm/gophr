package dtos

//go:generate ffjson $GOFILE

// PackageVersion is a tuple between a version name and the all-time downloads
// count for that version.
type PackageVersion struct {
	Name             string `json:"name"`
	AllTimeDownloads int64  `json:"allTimeDownloads"`
}
