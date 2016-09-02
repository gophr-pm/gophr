package models

import (
	"fmt"
	"regexp"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common/errors"
)

// Constants directly related to interacting with the package model in the
// cassandra database.
const (
	// TableNamePackages is the name of the table containing the package model.
	TableNamePackages = "gophr.packages"
	// IndexNamePackages is the name of the lucene index
	IndexNamePackages             = "packages_index"
	ColumnNamePackagesRepo        = "repo"
	ColumnNamePackagesStars       = "stars"
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
	cqlQueryFuzzySearchPackagesTemplate = fmt.Sprintf(
		`SELECT %s,%s,%s FROM %s WHERE expr(%s,'{query:{type:"fuzzy",field:"%s",value:"%s"}}') LIMIT 10`,
		ColumnNamePackagesRepo,
		ColumnNamePackagesAuthor,
		ColumnNamePackagesDescription,
		TableNamePackages,
		IndexNamePackages,
		ColumnNamePackagesSearchBlob,
		"%s",
	)

	cqlQuerySelectPackageVersions = fmt.Sprintf(
		`SELECT %s FROM %s WHERE %s = ? AND %s = ? LIMIT 1`,
		ColumnNamePackagesVersions,
		TableNamePackages,
		ColumnNamePackagesAuthor,
		ColumnNamePackagesRepo,
	)

	cqlQueryInsertPackage = fmt.Sprintf(
		`INSERT INTO %s (%s,%s,%s,%s,%s,%s,%s,%s,%s,%s) VALUES (?,?,?,?,?,?,?,?,?,?)`,
		TableNamePackages,
		ColumnNamePackagesRepo,
		ColumnNamePackagesStars,
		ColumnNamePackagesExists,
		ColumnNamePackagesAuthor,
		ColumnNamePackagesVersions,
		ColumnNamePackagesGodocURL,
		ColumnNamePackagesIndexTime,
		ColumnNamePackagesAwesomeGo,
		ColumnNamePackagesSearchBlob,
		ColumnNamePackagesDescription,
	)

	cqlQueryDeletePackage = fmt.Sprintf(
		`DELETE FROM %s WHERE %s = ? AND %s = ?`,
		TableNamePackages,
		ColumnNamePackagesAuthor,
		ColumnNamePackagesRepo,
	)
)

var (
	alphanumericFilterRegex = regexp.MustCompile(`[^\sa-zA-Z0-9\-_]+`)
)

// PackageModel is a struct representing one individual package in the database.
type PackageModel struct {
	Repo        *string
	Stars       *int
	Exists      *bool
	Author      *string
	Versions    []string
	GodocURL    *string
	IndexTime   *time.Time
	AwesomeGo   *bool
	SearchBlob  *string
	Description *string
}

