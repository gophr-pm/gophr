package verdeps

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
)

// shaFetcher is a function type that de-couples verdeps.processDeps from
// verdeps.fetchSHA.
type shaFetcher func(args fetchSHAArgs)

// depsReviser is a function type that de-couples verdeps.processDeps from
// verdeps.reviseDeps.
type depsReviser func(args reviseDepsArgs)

// packageDirReader is a function type that de-couples verdeps.processDeps from
// verdeps.readPackageDir.
type packageDirReader func(args readPackageDirArgs)

// specWaitingListCreator is a function type that de-couples verdeps.processDeps
// from newSpecWaitingList.
type specWaitingListCreator func(
	initialSpecs ...*importSpec,
) specWaitingList

type syncedStringMapCreator func() syncedStringMap

// processDepsArgs is the arguments struct for processDeps.
type processDepsArgs struct {
	io                 io.IO
	ghSvc              github.RequestService
	fetchSHA           shaFetcher
	reviseDeps         depsReviser
	packageSHA         string
	packagePath        string
	packageRepo        string
	packageAuthor      string
	readPackageDir     packageDirReader
	packageVersionDate time.Time
	newSpecWaitingList specWaitingListCreator
	newSyncedStringMap syncedStringMapCreator
}

// TODO(skeswa): add a descriptive comment.
func processDeps(args processDepsArgs) error {
	var (
		revisionChan             = make(chan *revision)
		waitingSpecs             = newSyncedWaitingListMap()
		importPathSHAs           = args.newSyncedStringMap()
		importSpecChan           = make(chan *importSpec)
		packageSpecChan          = make(chan *packageSpec)
		importPathSHAChan        = make(chan *importPathSHA)
		accumulatedErrors        = newSyncedErrors()
		revisionWaitGroup        = &sync.WaitGroup{}
		pendingSHARequests       = newSyncedInt()
		syncedImportCounts       = newSyncedImportCounts()
		generatedInternalDirName = generateInternalDirName()
	)

	// Read the package looking for import and package metadata.
	go args.readPackageDir(readPackageDirArgs{
		io:                       args.io,
		errors:                   accumulatedErrors,
		importCounts:             syncedImportCounts,
		packageDirPath:           args.packagePath,
		importSpecChan:           importSpecChan,
		packageSpecChan:          packageSpecChan,
		generatedInternalDirName: generatedInternalDirName,
	})

	// Revise dependencies in the go source files.
	revisionWaitGroup.Add(1)
	go args.reviseDeps(reviseDepsArgs{
		inputChan:          revisionChan,
		revisionWaitGroup:  revisionWaitGroup,
		accumulatedErrors:  accumulatedErrors,
		syncedImportCounts: syncedImportCounts,
	})

	// Process incoming data from the channels.
	for {
		select {
		case ips, alive := <-importPathSHAChan:
			// Nil the channel if no longer alive.
			if !alive {
				importPathSHAChan = nil
				break
			}

			// Create an entry in the map.
			importPathHash := importPathHashOf(ips.importPath)
			importPathSHAs.set(importPathHash, ips.sha)

			// Clear away the waiting specs.
			if waitingList, exists := waitingSpecs.get(importPathHash); exists {
				// There is a waiting list, so it needs to be cleared.
				if specs := waitingList.clear(); specs != nil {
					for _, spec := range specs {
						enqueueImportRevision(
							revisionChan,
							spec.imports.Path.Value,
							ips.sha,
							generatedInternalDirName,
							spec)
					}
				}
			}

			// Check if this is the last time that an import path sha will come through.
			// We check this by checking if the spec channel has already closed, and
			// if there are no pending sha requests. If so, close this channel.
			if pendingSHARequests.value() == 0 && importSpecChan == nil {
				closeImportPathSHAChan(importPathSHAChan, waitingSpecs)
			}

		case spec, alive := <-packageSpecChan:
			// Nil the channel if no longer alive.
			if !alive {
				packageSpecChan = nil
				break
			}

			// Pass the package revision along to the file-writing stage.
			enqueuePackageRevision(revisionChan, spec)

		case spec, alive := <-importSpecChan:
			// Nil the channel if no longer alive.
			if !alive {
				importSpecChan = nil

				// If there are no pending sha requests, then there never will be since
				// this channel is closing. Therefore, we can close the sha channel.
				if shouldCloseImportPathSHAChan(
					pendingSHARequests,
					importSpecChan,
					importPathSHAChan,
				) {
					closeImportPathSHAChan(importPathSHAChan, waitingSpecs)
				}

				break
			}

			// Signal that a request might be about to begin to prevent race conditions.
			// This is done regardless of whether a request actually occurs to prevent
			// pendingSHARequests.value() from being evaluated between the beginning of
			// the case statement and go args.fetchSHA(...).
			pendingSHARequests.increment()

			// For each incoming spec, make it wait keyed on the import path hash.
			importPath := spec.imports.Path.Value
			importPathHash := importPathHashOf(spec.imports.Path.Value)
			if sha, exists := importPathSHAs.get(importPathHash); !exists {
				// If we don't presently have the sha, then we have to go out and get
				// it.
				if specs, exists := waitingSpecs.get(importPathHash); !exists {
					// Create a new waiting list for this import path since it does not
					// yet exist.
					waitingSpecs.setIfAbsent(
						importPathHash,
						args.newSpecWaitingList(spec))

					// Start the request itself.
					go args.fetchSHA(fetchSHAArgs{
						ghSvc:              args.ghSvc,
						outputChan:         importPathSHAChan,
						importPath:         importPath,
						packageSHA:         args.packageSHA,
						packageRepo:        args.packageRepo,
						packageAuthor:      args.packageAuthor,
						pendingSHARequests: pendingSHARequests,
						packageVersionDate: args.packageVersionDate,
					})
				} else {
					// It is now clear that no pending request began here.
					pendingSHARequests.decrement()

					if ok := specs.add(spec); !ok {
						// If the add failed, assume that it is because the the sha was
						// obtained after we last checked.
						if sha, exists = importPathSHAs.get(importPathHash); !exists {
							accumulatedErrors.add(fmt.Errorf(
								"Could not version dependency %s because the SHA did not yet exist.",
								importPath))
						} else {
							enqueueImportRevision(
								revisionChan,
								importPath,
								sha,
								generatedInternalDirName,
								spec)
						}
					}
				}
			} else {
				// It is now clear that no pending request began here.
				pendingSHARequests.decrement()

				// If we got here, it means that the sha has already been obtained, so
				// the new import path exists.
				enqueueImportRevision(
					revisionChan,
					importPath,
					sha,
					generatedInternalDirName,
					spec)
			}
		}

		// If both of the channels being selected are nil, its time to stop
		// selecting.
		if importPathSHAChan == nil &&
			packageSpecChan == nil &&
			importSpecChan == nil {
			// At this point, the deps chan should be closed since all potential deps
			// have been seen and processed.
			close(revisionChan)
			revisionChan = nil

			// Wait until all the deps are revised before continuing.
			revisionWaitGroup.Wait()

			// Exit the inifinte loop.
			break
		}
	}

	// If there were any errors, compose all the errors into on error.
	if accumulatedErrors.len() > 0 {
		return concatErrors(accumulatedErrors)
	}

	// Otherwise, return without a hitch.
	return nil
}

