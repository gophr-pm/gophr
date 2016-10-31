package verdeps

import (
	"fmt"

	"log"
	"sync"
	"time"

	"github.com/gophr-pm/gophr/lib/errors"
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

// syncedStringMapCreator is function type that de-couples verdeps.processDeps
// from newSyncedStringMap.
type syncedStringMapCreator func() syncedStringMap

// syncedWaitingListMapCreator is function type that de-couples
// verdeps.processDeps from newSyncedWaitingListMap.
type syncedWaitingListMapCreator func() syncedWaitingListMap

// processDepsArgs is the arguments struct for processDeps.
type processDepsArgs struct {
	io                      io.IO
	ghSvc                   github.RequestService
	fetchSHA                shaFetcher
	reviseDeps              depsReviser
	packageSHA              string
	packagePath             string
	packageRepo             string
	packageAuthor           string
	readPackageDir          packageDirReader
	packageVersionDate      time.Time
	newSpecWaitingList      specWaitingListCreator
	newSyncedStringMap      syncedStringMapCreator
	newSyncedWaitingListMap syncedWaitingListMapCreator
}

// TODO(skeswa): add a descriptive comment.
func processDeps(args processDepsArgs) error {
	var (
		revisionChan             = make(chan *revision)
		waitingSpecs             = args.newSyncedWaitingListMap()
		importSpecChan           = make(chan *importSpec)
		fetchSHAResults          = args.newSyncedStringMap()
		packageSpecChan          = make(chan *packageSpec)
		accumulatedErrors        = errors.NewSyncedErrors()
		revisionWaitGroup        = &sync.WaitGroup{}
		syncedImportCounts       = newSyncedImportCounts()
		fetchSHAResultChan       = make(chan *fetchSHAResult)
		processedImportsCount    = newSyncedInt()
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
		traversePackageDir:       traversePackageDir,
		generatedInternalDirName: generatedInternalDirName,
	})

	// Revise dependencies in the go source files.
	revisionWaitGroup.Add(1)
	go args.reviseDeps(reviseDepsArgs{
		io:                 args.io,
		inputChan:          revisionChan,
		composeBytesDiffs:  composeBytesDiffs,
		revisionWaitGroup:  revisionWaitGroup,
		accumulatedErrors:  accumulatedErrors,
		syncedImportCounts: syncedImportCounts,
	})

	for {
		// Process incoming data from the channels.
		select {
		case result := <-fetchSHAResultChan:
			// Count this import as processed.
			processedImportsCount.increment()

			// Only continue if an actual importPath-sha mapping came through.
			if result.successful {
				// Create an entry in the map.
				importPathHash := importPathHashOf(result.importPath)
				fetchSHAResults.set(importPathHash, result.sha)

				// Clear away the waiting specs.
				if waitingList, exists := waitingSpecs.get(importPathHash); exists {
					// There is a waiting list, so it needs to be cleared.
					if specs := waitingList.clear(); specs != nil {
						for _, spec := range specs {
							enqueueImportRevision(
								revisionChan,
								spec.imports.Path.Value,
								result.sha,
								generatedInternalDirName,
								spec)
						}
					}
				}
			} else {
				// If not successful, log the error (don't return it since it isn't
				// fatal).
				log.Printf(
					`Failed to fetch SHA for import path "%s": %v`,
					result.sha,
					result.err)
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
				break
			}

			// For each incoming spec, make it wait keyed on the import path hash.
			importPath := spec.imports.Path.Value
			importPathHash := importPathHashOf(spec.imports.Path.Value)
			if sha, exists := fetchSHAResults.get(importPathHash); !exists {
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
						outputChan:         fetchSHAResultChan,
						importPath:         importPath,
						packageSHA:         args.packageSHA,
						packageRepo:        args.packageRepo,
						packageAuthor:      args.packageAuthor,
						packageVersionDate: args.packageVersionDate,
					})
				} else {
					if ok := specs.add(spec); !ok {
						// If the add failed, assume that it is because the the sha was
						// obtained after we last checked.
						if sha, exists = fetchSHAResults.get(importPathHash); !exists {
							accumulatedErrors.Add(fmt.Errorf(
								"Could not version dependency %s"+
									" because the SHA did not yet exist.",
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

					// Count this import as processed.
					processedImportsCount.increment()
				}
			} else {
				// If we got here, it means that the sha has already been obtained, so
				// the new import path exists.
				enqueueImportRevision(
					revisionChan,
					importPath,
					sha,
					generatedInternalDirName,
					spec)

				// Count this import as processed.
				processedImportsCount.increment()
			}
		}

		// Wait until the import spec channel in nil because that means that the
		// import count total in correct. Thereafter, check if the total number
		// of processed import path SHAs matches the total number of imports. In
		// this case, it is time to close the import path SHA channel.
		if fetchSHAResultChan != nil &&
			importSpecChan == nil &&
			processedImportsCount.value() == syncedImportCounts.totalCount() {
			// Clear all the remaining waiting lists since none of the remaining
			// import paths will ever be matched to the corresponding SHAs.
			// TODO(skeswa): find a way to log these.
			waitingSpecs.
				each(func(importPath string, waitingList specWaitingList) {
					waitingList.clear()
				}).
				clear()

			// Finally, close the channel.
			close(fetchSHAResultChan)
			fetchSHAResultChan = nil
		}

		// If both of the channels being selected are nil, its time to stop
		// selecting.
		if fetchSHAResultChan == nil &&
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
	if accumulatedErrors.Len() > 0 {
		return accumulatedErrors.Compose("Failed to process dependencies")
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