// NewPackageModelForInsert creates an instance of PackageModel that is
// optimized and validated for the insert operation in the database.
func NewPackageModelForInsert(
	author string,
	exists bool,
	repo string,
	versions []string,
	godocURL string,
	indexTime time.Time,
	awesomeGo bool,
	description string,
	stars int,
) (*PackageModel, error) {
	if len(repo) < 1 {
		return nil, errors.NewInvalidParameterError("repo", repo)
	}
	if len(author) < 1 {
		return nil, errors.NewInvalidParameterError("author", author)
	}
	if len(godocURL) < 1 {
		return nil, errors.NewInvalidParameterError("godocURL", godocURL)
	}

	searchBlob := fmt.Sprintf(
		packagesSearchBlobTemplate,
		author,
		repo,
		description,
	)

	return &PackageModel{
		Repo:        &repo,
		Stars:       &stars,
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

// NewPackageModelFromBulkSelect creates an instance of PackageModel that is
// optimized and validated for a select operation designed to get data about
// multiple packages from the database.
func NewPackageModelFromBulkSelect(
	author string,
	repo string,
	description string,
) (*PackageModel, error) {
	if len(repo) < 1 {
		return nil, errors.NewInvalidParameterError("repo", repo)
	}
	if len(author) < 1 {
		return nil, errors.NewInvalidParameterError("author", author)
	}

	return &PackageModel{
		Repo:        &repo,
		Author:      &author,
		Description: &description,
	}, nil
}

// TODO(Shikkic): get your shit together.
func NewPackageModelTest(
	author string,
	repo string,
	awesome_go bool,
	description string,
	exists bool,
	godoc_url string,
	index_time time.Time,
	search_blob string,
	versions []string,
	stars int,
) *PackageModel {
	return &PackageModel{
		Repo:        &repo,
		Stars:       &stars,
		Exists:      &exists,
		Author:      &author,
		Versions:    versions,
		GodocURL:    &godoc_url,
		IndexTime:   &index_time,
		AwesomeGo:   &awesome_go,
		SearchBlob:  &search_blob,
		Description: &description,
	}
}

// NewPackageModelFromSingleSelect creates an instance of PackageModel that is
// optimized and validated for a select operation designed to get data about
// a single package from the database.
func NewPackageModelFromSingleSelect(
	author string,
	exists bool,
	repo string,
	versions []string,
	godocURL string,
	awesomeGo bool,
	description string,
) (*PackageModel, error) {
	if len(repo) < 1 {
		return nil, errors.NewInvalidParameterError("repo", repo)
	}
	if len(author) < 1 {
		return nil, errors.NewInvalidParameterError("author", author)
	}
	if len(godocURL) < 1 {
		return nil, errors.NewInvalidParameterError("godocURL", godocURL)
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

// FindPackageVersions gets the versions of a package from the database. If
// no such package exists, or there were no versions for said package, then nil
// is returned.
func FindPackageVersions(session *gocql.Session, author string, repo string) ([]string, error) {
	var (
		err      error
		versions []string
	)

	iter := session.Query(cqlQuerySelectPackageVersions, author, repo).Iter()

	if !iter.Scan(&versions) {
		return nil, nil
	}

	if err = iter.Close(); err != nil {
		return nil, errors.NewQueryScanError(nil, err)
	}

	return versions, nil
}

// FuzzySearchPackages finds a list of packages relevant to a query phrase
// string. The search takes author, package and description into account.
func FuzzySearchPackages(
	session *gocql.Session,
	searchText string,
) ([]*PackageModel, error) {
	// First, remove all non-essential characters
	searchText = alphanumericFilterRegex.ReplaceAllString(searchText, "")
	// Next put the search text into a query string
	query := fmt.Sprintf(cqlQueryFuzzySearchPackagesTemplate, searchText)
	// Return the processed results of the query
	return scanPackageModels(session.Query(query))
}

// InsertPackage inserts an individual package into the database.
func InsertPackage(
	session *gocql.Session,
	packageModel *PackageModel,
) error {
	err := session.Query(cqlQueryInsertPackage,
		*packageModel.Repo,
		*packageModel.Stars,
		*packageModel.Exists,
		*packageModel.Author,
		packageModel.Versions,
		*packageModel.GodocURL,
		*packageModel.IndexTime,
		*packageModel.AwesomeGo,
		*packageModel.SearchBlob,
		*packageModel.Description,
	).Exec()

	return err
}

// InsertPackages inserts a slice of package models into the database.
func InsertPackages(
	session *gocql.Session,
	packageModels []*PackageModel,
) error {
	batch := gocql.NewBatch(gocql.LoggedBatch)

	if packageModels == nil || len(packageModels) == 0 {
		return errors.NewInvalidParameterError("packageModels", packageModels)
	}

	for _, packageModel := range packageModels {
		if packageModel != nil &&
			packageModel.Repo != nil &&
			packageModel.Exists != nil &&
			packageModel.Author != nil &&
			packageModel.GodocURL != nil &&
			packageModel.IndexTime != nil &&
			packageModel.AwesomeGo != nil &&
			packageModel.SearchBlob != nil &&
			packageModel.Description != nil {
			batch.Query(
				cqlQueryInsertPackage,
				*packageModel.Repo,
				*packageModel.Exists,
				*packageModel.Author,
				packageModel.Versions,
				*packageModel.GodocURL,
				*packageModel.IndexTime,
				*packageModel.AwesomeGo,
				*packageModel.SearchBlob,
				*packageModel.Description,
			)
		} else {
			return errors.NewInvalidParameterError(
				"packageModels",
				fmt.Sprintf("[ ..., %v, ... ]", packageModel),
			)
		}
	}

	err := session.ExecuteBatch(batch)
	if err != nil {
		return err
	}

	return nil
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
		packageModel, err = NewPackageModelFromBulkSelect(author, repo, description)
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
		return nil, errors.NewQueryScanError(scanError, closeError)
	}

	return packageModels, nil
}

func ScanAllPackageModels(session *gocql.Session) ([]*PackageModel, error) {
	var (
		err          error
		scanError    error
		closeError   error
		packageModel *PackageModel

		author      string
		repo        string
		awesome_go  bool
		description string
		exists      bool
		godoc_url   string
		index_time  time.Time
		search_blob string
		versions    []string
		stars       int

		query = session.Query(`SELECT
			author,
			repo,
			awesome_go,
			description,
			exists,
			godoc_url,
			index_time,
			search_blob,
			versions,
			stars
			FROM gophr.packages`)
		iter          = query.Iter()
		packageModels = make([]*PackageModel, 0)
	)

	for iter.Scan(&author, &repo, &awesome_go, &description, &exists, &godoc_url, &index_time, &search_blob, &versions, &stars) {
		packageModel = NewPackageModelTest(author, repo, awesome_go, description, exists, godoc_url, index_time, search_blob, versions, stars)
		packageModels = append(packageModels, packageModel)
	}

	if err = iter.Close(); err != nil {
		closeError = err
	}

	if scanError != nil || closeError != nil {
		return nil, errors.NewQueryScanError(scanError, closeError)
	}

	return packageModels, nil
}

func DeletePackageModel(session *gocql.Session, packageModel *PackageModel) error {
	author := *packageModel.Author
	repo := *packageModel.Repo
	query := session.Query(cqlQueryDeletePackage, author, repo)
	err := query.Exec()

	return err
}
