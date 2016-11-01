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

		Convey("if we fail to fetch packages from awesome go, we should return an error", func() {
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

		Convey("if we successfully retrieve 10, we should persist every package and return nil", func() {
			var (
				pkgs    = generateRandomAwesomePackages(10)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
				c       = make(chan PackageTuple)
			)

			// Run the main awsome index code.
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				PersistPackages: func(args PersistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to verify
					// all packages are accounted for.
					go func(pkgs []PackageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.PackageTuples)
					return nil
				},
			})

			// Check to make sure every package has been persisted.
			pkgCount := 0
			for pkg := range c {
				if pkgCount == len(pkgs)-1 {
					close(c)
					break
				}
				lookUpKey := pkg.author + "/" + pkg.repo
				So(pkg, ShouldResemble, pkgsMap[lookUpKey])
				pkgCount++
			}

			// Verify number of packages
			So(len(pkgs), ShouldEqual, 10)
			// Finally verify there haven't been any errors.
			So(err, ShouldBeNil)
		})

		Convey("if we successfully retrieve 50, we should persist every package and return nil", func() {
			var (
				pkgs    = generateRandomAwesomePackages(50)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
				c       = make(chan PackageTuple)
			)

			// Run the main awsome index code.
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				PersistPackages: func(args PersistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to verify
					// all packages are accounted for.
					go func(pkgs []PackageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.PackageTuples)
					return nil
				},
			})

			// Check to make sure every package has been persisted.
			pkgCount := 0
			for pkg := range c {
				if pkgCount == len(pkgs)-1 {
					close(c)
					break
				}
				lookUpKey := pkg.author + "/" + pkg.repo
				So(pkg, ShouldResemble, pkgsMap[lookUpKey])
				pkgCount++
			}

			// Verify number of packages
			So(len(pkgs), ShouldEqual, 50)
			// Finally verify there haven't been any errors.
			So(err, ShouldBeNil)
		})

		Convey("if we successfully retrieve 51, we should persist every package and return nil", func() {
			var (
				pkgs    = generateRandomAwesomePackages(51)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
				c       = make(chan PackageTuple)
			)

			// Run the main awsome index code.
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				PersistPackages: func(args PersistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to verify
					// all packages are accounted for.
					go func(pkgs []PackageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.PackageTuples)
					return nil
				},
			})

			// Check to make sure every package has been persisted.
			pkgCount := 0
			for pkg := range c {
				if pkgCount == len(pkgs)-1 {
					close(c)
					break
				}
				lookUpKey := pkg.author + "/" + pkg.repo
				So(pkg, ShouldResemble, pkgsMap[lookUpKey])
				pkgCount++
			}

			// Verify number of packages
			So(len(pkgs), ShouldEqual, 51)
			// Finally verify there haven't been any errors.
			So(err, ShouldBeNil)
		})

		Convey("if we successfully retrieve 107, we should persist every package and return nil", func() {
			var (
				pkgs    = generateRandomAwesomePackages(107)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
				c       = make(chan PackageTuple)
			)

			// Run the main awsome index code.
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				PersistPackages: func(args PersistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to verify
					// all packages are accounted for.
					go func(pkgs []PackageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.PackageTuples)
					return nil
				},
			})

			// Check to make sure every package has been persisted.
			pkgCount := 0
			for pkg := range c {
				if pkgCount == len(pkgs)-1 {
					close(c)
					break
				}
				lookUpKey := pkg.author + "/" + pkg.repo
				So(pkg, ShouldResemble, pkgsMap[lookUpKey])
				pkgCount++
			}

			// Verify number of packages
			So(len(pkgs), ShouldEqual, 107)
			// Finally verify there haven't been any errors.
			So(err, ShouldBeNil)
		})

		Convey("if we succeed retrieve 201, we should persist every package and return nil", func() {
			var (
				pkgs    = generateRandomAwesomePackages(201)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
				c       = make(chan PackageTuple)
			)

			// Run the main awsome index code.
			err := Index(IndexArgs{
				Init: func() (*config.Config, db.Client) {
					c := db.NewMockClient()
					c.On("Close").Return()
					return &config.Config{}, c
				},
				PackageFetcher: func(FetchAwesomeGoListArgs) ([]PackageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				PersistPackages: func(args PersistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to verify
					// all packages are accounted for.
					go func(pkgs []PackageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.PackageTuples)
					return nil
				},
			})

			// Check to make sure every package has been persisted.
			pkgCount := 0
			for pkg := range c {
				if pkgCount == len(pkgs)-1 {
					close(c)
					break
				}
				lookUpKey := pkg.author + "/" + pkg.repo
				So(pkg, ShouldResemble, pkgsMap[lookUpKey])
				pkgCount++
			}

			// Verify number of packages
			So(len(pkgs), ShouldEqual, 201)
			// Finally verify there haven't been any errors.
			So(err, ShouldBeNil)
		})
	})
}