// enqueueImportRevision is a helper function that puts a revision into the
// revision channel that revises an import statement.
func enqueueImportRevision(
	revisionChan chan *revision,
	importPath string,
	sha string,
	generatedInternalDirName string,
	spec *importSpec,
) {
	author, repo, subpath := parseImportPath(importPath)
	newImportPath := composeNewImportPath(
		author,
		repo,
		sha,
		subpath,
		generatedInternalDirName)

	revisionChan <- newImportRevision(spec, newImportPath)
}

// enqueuePackageRevision is a helper function that puts a revision into the
// revision channel that (potentially) revises a package statement.
func enqueuePackageRevision(revisionChan chan *revision, spec *packageSpec) {
	revisionChan <- newPackageRevision(spec)
}

// closeImportPathSHAChan closes the importPathSHAChan and clears all spec
// waiting list. The waiting lists are cleared because they're waiting for a SHA
// that will presumably never come.
func closeImportPathSHAChan(
	importPathSHAChan chan *importPathSHA,
	waitingSpecs *syncedWaitingListMap) {
	waitingSpecs.each(func(importPath string, waitingList specWaitingList) {
		waitingList.clear()
	}).clear()

	// Finally, close the channel.
	close(importPathSHAChan)
}

// shouldCloseImportPathSHAChan returns true if the importPathSHAChan should be
// closed.
func shouldCloseImportPathSHAChan(
	pendingSHARequests *syncedInt,
	importSpecChan chan *importSpec,
	importPathSHAChan chan *importPathSHA) bool {
	return pendingSHARequests.value() == 0 &&
		importSpecChan == nil &&
		importPathSHAChan != nil
}

// concatErrors joins all the accumulated individual errors into one combined
// error.
func concatErrors(accumulatedErrors *syncedErrors) error {
	errs := accumulatedErrors.get()
	buffer := bytes.Buffer{}

	buffer.WriteString("Failed to process dependencies. Bumped into ")
	buffer.WriteString(strconv.Itoa(len(errs)))
	buffer.WriteString(" problems: [ ")

	for i, err := range errs {
		if i > 0 {
			buffer.WriteString(", ")
		}
		buffer.WriteString(err.Error())
	}

	buffer.WriteString(" ].")

	return errors.New(buffer.String())
}
