package verdeps

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/skeswa/gophr/common/github"
)

type processDepsArgs struct {
	ghSvc              *github.RequestService
	packageSHA         string
	packagePath        string
	packageRepo        string
	packageAuthor      string
	packageVersionDate time.Time
}

func processDeps(args processDepsArgs) error {
	var (
		revisionChan       = make(chan *revision)
		waitingSpecs       = make(map[string]*specWaitingList)
		importPathSHAs     = make(map[string]string)
		importSpecChan     = make(chan *importSpec)
		importPathSHAChan  = make(chan *importPathSHA)
		accumulatedErrors  = newSyncedErrors()
		revisionWaitGroup  = &sync.WaitGroup{}
		pendingSHARequests = newSyncedInt()
	)

	// Start reading the deps.
	go readDeps(readDepsArgs{
		outputChan:        importSpecChan,
		packagePath:       args.packagePath,
		accumulatedErrors: accumulatedErrors,
	})

	// Revise dependencies in the go source files.
	go reviseDeps(reviseDepsArgs{
		inputChan:         revisionChan,
		revisionWaitGroup: revisionWaitGroup,
		accumulatedErrors: accumulatedErrors,
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
			importPathSHAs[importPathHash] = ips.sha

			// Clear away the waiting specs.
			if waitingList, exists := waitingSpecs[importPathHash]; exists {
				// There is a waiting list, so it needs to be cleared.
				if specs := waitingList.clear(); specs != nil {
					for _, spec := range specs {
						enqueueRevision(revisionChan, spec.imports.Path.Value, ips.sha, spec)
					}
				}
			}

			// Check if this is the last time tme an import path sha will come through.
			// We check this by checking if the spec channel has already closed, and
			// if there are no pending sha requests. If so, close this channel.
			if pendingSHARequests.value() == 0 && importSpecChan == nil {
				closeImportPathSHAChan(importPathSHAChan, waitingSpecs)
			}

		case spec, alive := <-importSpecChan:
			// Nil the channel if no longer alive.
			if !alive {
				importSpecChan = nil

				// If there are no pending sha requests, then there never will be since
				// this channel is closing. Therefore, we can close the sha channel.
				if shouldCloseImportPathSHAChan(pendingSHARequests, importSpecChan, importPathSHAChan) {
					closeImportPathSHAChan(importPathSHAChan, waitingSpecs)
				}

				break
			}

			// For each incoming spec, make it wait keyed on the import path hash.
			importPath := spec.imports.Path.Value
			importPathHash := importPathHashOf(spec.imports.Path.Value)
			if sha, exists := importPathSHAs[importPathHash]; !exists {
				// If we don't presently have the sha, then we have to go out and get it.
				if specs, exists := waitingSpecs[importPathHash]; !exists {
					waitingSpecs[importPathHash] = newSpecWaitingList(spec)
					go fetchSHA(fetchSHAArgs{
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
					if ok := specs.add(spec); !ok {
						// If the add failed, assume that it is because the the sha was
						// obtained after we last checked.
						if sha, exists = importPathSHAs[importPathHash]; !exists {
							accumulatedErrors.add(fmt.Errorf("Could not version dependency %s because the SHA did not yet exist.", importPath))
						} else {
							enqueueRevision(revisionChan, importPath, sha, spec)
						}
					}
				}
			} else {
				// If we got here, it means that the sha has already been obtained, so the
				// new import path exists.
				enqueueRevision(revisionChan, importPath, sha, spec)
			}
		}

		// If both of the channels being selected are nil, its time to stop selecting.
		if importPathSHAChan == nil && importSpecChan == nil {
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

// enqueueRevision is a helper function that puts a revision into the revision channel.
func enqueueRevision(revisionChan chan *revision, importPath, sha string, spec *importSpec) {
	author, repo, subpath := parseImportPath(importPath)
	newImportPath := composeNewImportPath(author, repo, sha, subpath)

	revisionChan <- newRevision(spec, newImportPath)
}

// closeImportPathSHAChan closes the importPathSHAChan and clears all spec waiting list.
// The waiting lists are cleared because they're waiting for a SHA that will
// presumably never come.
func closeImportPathSHAChan(importPathSHAChan chan *importPathSHA, importPathSHAs map[string]*specWaitingList) {
	// Clear all the waiting lists.
	for importPath, waitingList := range importPathSHAs {
		waitingList.clear()
		// Remove this key-value pair from the map.
		delete(importPathSHAs, importPath)
	}

	// Finally, close the channel.
	close(importPathSHAChan)
}

// shouldCloseImportPathSHAChan returns true if the importPathSHAChan should be
// closed.
func shouldCloseImportPathSHAChan(
	pendingSHARequests *syncedInt,
	importSpecChan chan *importSpec,
	importPathSHAChan chan *importPathSHA) bool {
	return pendingSHARequests.value() == 0 && importSpecChan == nil && importPathSHAChan != nil
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
