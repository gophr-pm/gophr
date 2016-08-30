package verdeps

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteDiffs(t *testing.T) {
	a := []byte{1, 2, 3}
	b := []byte{4, 5, 6, 7}
	c := []byte{8}
	base := []byte{100, 101, 102}

	_, err := composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         a,
			inclusiveFrom: 1,
			exclusiveTo:   1,
		},
	})
	assert.NotNil(t, err, "should fail on invalid diffs (to=from)")

	_, err = composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         a,
			inclusiveFrom: 2,
			exclusiveTo:   1,
		},
	})
	assert.NotNil(t, err, "should fail on invalid diffs (to<from)")

	expectedOutput := []byte{100, 1, 2, 3, 102}
	actualOutput, err := composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         a,
			inclusiveFrom: 1,
			exclusiveTo:   2,
		},
	})
	assert.Nil(t, err, "should be no error")
	assert.True(t, bytesEqual(expectedOutput, actualOutput), "single index subs should work")

	expectedOutput = []byte{1, 2, 3, 102}
	actualOutput, err = composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         a,
			inclusiveFrom: 0,
			exclusiveTo:   2,
		},
	})
	assert.Nil(t, err, "should be no error")
	assert.True(t, bytesEqual(expectedOutput, actualOutput), "multi index subs should work")

	expectedOutput = []byte{1, 2, 3}
	actualOutput, err = composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         a,
			inclusiveFrom: 0,
			exclusiveTo:   3,
		},
	})
	assert.Nil(t, err, "should be no error")
	assert.True(t, bytesEqual(expectedOutput, actualOutput), "whole slice subs should work")

	expectedOutput = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	actualOutput, err = composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         a,
			inclusiveFrom: 0,
			exclusiveTo:   1,
		},
		{
			bytes:         b,
			inclusiveFrom: 1,
			exclusiveTo:   2,
		},
		{
			bytes:         c,
			inclusiveFrom: 2,
			exclusiveTo:   3,
		},
	})
	assert.Nil(t, err, "should be no error")
	assert.True(t, bytesEqual(expectedOutput, actualOutput), "subsequence slice subs should work")

	expectedOutput = []byte{1, 2, 3, 4, 5, 6, 7, 8}
	actualOutput, err = composeBytesDiffs(base, []bytesDiff{
		{
			bytes:         c,
			inclusiveFrom: 2,
			exclusiveTo:   3,
		},
		{
			bytes:         b,
			inclusiveFrom: 1,
			exclusiveTo:   2,
		},
		{
			bytes:         a,
			inclusiveFrom: 0,
			exclusiveTo:   1,
		},
	})
	assert.Nil(t, err, "should be no error")
	assert.True(t, bytesEqual(expectedOutput, actualOutput), "subsequence slice subs should work out of order")
}

func bytesEqual(a, b []byte) bool {
	if a == nil && b == nil {
		// fmt.Printf("%v equals %v.\n", a, b)
		return true
	}

	if a == nil || b == nil {
		// fmt.Printf("%v does not equal %v.\n", a, b)
		return false
	}

	if len(a) != len(b) {
		// fmt.Printf("%v does not equal %v.\n", a, b)
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			// fmt.Printf("%v does not equal %v.\n", a, b)
			return false
		}
	}

	// fmt.Printf("%v equals %v.\n", a, b)
	return true
}
