package verdeps

import (
	"errors"
	"testing"

	"github.com/gophr-pm/gophr/lib/io"
	. "github.com/smartystreets/goconvey/convey"
)

func TestReadPackageDir(t *testing.T) {
	Convey("Given a package dir to be read", t, func() {
		Convey("When there are no errors in traversal, no errors should be returned", func() {
			var (
				io                           = io.NewMockIO()
				errors                       = newSyncedErrors()
				importCounts                 = newSyncedImportCounts()
				packageDirPath               = "/this/is/a/package/dir"
				importSpecChan               = make(chan *importSpec)
				packageSpecChan              = make(chan *packageSpec)
				traversePackageDir           packageDirTraverser
				generatedInternalDirName     = "thisisagenerateddirname"
				actualTraversePackageDirArgs traversePackageDirArgs
			)

			traversePackageDir = func(args traversePackageDirArgs) {
				actualTraversePackageDirArgs = args
				args.waitGroup.Done()
			}

			readPackageDir(readPackageDirArgs{
				io:                       io,
				errors:                   errors,
				importCounts:             importCounts,
				packageDirPath:           packageDirPath,
				importSpecChan:           importSpecChan,
				packageSpecChan:          packageSpecChan,
				traversePackageDir:       traversePackageDir,
				generatedInternalDirName: generatedInternalDirName,
			})

			So(actualTraversePackageDirArgs.dirPath, ShouldEqual, packageDirPath)
			So(actualTraversePackageDirArgs.errors, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.generatedInternalDirName, ShouldEqual, generatedInternalDirName)
			So(actualTraversePackageDirArgs.importCounts, ShouldEqual, importCounts)
			So(actualTraversePackageDirArgs.importSpecChan, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.packageSpecChan, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.parseGoFile, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.subDirPath, ShouldBeEmpty)
			So(actualTraversePackageDirArgs.waitGroup, ShouldNotBeNil)

			So(actualTraversePackageDirArgs.errors.len(), ShouldEqual, 0)
		})

		Convey("When there are errors in traversal, an error should be returned", func() {
			var (
				io                           = io.NewMockIO()
				errs                         = newSyncedErrors()
				importCounts                 = newSyncedImportCounts()
				packageDirPath               = "/this/is/a/package/dir"
				importSpecChan               = make(chan *importSpec)
				packageSpecChan              = make(chan *packageSpec)
				traversePackageDir           packageDirTraverser
				generatedInternalDirName     = "thisisagenerateddirname"
				actualTraversePackageDirArgs traversePackageDirArgs
			)

			traversePackageDir = func(args traversePackageDirArgs) {
				actualTraversePackageDirArgs = args
				args.errors.add(errors.New("error1"))
				args.errors.add(errors.New("error2"))
				args.errors.add(errors.New("error3"))
				args.waitGroup.Done()
			}

			readPackageDir(readPackageDirArgs{
				io:                       io,
				errors:                   errs,
				importCounts:             importCounts,
				packageDirPath:           packageDirPath,
				importSpecChan:           importSpecChan,
				packageSpecChan:          packageSpecChan,
				traversePackageDir:       traversePackageDir,
				generatedInternalDirName: generatedInternalDirName,
			})

			So(actualTraversePackageDirArgs.dirPath, ShouldEqual, packageDirPath)
			So(actualTraversePackageDirArgs.errors, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.generatedInternalDirName, ShouldEqual, generatedInternalDirName)
			So(actualTraversePackageDirArgs.importCounts, ShouldEqual, importCounts)
			So(actualTraversePackageDirArgs.importSpecChan, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.packageSpecChan, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.parseGoFile, ShouldNotBeNil)
			So(actualTraversePackageDirArgs.subDirPath, ShouldBeEmpty)
			So(actualTraversePackageDirArgs.waitGroup, ShouldNotBeNil)

			So(errs.len(), ShouldEqual, 1)
			errStr := (errs.get()[0]).Error()
			So(errStr, ShouldContainSubstring, "error1")
			So(errStr, ShouldContainSubstring, "error2")
			So(errStr, ShouldContainSubstring, "error3")
		})
	})
}
