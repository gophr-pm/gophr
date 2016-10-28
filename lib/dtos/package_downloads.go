package dtos

//go:generate ffjson $GOFILE

// PackageDownloads holds all the download counts of a package, separated into
// time splits.
type PackageDownloads struct {
	Daily   int64 `json:"daily"`
	Weekly  int64 `json:"weekly"`
	Monthly int64 `json:"monthly"`
	AllTime int64 `json:"allTime"`
}
