package verdeps

import (
	"fmt"

	"sync"

	errs "github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/io"
)

// packageDirTraverser is a function type that de-couples verdeps.readPackageDir
// from verdeps.traversePackageDir.
type packageDirTraverser func(args traversePackageDirArgs)

// readPackageDirArgs is the arguments struct for readPackageDirArgsArgs.
type readPackageDirArgs struct {
	io                       io.IO
	errors                   *errs.SyncedErrors
	importCounts             *syncedImportCounts
	packageDirPath           string
	importSpecChan           chan *importSpec
	packageSpecChan          chan *packageSpec
	traversePackageDir       packageDirTraverser
	generatedInternalDirName string
}

// readPackageDir deduces important Go import and package metadata in order for
// the remainder of the dependency versioner codebase to "gophrify" the
// appropriate imports in said package.
func readPackageDir(args readPackageDirArgs) {
	// Create a localized error list.
	var (
		errs      = errs.NewSyncedErrors()
		waitGroup = &sync.WaitGroup{}
	)

	// Traverse the directory tree looking for go files - all the while properly
	// handling go package vendoring.
	waitGroup.Add(1)
	args.traversePackageDir(traversePackageDirArgs{
		io:                       args.io,
		errors:                   errs,
		dirPath:                  args.packageDirPath,
		waitGroup:                waitGroup,
		subDirPath:               "",
		parseGoFile:              parseGoFile,
		inVendorDir:              false,
		importCounts:             args.importCounts,
		vendorContext:            newVendorContext(),
		importSpecChan:           args.importSpecChan,
		packageSpecChan:          args.packageSpecChan,
		generatedInternalDirName: args.generatedInternalDirName,
	})

	// Wait for traversal to finish before concocting an error.
	waitGroup.Wait()

	// Compose any errors that there may be into one error.
	if errs.Len() > 0 {
		args.errors.Add(errs.Compose(fmt.Sprintf("Failed to read package directory \"%s\"", args.packageDirPath)))
	}

	// We're done. Time to close the output channels.
	close(args.importSpecChan)
	close(args.packageSpecChan)
}
