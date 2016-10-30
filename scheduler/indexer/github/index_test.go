package githubIndexer

import (
	"errors"
	"math/rand"
	"testing"
	"time"

	"github.com/gocql/gocql"
	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/model"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestGithubIndex(t *testing.T) {
	Convey("The github indexer should run", t, func() {
		Convey("if PackageRetriever fails, we should return an error", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, *gocql.Session) {
					return &config.Config{}, &gocql.Session{}
				},
				PackageRetriever: func(session *gocql.Session) ([]*models.PackageModel, error) {
					return nil, errors.New("Failed to query package models")
				},
			})

			So(err, ShouldNotBeNil)
		})
		// TODO finish
		Convey("if PackageRetriever returns no models, we should return an error", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, *gocql.Session) {
					return &config.Config{}, &gocql.Session{}
				},
				PackageRetriever: func(session *gocql.Session) ([]*models.PackageModel, error) {
					var pkgs []*models.PackageModel
					return pkgs, nil
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("if PackageRetriever succeeds, we should return nil", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, *gocql.Session) {
					return &config.Config{}, &gocql.Session{}
				},
				PackageDeleter: func(session *gocql.Session, packageModel *models.PackageModel) error {
					return nil
				},
				PackageInserter: func(session *gocql.Session, packageModel *models.PackageModel) error {
					return nil
				},
				PackageRetriever: func(session *gocql.Session) ([]*models.PackageModel, error) {
					return generateRandomPackageModels(10), nil
				},
				RequestTimeBuffer: 0,
				NewGithubRequestService: func(args github.RequestServiceArgs) github.RequestService {
					m := github.NewMockRequestService()
					m.On("FetchGitHubDataForPackageModel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dtos.GithubRepoDTO{}, errors.New("this is an error"))
					return m
				},
			})

			So(err, ShouldBeNil)
		})
	})
}

var letterRunes = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ")

// randStringRunes generates a random string n runes long.
func randStringRunes(n int) string {
	b := make([]rune, n)
	for i := range b {
		b[i] = letterRunes[rand.Intn(len(letterRunes))]
	}
	return string(b)
}

func generateRandomPackageModels(numPackages int) []*models.PackageModel {
	var pkgModels []*models.PackageModel
	for i := 0; i < numPackages; i++ {
		repo := randStringRunes(5)
		stars := 100
		exists := true
		author := randStringRunes(8)
		versions := []string{randStringRunes(8), randStringRunes(3)}
		godocURL := randStringRunes(16)
		indexTime := time.Now()
		awesomeGo := false
		searchBlob := randStringRunes(20)
		description := randStringRunes(40)
		pkgModel := models.PackageModel{
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
		}
		pkgModels = append(pkgModels, &pkgModel)
	}
	return pkgModels
}
