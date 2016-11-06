package awesome

import (
	"errors"
	"testing"

	"github.com/gophr-pm/gophr/scheduler/worker/common"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestAwesomeIndex(t *testing.T) {
	Convey("The awesome indexer should run", t, func() {

		Convey("if we fail to fetch packages from awesome go, we should return an error", func() {
			err := errors.New("Failed to retrieve awesome-go markdown")
			logger := common.NewMockJobLogger()
			// Should communicate that the job has started.
			logger.On("Info", mock.AnythingOfType("string"))
			// Should try to log the error.
			logger.On("Errorf", mock.AnythingOfType("string"), err)

			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					return nil, err
				},
			})
		})

		Convey("if we fail to persist packages, we should return an error", func() {
			err := errors.New("Failed to persist packages")
			logger := common.NewMockJobLogger()
			// Should communicate that the job has started.
			logger.On("Info", mock.AnythingOfType("string"))
			// Should try to log the error.
			logger.On("Errorf", mock.AnythingOfType("string"), err)

			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					return generateRandomAwesomePackages(201), nil
				},
				persistPackages: func(persistAwesomePackagesArgs) error {
					return err
				},
			})
		})

		Convey("if we successfully retrieve 10 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				pkgs    = generateRandomAwesomePackages(10)
				logger  = common.NewMockJobLogger()
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Should communicate that the job has started and then ended.
			logger.On("Info", mock.AnythingOfType("string"))

			// Run the main awsome index code.
			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				persistPackages: func(args persistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to
					// verify all packages are accounted for.
					go func(pkgs []packageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.packageTuples)
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
		})

		Convey("if we successfully retrieve 50 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				pkgs    = generateRandomAwesomePackages(50)
				logger  = common.NewMockJobLogger()
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Should communicate that the job has started and then ended.
			logger.On("Info", mock.AnythingOfType("string"))

			// Run the main awsome index code.
			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				persistPackages: func(args persistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to
					// verify all packages are accounted for.
					go func(pkgs []packageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.packageTuples)
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
		})

		Convey("if we successfully retrieve 51 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				pkgs    = generateRandomAwesomePackages(51)
				logger  = common.NewMockJobLogger()
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Should communicate that the job has started and then ended.
			logger.On("Info", mock.AnythingOfType("string"))

			// Run the main awsome index code.
			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				persistPackages: func(args persistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to
					// verify all packages are accounted for.
					go func(pkgs []packageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.packageTuples)
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
		})

		Convey("if we successfully retrieve 107 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				pkgs    = generateRandomAwesomePackages(107)
				logger  = common.NewMockJobLogger()
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Should communicate that the job has started and then ended.
			logger.On("Info", mock.AnythingOfType("string"))

			// Run the main awsome index code.
			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				persistPackages: func(args persistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to
					// verify all packages are accounted for.
					go func(pkgs []packageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.packageTuples)
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
		})

		Convey("if we succeed retrieve 201 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				pkgs    = generateRandomAwesomePackages(201)
				logger  = common.NewMockJobLogger()
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Should communicate that the job has started and then ended.
			logger.On("Info", mock.AnythingOfType("string"))

			// Run the main awsome index code.
			index(indexArgs{
				logger: logger,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					// Return the list of random packages generate before the index call.
					return pkgs, nil
				},
				persistPackages: func(args persistAwesomePackagesArgs) error {
					// For each batch sent to be persisted, send them through a channel to
					// verify all packages are accounted for.
					go func(pkgs []packageTuple) {
						for _, pkg := range pkgs {
							c <- pkg
						}
					}(args.packageTuples)
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
		})
	})
}
