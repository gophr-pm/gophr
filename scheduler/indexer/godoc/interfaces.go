package godoc

import "github.com/PuerkitoBio/goquery"

// PackageMetadata lol
type PackageMetadata struct {
	githubURL   string
	description string
	author      string
	repo        string
}

// FetchPackageMetadataArgs lol
type FetchPackageMetadataArgs struct {
	ParseHTML htmlParser
}

type htmlParser func(url string) (*goquery.Document, error)
