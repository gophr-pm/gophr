package godoc

import "github.com/PuerkitoBio/goquery"

// PackageMetadata lol
type PackageMetadata struct {
	githubURL   string
	description string
	author      string
	repo        string
}

// FetchPackageMetadataArgs is the args struct for fetching package metadata
// from godoc.
type FetchPackageMetadataArgs struct {
	ParseHTML htmlParser
}

// htmlParser parses an HTML doc from a url and returns a goquery doc.
type htmlParser func(url string) (*goquery.Document, error)
