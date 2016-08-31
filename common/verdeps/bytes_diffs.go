package verdeps

import (
	"fmt"
	"sort"
)

type bytesDiff struct {
	bytes         []byte
	exclusiveTo   int
	inclusiveFrom int
}

type bytesDiffs []bytesDiff

func (bd bytesDiffs) Len() int {
	return len(bd)
}

func (bd bytesDiffs) Less(i, j int) bool {
	return bd[i].inclusiveFrom >= bd[j].exclusiveTo
}

func (bd bytesDiffs) Swap(i, j int) {
	oldI := bd[i]
	bd[i] = bd[j]
	bd[j] = oldI
}

// sortBytesDiffs returns a last-to-first-sorted slice of bytes diffs.
func sortBytesDiffs(diffs []bytesDiff) []bytesDiff {
	sort.Sort(bytesDiffs(diffs))
	return diffs
}

// sortBytesDiffs returns a last-to-first-sorted slice of bytes diffs.
func calcBytesDiffsDelta(diffs []bytesDiff) int {
	delta := 0
	for _, diff := range diffs {
		delta = delta + (len(diff.bytes) - (diff.exclusiveTo - diff.inclusiveFrom))
	}

	return delta
}

// isBytesDiffValid returns true if diff is valid.
func isBytesDiffValid(diff bytesDiff) bool {
	return diff.exclusiveTo > diff.inclusiveFrom
}

// cursorsInBounds returns true if the bytesCursor and newBytesCursor are both
// within the bounds of slices of size bytesLen and newBytesLen respectively.
func cursorsInBounds(
	bytesLen,
	newBytesLen,
	bytesCursor,
	newBytesCursor int) bool {
	return bytesCursor >= 0 &&
		newBytesCursor >= 0 &&
		bytesCursor <= bytesLen &&
		newBytesCursor <= newBytesLen
}

// composeBytesDiffs takes a collection of diffs, sorts them in reverse order,
// and applies them to a copy of bytes. Afterwards, the bytes copy is returned.
// If any of the diffs are out of bounds, then an error is returned.
func composeBytesDiffs(bytes []byte, diffs []bytesDiff) ([]byte, error) {
	var (
		i                  int
		diff               bytesDiff
		diffsLen           int
		bytesCursor        int
		newBytesCursor     int
		nextBytesCursor    int
		nextNewBytesCursor int
	)

	// If there are no diffs, then bytes is already correct.
	if diffsLen = len(diffs); diffsLen < 1 {
		return bytes, nil
	}

	// Make sure that we have a correctly sorted set of diffs.
	if diffsLen > 1 {
		diffs = sortBytesDiffs(diffs)
	}

	// Perform state initialziation.
	var (
		bytesLen    = len(bytes)
		newBytesLen = len(bytes) + calcBytesDiffsDelta(diffs)
		newBytes    = make([]byte, newBytesLen)
	)

	// Start the cursors at the back.
	bytesCursor = bytesLen
	newBytesCursor = newBytesLen

	// Apply all of the diffs in order.
	for i, diff = range diffs {
		// First, validate the diff.
		if !isBytesDiffValid(diff) {
			return nil, fmt.Errorf("A bytes diff was invalid: %v.", diff)
		}

		// Copy all the stuff in between the previous diff and this one.
		nextBytesCursor = diff.exclusiveTo
		nextNewBytesCursor = newBytesCursor - (bytesCursor - nextBytesCursor)
		// Check that we're in bounds.
		if !cursorsInBounds(bytesLen, newBytesLen, nextBytesCursor, nextNewBytesCursor) {
			return nil, fmt.Errorf("A byte diffs was out of bounds: %v.", diff)
		}
		// Perform the copy.
		copy(newBytes[nextNewBytesCursor:newBytesCursor], bytes[nextBytesCursor:bytesCursor])
		// Advance the cursors.
		bytesCursor = nextBytesCursor
		newBytesCursor = nextNewBytesCursor

		// Use the cursors to continue copying.
		nextBytesCursor = bytesCursor - (diff.exclusiveTo - diff.inclusiveFrom)
		nextNewBytesCursor = newBytesCursor - len(diff.bytes)
		// Check that we're in bounds.
		if !cursorsInBounds(bytesLen, newBytesLen, nextBytesCursor, nextNewBytesCursor) {
			return nil, fmt.Errorf("A byte diffs was out of bounds: %v.", diff)
		}
		// Perform the copy.
		copy(newBytes[nextNewBytesCursor:newBytesCursor], diff.bytes[:])
		// Advance the cursors.
		bytesCursor = nextBytesCursor
		newBytesCursor = nextNewBytesCursor

		// If this is the last diff, then we have to copy the bytes in front of
		// all the diffs. This means the bytes at the _beginning_ of the slice
		// (due to the sort direction).
		if i == (diffsLen - 1) {
			// Perform the copy.
			copy(newBytes[:bytesCursor], bytes[:newBytesCursor])
		}
	}

	return newBytes, nil
}
