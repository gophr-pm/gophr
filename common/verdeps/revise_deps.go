package verdeps

import (
	"errors"
	"io/ioutil"
	"log"
	"sync"
)

const charDoubleQuote = '"'

type reviseDepsArgs struct {
	inputChan          chan *revision
	revisionWaitGroup  *sync.WaitGroup
	accumulatedErrors  *syncedErrors
	syncedImportCounts *syncedImportCounts
}

func reviseDeps(args reviseDepsArgs) {
	var (
		path             string
		revs             []*revision
		importCount      int
		pathRevisionsMap = make(map[string][]*revision)
	)

	// Take care of our wait group responsibilities first and foremost.
	args.revisionWaitGroup.Add(1)
	defer args.revisionWaitGroup.Done()

	// Process every revision that comes in from the revision channel.
	for rev := range args.inputChan {
		// Get the rev slice, and add this rev.
		path = rev.path
		revs = pathRevisionsMap[path]
		revs = append(revs, rev)
		importCount = args.syncedImportCounts.importCountOf(path)

		// Decide whether its time to apply the revs.
		if importCount == len(revs) {
			// Apply the revisions now that we have all the appropriate revisions.
			if len(revs) > 0 {
				go applyRevisions(path, revs, args.revisionWaitGroup, args.accumulatedErrors)
			}
			// Get rids of the revs from the map since we don't need them anymore.
			delete(pathRevisionsMap, path)
		} else {
			// Update the revs in the map - keep waiting for mroe revs.
			pathRevisionsMap[path] = revs
		}
	}

	var (
		missedImports           = 0
		filesWithMissingImports = len(pathRevisionsMap)
	)

	// Apply all remaining revisions, and log the files that don't have every
	// import versioned.
	for path, revs = range pathRevisionsMap {
		// Record how many imports we missed
		missedImports = missedImports + (args.syncedImportCounts.importCountOf(path) - len(revs))
		// Apply the revisions that we have (given we have any).
		if len(revs) > 0 {
			go applyRevisions(path, revs, args.revisionWaitGroup, args.accumulatedErrors)
		}
		// Get rids of the revs from the map since we don't need them anymore.
		delete(pathRevisionsMap, path)
	}

	// Summarize what we missed in a log message.
	if missedImports > 0 {
		log.Printf("Missed %d imports in %d files.\n", missedImports, filesWithMissingImports)
	}
}

// applyRevisions applies all the provided revisions to the appropriate files.
func applyRevisions(
	path string,
	revs []*revision,
	waitGroup *sync.WaitGroup,
	accumulatedErrors *syncedErrors) {
	var (
		err      error
		diffs    []bytesDiff
		from, to int
		fileData []byte
	)

	// Take care of our wait group responsibilities first and foremost.
	waitGroup.Add(1)
	defer waitGroup.Done()

	// Read the file data at the specified path.
	if fileData, err = ioutil.ReadFile(path); err != nil {
		accumulatedErrors.add(err)
		return
	}

	// Create bytes diffs for each of the revisions.
	for _, rev := range revs {
		// Adjust from and to so that they fall on quote bytes.
		if from, to, err = findImportPathBoundaries(fileData, rev.fromIndex, rev.toIndex); err != nil {
			// Exit if the import path boundaries could not be adjusted.
			accumulatedErrors.add(err)
			return
		}

		diffs = append(diffs, bytesDiff{
			bytes:         rev.gophrURL,
			exclusiveTo:   to,
			inclusiveFrom: from,
		})
	}

	// Combine the diffs and the file data.
	if fileData, err = composeBytesDiffs(fileData, diffs); err != nil {
		accumulatedErrors.add(err)
		return
	}

	// After the file data has been adequately tampered with. Write back to the
	// file.
	if err = ioutil.WriteFile(path, fileData, 0644); err != nil {
		accumulatedErrors.add(err)
		return
	}
}

// findImportPathBoundaries adjusts from and to to align perfectly with a
// quoted import path. If the import path cannot be found, then an error is
// returned.
func findImportPathBoundaries(data []byte, from, to int) (int, int, error) {
	var (
		i            int
		adjustedTo   int
		adjustedFrom int
	)

	// Firstly, read backwards until we hit a quote on the left.
	for i = from + 2; isInBounds(data, i) && data[i] != charDoubleQuote; i-- {
	}

	// Exit if out of bounds.
	if !isInBounds(data, i) {
		return -1, -1, errors.New("Could not find the beginning of the import path")
	}

	// We now have the adjusted from.
	adjustedFrom = i

	// Last, read forwards until we hit a quote on the right.
	for i = to - 2; isInBounds(data, i) && data[i] != charDoubleQuote; i++ {
	}

	// Exit if out of bounds.
	if !isInBounds(data, i) || !isInBounds(data, i+1) {
		return -1, -1, errors.New("Could not find the end of the import path")
	}

	// We now have the adjusted to.
	adjustedTo = i + 1

	return adjustedFrom, adjustedTo, nil
}

// isInBounds returns true if i is an index of data.
func isInBounds(data []byte, i int) bool {
	return i >= 0 && i < len(data)
}
