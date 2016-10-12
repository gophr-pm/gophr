package verdeps

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sync"
	"github.com/gophr-pm/gophr/common/io"
)

// traversePackageDir is the arguments struct for traversePackageDirArgs.
type traversePackageDirArgs struct {
	errors                   *syncedErrors
	dirPath                  string
	waitGroup                *sync.WaitGroup
	subDirPath               string
	inVendorDir              bool
	importCounts             *syncedImportCounts
	vendorContext            *vendorContext
	importSpecChan           chan *importSpec
	packageSpecChan          chan *packageSpec
	generatedInternalDirName string
}

// traversePackageDir recursively traverses a package directory in order to find
// Go imports and package metadata. Once the aforesaid has been found, it is
// returned through the import and package spec channels. traversePackageDir
// also is fully aware of Go package vendoring and will ignore vendored
// packages accordingly.
func traversePackageDir(args traversePackageDirArgs) {
	var (
		err         error
		files       []os.FileInfo
		fullDirPath = filepath.Join(args.dirPath, args.subDirPath)
	)

	// Signal that the wait group should continue when this function exits.
	defer args.waitGroup.Done()

	// Read all the files of this directory.
	if files, err = ioutil.ReadDir(fullDirPath); err != nil {
		args.errors.add(err)
		return
	}

	verdepsHelper := &verdepsHelperArgs{io: io.NewIO()}
	// Get all relevant pathing information in one fell swoop.
	vendorDirPath, subDirNames, goFilePaths := verdepsHelper.getPackageDirPaths(files, fullDirPath)

	// Record this subpath as a vendored package.
	if args.inVendorDir &&
		len(args.subDirPath) > 0 &&
		args.vendorContext != nil {
		// Add this vendored package to the *next* context since vendored packages
		// in the same vendor directory do not vendor for themselves.
		args.vendorContext.add(args.subDirPath)
	}

	// If there is a vendor dir, then traverse it first.
	if len(vendorDirPath) > 0 {
		// Create or spawn a vendor context.
		currentVendorContext := args.vendorContext.spawnChildContext()

		// Update the vendor context of this directory traversal.
		args.vendorContext = currentVendorContext

		// Explore the new vendor directory synchronously since it needs to be
		// traversed before every other directory.
		subErrors := newSyncedErrors()
		subImportSpecChan := make(chan *importSpec)
		subPackageSpecChan := make(chan *packageSpec)
		traversePackageDirWaitGroup := &sync.WaitGroup{}
		bufferVendorablesWaitGroup := &sync.WaitGroup{}

		// Traverse the vendor dir.
		traversePackageDirWaitGroup.Add(1)
		go traversePackageDir(traversePackageDirArgs{
			errors:          subErrors,
			dirPath:         vendorDirPath,
			waitGroup:       traversePackageDirWaitGroup,
			subDirPath:      "",
			inVendorDir:     true,
			importCounts:    nil,
			vendorContext:   currentVendorContext,
			importSpecChan:  subImportSpecChan,
			packageSpecChan: subPackageSpecChan,
		})

		// The specs must be buffered since vendored packages are self-referencial.
		bufferVendorablesWaitGroup.Add(1)
		go bufferVendorables(bufferVendorablesArgs{
			waitGroup:             bufferVendorablesWaitGroup,
			importCounts:          args.importCounts,
			inputImportSpecChan:   subImportSpecChan,
			outputImportSpecChan:  args.importSpecChan,
			currentVendorContext:  currentVendorContext,
			inputPackageSpecChan:  subPackageSpecChan,
			outputPackageSpecChan: args.packageSpecChan,
		})

		// Wait for the vendor traversal to finish.
		traversePackageDirWaitGroup.Wait()

		// Finalize the vendor context so it can no longer be modified.
		currentVendorContext.finalize()

		// Close the subImportSpecChan so that bufferVendorables can exit.
		close(subImportSpecChan)
		close(subPackageSpecChan)

		// Now, wait for all the vendorable import specs to be retried.
		bufferVendorablesWaitGroup.Wait()

		// Exit if there were problems.
		if subErrors.len() > 0 {
			fmt.Printf("error founds found while traversing %s\n", vendorDirPath)
			args.errors.add(subErrors.get()...)
			return
		}
	}

	// Create wait group on child file operations.
	subWaitGroup := &sync.WaitGroup{}

	// Explore all child directories that aren't "vendor".
	for _, subDirName := range subDirNames {
		// Internal directories must be renamed for gophr to function correctly.
		if subDirName == internalDirName {
			if err := os.Rename(
				filepath.Join(args.dirPath, args.subDirPath, internalDirName),
				filepath.Join(args.dirPath, args.subDirPath, args.generatedInternalDirName),
			); err != nil {
				// If there was a problem performing the rename, exit immediately.
				args.errors.add(err)
				return
			}

			// Obviously, the subDirName must now change accordingly.
			subDirName = args.generatedInternalDirName
		}

		subWaitGroup.Add(1)
		go traversePackageDir(traversePackageDirArgs{
			errors:          args.errors,
			dirPath:         args.dirPath,
			waitGroup:       subWaitGroup,
			subDirPath:      filepath.Join(args.subDirPath, subDirName),
			inVendorDir:     args.inVendorDir,
			importCounts:    args.importCounts,
			vendorContext:   args.vendorContext,
			importSpecChan:  args.importSpecChan,
			packageSpecChan: args.packageSpecChan,
		})
	}

	// Process all of the go files, looking for transformable imports.
	for _, goFilePath := range goFilePaths {
		subWaitGroup.Add(1)
		go parseGoFile(parseGoFileArgs{
			errors:          args.errors,
			filePath:        goFilePath,
			waitGroup:       subWaitGroup,
			importCounts:    args.importCounts,
			vendorContext:   args.vendorContext,
			importSpecChan:  args.importSpecChan,
			packageSpecChan: args.packageSpecChan,
		})
	}

	// Wait for child file operations to complete before exiting.
	subWaitGroup.Wait()
}
