package verdeps

import (
	"errors"
	"io/ioutil"
	"sort"
	"sync"
)

const charDoubleQuote = '"'

type reviseDepsArgs struct {
	inputChan         chan *revision
	revisionWaitGroup *sync.WaitGroup
	accumulatedErrors *syncedErrors
}

func reviseDeps(args reviseDepsArgs) {
	var (
		lastPath       string
		revisionBuffer []*revision
	)

	// Take care of our wait group responsibilities first and foremost.
	args.revisionWaitGroup.Add(1)
	defer args.revisionWaitGroup.Done()

	// Process every revision that comes in from the revision channel.
	for rev := range args.inputChan {
		// If the last path is different from the current one, flush the buffer.
		if len(lastPath) > 0 && rev.path != lastPath {
			// Group all the revisions that apply to the same file.
			go applyRevisions(lastPath, revisionBuffer, args.revisionWaitGroup, args.accumulatedErrors)
			// Clear the buffer again.
			revisionBuffer = nil
		}

		// Adjust the last path to keep it accurate.
		lastPath = rev.path
		// Accumulate co-located revisions.
		revisionBuffer = append(revisionBuffer, rev)
	}

	// Take care of the remaining revisions.
	go applyRevisions(lastPath, revisionBuffer, args.revisionWaitGroup, args.accumulatedErrors)
	revisionBuffer = nil
}

// applyRevisions applies all the provided revisions to the appropriate files.
func applyRevisions(
	path string,
	revs []*revision,
	waitGroup *sync.WaitGroup,
	accumulatedErrors *syncedErrors) {
	var (
		err      error
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

	// Sort the revs so that the last revs come first. If we do it in this order,
	// then offsets don't have to be managed.
	sort.Sort(sortableRevisions(revs))

	// Iterate through each revision, applying changes from last index to first
	// index.
	for _, rev := range revs {
		// Find the exact boundaries.
		from, to, err = findImportPathBoundaries(fileData, rev.fromIndex, rev.toIndex)
		if err != nil {
			accumulatedErrors.add(err)
			continue
		}

		// Perform the file data changes.
		fileData = embedByteSlice(fileData, rev.gophrURL, from, to)
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

// embedByteSlice replaces the bytes of outer from within the specifies indicies
// with the bytes from inner.
func embedByteSlice(outer, inner []byte, from, to int) []byte {
	delta := len(inner) - (to - from)
	if delta > 0 {
		oldOuterLen := len(outer)
		outerExtension := make([]byte, delta)
		newOuter := append(outer, outerExtension...)
		copy(newOuter[from+delta:oldOuterLen+delta], outer[from:oldOuterLen])
		copy(newOuter[from:from+len(inner)], inner[:])
		return newOuter
	} else if delta < 0 {
		newOuter := make([]byte, len(outer)+delta)
		copy(newOuter[:from], outer[:from])
		copy(newOuter[from:from+len(inner)], inner[:])
		copy(newOuter[from+len(inner):], outer[to:])
		return newOuter
	} else {
		copy(outer[from:to], inner[:])
		return outer
	}
}
