package common

//go:generate ffjson $GOFILE

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocql/gocql"
)

const (
	TableNamePackages             = "packages"
	IndexNamePackages             = "packages_index"
	ColumnNamePackagesRepo        = "repo"
	ColumnNamePackagesExists      = "exists"
	ColumnNamePackagesAuthor      = "author"
	ColumnNamePackagesVersions    = "versions"
	ColumnNamePackagesGodocURL    = "godoc_url"
	ColumnNamePackagesIndexTime   = "index_time"
	ColumnNamePackagesAwesomeGo   = "awesome_go"
	ColumnNamePackagesSearchBlob  = "search_blob"
	ColumnNamePackagesDescription = "description"
)

const (
	packagesSearchBlobTemplate = "%s %s %s"
)

var (
	cqlQueryFormatFuzzySearchPackages = fmt.Sprintf(
		`SELECT %s,%s,%s FROM %s WHERE expr(%s,'{query:{type:"fuzzy",field:"%s",value:"%s"}}') LIMIT 10`,
		ColumnNamePackagesRepo,
		ColumnNamePackagesAuthor,
		ColumnNamePackagesDescription,
		TableNamePackages,
		IndexNamePackages,
		ColumnNamePackagesSearchBlob,
		"%s",
	)
)

var (
	alphanumericFilterRegex = regexp.MustCompile(`[^\sa-zA-Z0-9\-_]+`)
)

type PackageModel struct {
	Repo        *string    `json:"repo,omitempty"`
	Exists      *bool      `json:"exists,omitempty"`
	Author      *string    `json:"author,omitempty"`
	Versions    []string   `json:"versions,omitempty"`
	GodocURL    *string    `json:"godocURL,omitempty"`
	IndexTime   *time.Time `json:"-"`
	AwesomeGo   *bool      `json:"awesome,omitempty"`
	SearchBlob  *string    `json:"-"`
	Description *string    `json:"description,omitempty"`
}

func NewPackageModelForInsert(
	repo string,
	exists bool,
	author string,
	versions []string,
	godocURL string,
	indexTime time.Time,
	awesomeGo bool,
	description string,
) (*PackageModel, error) {
	if len(repo) < 1 {
		return nil, NewInvalidParameterError("repo", repo)
	}
	if len(author) < 1 {
		return nil, NewInvalidParameterError("author", author)
	}
	if len(godocURL) < 1 {
		return nil, NewInvalidParameterError("godocURL", godocURL)
	}

	searchBlob := fmt.Sprintf(
		packagesSearchBlobTemplate,
		author,
		repo,
		description,
	)

	return &PackageModel{
		Repo:        &repo,
		Exists:      &exists,
		Author:      &author,
		Versions:    versions,
		GodocURL:    &godocURL,
		IndexTime:   &indexTime,
		AwesomeGo:   &awesomeGo,
		SearchBlob:  &searchBlob,
		Description: &description,
	}, nil
}

func NewPackageModelFromBulkSelect(
	repo string,
	author string,
	description string,
) (*PackageModel, error) {
	if len(repo) < 1 {
		return nil, NewInvalidParameterError("repo", repo)
	}
	if len(author) < 1 {
		return nil, NewInvalidParameterError("author", author)
	}

	return &PackageModel{
		Repo:        &repo,
		Author:      &author,
		Description: &description,
	}, nil
}

func NewPackageModelFromSingleSelect(
	repo string,
	exists bool,
	author string,
	versions []string,
	godocURL string,
	awesomeGo bool,
	description string,
) (*PackageModel, error) {
	if len(repo) < 1 {
		return nil, NewInvalidParameterError("repo", repo)
	}
	if len(author) < 1 {
		return nil, NewInvalidParameterError("author", author)
	}
	if len(godocURL) < 1 {
		return nil, NewInvalidParameterError("godocURL", godocURL)
	}

	return &PackageModel{
		Repo:        &repo,
		Exists:      &exists,
		Author:      &author,
		Versions:    versions,
		GodocURL:    &godocURL,
		AwesomeGo:   &awesomeGo,
		Description: &description,
	}, nil
}

// FuzzySearchPackages fuzzy searches packages by author, package & description.
func FuzzySearchPackages(
	session *gocql.Session,
	searchText string,
) ([]*PackageModel, error) {
	// First, remove all non-essential characters
	searchText = alphanumericFilterRegex.ReplaceAllString(searchText, "")
	// Next put the search text into a query string
	query := fmt.Sprintf(cqlQueryFormatFuzzySearchPackages, searchText)
	// Return the processed results of the query
	return scanPackageModels(session.Query(query))
}

/********************************** HELPERS ***********************************/

// TODO(skeswa): implement this for querying single packages
func scanPackageModel(query *gocql.Query) ([]*PackageModel, error) {
	return nil, nil
}

func scanPackageModels(query *gocql.Query) ([]*PackageModel, error) {
	var (
		err          error
		scanError    error
		closeError   error
		packageModel *PackageModel

		repo        string
		author      string
		description string

		iter          = query.Iter()
		packageModels = make([]*PackageModel, 0)
	)

	for iter.Scan(&repo, &author, &description) {
		packageModel, err = NewPackageModelFromBulkSelect(repo, author, description)
		if err != nil {
			scanError = err
			break
		} else {
			packageModels = append(packageModels, packageModel)
		}
	}

	if err = iter.Close(); err != nil {
		closeError = err
	}

	if scanError != nil || closeError != nil {
		return nil, NewQueryScanError(scanError, closeError)
	}

	return packageModels, nil
}
