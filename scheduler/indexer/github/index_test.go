package githubIndexer

import (
	"errors"
	"testing"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
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
				Init: func() (*config.Config, db.Client) {
					s := db.NewMockClient()
					s.On("Close")
					return &config.Config{}, s
				},
				PackageRetriever: func(session db.Client) ([]*models.PackageModel, error) {
					return nil, errors.New("Failed to query package models")
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("if PackageRetriever returns no models, we should return an error", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					s := db.NewMockClient()
					s.On("Close")
					return &config.Config{}, s
				},
				PackageRetriever: func(session db.Client) ([]*models.PackageModel, error) {
					var pkgs []*models.PackageModel
					return pkgs, nil
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("if NewGithubRequestService returns and error, we should return nil", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					s := db.NewMockClient()
					s.On("Close")
					return &config.Config{}, s
				},
				PackageDeleter: func(session db.Client, packageModel *models.PackageModel) error {
					return nil
				},
				PackageInserter: func(session db.Client, packageModel *models.PackageModel) error {
					return nil
				},
				PackageRetriever: func(session db.Client) ([]*models.PackageModel, error) {
					return generateRandomPackageModels(10), nil
				},
				RequestTimeBuffer: 0,
				NewGithubRequestService: func(args github.RequestServiceArgs) github.RequestService {
					m := github.NewMockRequestService()
					m.On("FetchGitHubDataForPackageModel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dtos.GithubRepo{}, errors.New("this is an error"))
					return m
				},
			})

			So(err, ShouldBeNil)
		})

		Convey("if PackageRetriever succeeds, we should return nil", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					s := db.NewMockClient()
					s.On("Close")
					return &config.Config{}, s
				},
				PackageDeleter: func(session db.Client, packageModel *models.PackageModel) error {
					return nil
				},
				PackageInserter: func(session db.Client, packageModel *models.PackageModel) error {
					return nil
				},
				PackageRetriever: func(session db.Client) ([]*models.PackageModel, error) {
					return generateRandomPackageModels(10), nil
				},
				RequestTimeBuffer: 0,
				NewGithubRequestService: func(args github.RequestServiceArgs) github.RequestService {
					m := github.NewMockRequestService()
					m.On("FetchGitHubDataForPackageModel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dtos.GithubRepo{}, nil)
					return m
				},
			})

			So(err, ShouldBeNil)
		})

		Convey("if PackageIndexer fails, we should return nil", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					s := db.NewMockClient()
					s.On("Close")
					return &config.Config{}, s
				},
				PackageDeleter: func(session db.Client, packageModel *models.PackageModel) error {
					return nil
				},
				PackageInserter: func(session db.Client, packageModel *models.PackageModel) error {
					return errors.New("Failed to insert package")
				},
				PackageRetriever: func(session db.Client) ([]*models.PackageModel, error) {
					return generateRandomPackageModels(10), nil
				},
				RequestTimeBuffer: 0,
				NewGithubRequestService: func(args github.RequestServiceArgs) github.RequestService {
					m := github.NewMockRequestService()
					m.On("FetchGitHubDataForPackageModel", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Return(dtos.GithubRepo{}, errors.New("this is an error"))
					return m
				},
			})

			So(err, ShouldBeNil)
		})
	})
}
