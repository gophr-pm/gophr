package verdeps

import (
	"bytes"
	"errors"
	"strconv"
	"sync"
)

// readPackageDirArgs is the arguments struct for readPackageDirArgsArgs.
type readPackageDirArgs struct {
	errors          *syncedErrors
	importCounts    *syncedImportCounts
	packageDirPath  string
	importSpecChan  chan *importSpec
	packageSpecChan chan *packageSpec
}

// readPackageDir deduces important Go import and package metadata in order for
// the remainder of the dependency versioner codebase to "gophrify" the
// appropriate imports in said package.
func readPackageDir(args readPackageDirArgs) {
	// Close the input channels on exit.
	defer close(args.importSpecChan)
	defer close(args.packageSpecChan)

	// Create a localized error list.
	errs := newSyncedErrors()
	waitGroup := &sync.WaitGroup{}

	// Traverse the directory tree looking for go files - all the while properly
	// handling go package vendoring.
	waitGroup.Add(1)
	traversePackageDir(traversePackageDirArgs{
		errors:          errs,
		dirPath:         args.packageDirPath,
		waitGroup:       waitGroup,
		subDirPath:      "",
		inVendorDir:     false,
		importCounts:    args.importCounts,
		vendorContext:   newVendorContext(),
		importSpecChan:  args.importSpecChan,
		packageSpecChan: args.packageSpecChan,
	})

	// Wait for traversal to finish before concocting an error.
	waitGroup.Wait()

	// Compose any errors that there may be into one error.
	if errs.len() > 0 {
		var buffer bytes.Buffer
		buffer.WriteString("Failed to read package directory \"")
		buffer.WriteString(args.packageDirPath)
		buffer.WriteString("\" due to ")
		buffer.WriteString(strconv.Itoa(errs.len()))
		buffer.WriteString(" error(s) with file system traversal: [ ")
		rawErrors := errs.errors
		for i, err := range rawErrors {
			if i > 0 {
				buffer.WriteString(", ")
			}

			buffer.WriteString(err.Error())
		}
		buffer.WriteString(" ]")

		args.errors.add(errors.New(buffer.String()))
	}
}

// func renameInternalDirectory(dirPath, parentDirPath, internalRenameTarget string) error {
// 	if err := os.Rename(
// 		dirPath,
// 		filepath.Join(parentDirPath, internalRenameTarget),
// 	); err != nil {
// 		return err
// 	}
//
// 	return nil
// }
