package verdeps

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestByteDiffs__Basic(t *testing.T) {
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

const (
	textInputStr1 = `package main

import "github.com/a/b@x" // lkfdjs
// sdfkl
import (
"fiॺॺॺॺne"
"arॺe"
"yॺॺou"
"good?"



_ "github.com/how/are/you@doing"
"github.com/how/are/you@doing"
"github.com/how/are/you@doing"
)

this is some stuff
NOW ITS OVER
`
	textOutputStr1 = `package main

import REPLACED // lkfdjs
// sdfkl
import (
"fiॺॺॺॺne"
"arॺe"
"yॺॺou"
"good?"



_ REPLACED
REPLACED
REPLACED
)

this is some stuff
NOW ITS OVER
`
	textReplacement         = "REPLACED"
	textReplaceTarget1      = `"github.com/a/b@x"`
	textReplaceTarget2      = `"github.com/how/are/you@doing"`
	textReplaceTargetIndex1 = 21
	textReplaceTargetIndex2 = 121
)

var (
	target1StartIndex    = strings.Index(textInputStr1, textReplaceTarget1)
	target1EndIndex      = target1StartIndex + len(textReplaceTarget1)
	target2StartIndex    = strings.Index(textInputStr1, textReplaceTarget2)
	target2EndIndex      = target2StartIndex + len(textReplaceTarget2)
	target3StartIndex    = target2StartIndex + len(textReplaceTarget2) + 1
	target3EndIndex      = target3StartIndex + len(textReplaceTarget2)
	target4StartIndex    = target3StartIndex + len(textReplaceTarget2) + 1
	target4EndIndex      = target4StartIndex + len(textReplaceTarget2)
	textInputBytes1      = []byte(textInputStr1)
	textReplacementBytes = []byte(textReplacement)
)

func TestByteDiffs__Text(t *testing.T) {
	outputBytes, err := composeBytesDiffs(textInputBytes1, []bytesDiff{
		{
			bytes:         textReplacementBytes,
			inclusiveFrom: target1StartIndex,
			exclusiveTo:   target1EndIndex,
		},
		{
			bytes:         textReplacementBytes,
			inclusiveFrom: target2StartIndex,
			exclusiveTo:   target2EndIndex,
		},
		{
			bytes:         textReplacementBytes,
			inclusiveFrom: target3StartIndex,
			exclusiveTo:   target3EndIndex,
		},
		{
			bytes:         textReplacementBytes,
			inclusiveFrom: target4StartIndex,
			exclusiveTo:   target4EndIndex,
		},
	})
	assert.Nil(t, err, "should be no error")
	assert.Equal(t, textOutputStr1, string(outputBytes[:]), "text replacement should work")
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
