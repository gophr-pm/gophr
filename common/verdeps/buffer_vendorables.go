package verdeps

import (
	"strings"
	"sync"
)

// bufferVendorablesArgs is the arguments struct for bufferVendorables.
type bufferVendorablesArgs struct {
	waitGroup             *sync.WaitGroup
	importCounts          *syncedImportCounts
	inputImportSpecChan   chan *importSpec
	outputImportSpecChan  chan *importSpec
	currentVendorContext  *vendorContext
	inputPackageSpecChan  chan *packageSpec
	outputPackageSpecChan chan *packageSpec
}

// bufferVendorables buffers the import and package spec channels until
// the vendoring context has been finalized. At which point, all the buffered
// specs are sent out through the output channels.
func bufferVendorables(args bufferVendorablesArgs) {
	var (
		collectedImportSpecs  []*importSpec
		collectedPackageSpecs []*packageSpec
		unvendoredImportSpecs = make(map[string][]*importSpec)
	)

	// Loop until all of the specs have been buffered.
	for {
		select {
		// Collect specs that could have slipped through the cracks due to the vendor
		// directory not being fully indexed yet.
		case spec, alive := <-args.inputImportSpecChan:
			if !alive {
				args.inputImportSpecChan = nil
				break
			}

			// If the import is a github URL, then it could have been mistaken for a
			// gophrizable import.
			if strings.HasPrefix(spec.imports.Path.Value, "\"github.com/") {
				collectedImportSpecs = append(collectedImportSpecs, spec)
			} else {
				// Let anything not eligible to be collected pass through.
				unvendoredImportSpecs[spec.filePath] = append(
					unvendoredImportSpecs[spec.filePath],
					spec)
			}

		case spec, alive := <-args.inputPackageSpecChan:
			if !alive {
				args.inputPackageSpecChan = nil
				break
			}

			// Collect the package specs.
			collectedPackageSpecs = append(collectedPackageSpecs, spec)
		}

		// Break if both channels have been closed.
		if args.inputImportSpecChan == nil && args.inputPackageSpecChan == nil {
			break
		}
	}

	// Now that the channels have closed, iterate through the collected specs
	// so that unvendored packages can be passed through.
	for _, spec := range collectedImportSpecs {
		// Trim off the quotes from the import string so it can be used in a
		// vendor context look up.
		importString := strings.Trim(spec.imports.Path.Value, "\"")
		// If unvendored, throw the import spec into the map.
		if !args.currentVendorContext.contains(importString) {
			unvendoredImportSpecs[spec.filePath] = append(
				unvendoredImportSpecs[spec.filePath],
				spec)
		}
	}

	// Pass the unvendored packages through to the output chan.
	for filePath, specs := range unvendoredImportSpecs {
		// Signal how many imports to expect.
		args.importCounts.setImportCount(filePath, len(specs))
		// After the import count is set, enqueue both kinds of the specs.
		for _, spec := range specs {
			args.outputImportSpecChan <- spec
		}
	}

	// Get rid of the package specs now that there are import counts.
	for _, spec := range collectedPackageSpecs {
		args.outputPackageSpecChan <- spec
	}

	// Signal that the task is now complete.
	args.waitGroup.Done()
}
