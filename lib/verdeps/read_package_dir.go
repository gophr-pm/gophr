package verdeps

import (
	"fmt"
	"runtime"
	"sync"

	"github.com/gophr-pm/gophr/lib/errors"

	errs "github.com/gophr-pm/gophr/lib/errors"
	"github.com/gophr-pm/gophr/lib/io"
)

// readPackageDirArgs is the arguments struct for readPackageDirArgsArgs.
type readPackageDirArgs struct {
	io                       io.IO
	errors                   *errs.SyncedErrors
	importCounts             *syncedImportCounts
	packageDirPath           string
	importSpecChan           chan *importSpec
	packageSpecChan          chan *packageSpec
	generatedInternalDirName string
}

// readPackageDir deduces important Go import and package metadata in order for
// the remainder of the dependency versioner codebase to "gophrify" the
// appropriate imports in said package.
func readPackageDir(args readPackageDirArgs) {
	// Create a localized error list.
	var (
		errs      = newSyncedErrors()
		waitGroup = &sync.WaitGroup{}
	)

	// Traverse the directory tree looking for go files - all the while properly
	// handling go package vendoring.
	waitGroup.Add(1)
	traversePackageDir(traversePackageDirArgs{
		io:                       args.io,
		errors:                   errs,
		dirPath:                  args.packageDirPath,
		waitGroup:                waitGroup,
		subDirPath:               "",
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

	// Clean up after ourselves (traversal allocates a lot of memory).
	runtime.GC()
}
