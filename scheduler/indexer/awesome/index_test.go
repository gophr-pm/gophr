package awesome

import (
	"errors"
	"testing"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	. "github.com/smartystreets/goconvey/convey"
)

func TestAwesomeIndex(t *testing.T) {
	Convey("The awesome indexer should run", t, func() {

		Convey("if we fail to fetch packages from awesome go, we should fail", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					return nil, errors.New("Failed to retrieve awesome-go markdown")
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("if we fail to persist packages, we should fail", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					return generateRandomAwesomePackages(201), nil
				},
				PersistPackages: func(PersistAwesomePackagesArgs) error {
					return errors.New("Failed to persist packages")
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("if we succeed, we should return nil", func() {
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					return generateRandomAwesomePackages(201), nil
				},
				PersistPackages: func(PersistAwesomePackagesArgs) error {
					return nil
				},
			})

			So(err, ShouldBeNil)
		})
	})
}
