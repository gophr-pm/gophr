package awesome

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAwesomeIndex(t *testing.T) {
	Convey("The awesome indexer should run", t, func() {

		Convey("if we fail to fetch packages from awesome go, we should return an error", func() {
			err := errors.New("Failed to retrieve awesome-go markdown")
			errs := make(chan error, 1)

			index(indexArgs{
				errs: errs,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					return nil, err
				},
			})

			So((<-errs).Error(), ShouldContainSubstring, err.Error())
			close(errs)
		})

		Convey("if we fail to persist packages, we should return an error", func() {
			err := errors.New("Failed to persist packages")
			errs := make(chan error, 1)

			index(indexArgs{
				errs: errs,
				packageFetcher: func(fetchAwesomeGoListArgs) ([]packageTuple, error) {
					return generateRandomAwesomePackages(201), nil
				},
				persistPackages: func(persistAwesomePackagesArgs) error {
					return err
				},
			})

			So((<-errs).Error(), ShouldContainSubstring, err.Error())
			close(errs)
		})

		Convey("if we successfully retrieve 10 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				errs    = make(chan error, 1)
				pkgs    = generateRandomAwesomePackages(10)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Run the main awsome index code.
			index(indexArgs{
				errs: errs,
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

			// Verify number of packages.
			So(len(pkgs), ShouldEqual, 10)
			// Verify that there were no errors.
			So(len(errs), ShouldEqual, 0)

			close(errs)
		})

		Convey("if we successfully retrieve 50 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				errs    = make(chan error, 1)
				pkgs    = generateRandomAwesomePackages(50)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Run the main awsome index code.
			index(indexArgs{
				errs: errs,
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

			// Verify number of packages.
			So(len(pkgs), ShouldEqual, 50)
			// Verify that there were no errors.
			So(len(errs), ShouldEqual, 0)

			close(errs)
		})

		Convey("if we successfully retrieve 51 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				errs    = make(chan error, 1)
				pkgs    = generateRandomAwesomePackages(51)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Run the main awsome index code.
			index(indexArgs{
				errs: errs,
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

			// Verify number of packages.
			So(len(pkgs), ShouldEqual, 51)
			// Verify that there were no errors.
			So(len(errs), ShouldEqual, 0)

			close(errs)
		})

		Convey("if we successfully retrieve 107 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				errs    = make(chan error, 1)
				pkgs    = generateRandomAwesomePackages(107)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Run the main awsome index code.
			index(indexArgs{
				errs: errs,
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

			// Verify number of packages.
			So(len(pkgs), ShouldEqual, 107)
			// Verify that there were no errors.
			So(len(errs), ShouldEqual, 0)

			close(errs)
		})

		Convey("if we succeed retrieve 201 packages, we should persist every package and return nil", func() {
			var (
				c       = make(chan packageTuple)
				errs    = make(chan error, 1)
				pkgs    = generateRandomAwesomePackages(201)
				pkgsMap = generateMapOfAwesomePackages(pkgs)
			)

			// Run the main awsome index code.
			index(indexArgs{
				errs: errs,
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
			// Verify that there were no errors.
			So(len(errs), ShouldEqual, 0)

			close(errs)
		})
	})
}
