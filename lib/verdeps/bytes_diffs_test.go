package verdeps

import (
	"strings"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestByteDiffs(t *testing.T) {
	Convey("Given zero diffs and a simple base byte slice", t, func() {
		var base = []byte{100, 101, 102}

		Convey("Should do nothing, since is nothing to change", func() {
			expectedOutput := base
			actualOutput, err := composeBytesDiffs(base, []bytesDiff{})

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, expectedOutput)
		})
	})

	Convey("Given one diff and a simple base byte slice", t, func() {
		var (
			a    = []byte{1, 2, 3}
			base = []byte{100, 101, 102}
		)

		Convey("Diffs with out-of-bounds indices should cause a failure", func() {
			_, err := composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: 12378,
					exclusiveTo:   12379,
				},
			})

			So(err, ShouldNotBeNil)

			_, err = composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: -2,
					exclusiveTo:   1,
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("Invalid diffs (where to = from) should cause a failure", func() {
			_, err := composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: 1,
					exclusiveTo:   1,
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("Invalid diffs (where to < from) should cause a failure", func() {
			_, err := composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: 2,
					exclusiveTo:   1,
				},
			})

			So(err, ShouldNotBeNil)
		})

		Convey("Single-indexed sub-slices should work", func() {
			expectedOutput := []byte{100, 1, 2, 3, 102}
			actualOutput, err := composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: 1,
					exclusiveTo:   2,
				},
			})

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, expectedOutput)
		})

		Convey("Multi-indexed sub-slices should work", func() {
			expectedOutput := []byte{1, 2, 3, 102}
			actualOutput, err := composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: 0,
					exclusiveTo:   2,
				},
			})

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, expectedOutput)
		})

		Convey("All-encompassing sub-slices should work", func() {
			expectedOutput := []byte{1, 2, 3}
			actualOutput, err := composeBytesDiffs(base, []bytesDiff{
				{
					bytes:         a,
					inclusiveFrom: 0,
					exclusiveTo:   3,
				},
			})

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, expectedOutput)
		})
	})

	Convey("Given multiple diffs and a simple base byte slice", t, func() {
		var (
			a    = []byte{1, 2, 3}
			b    = []byte{4, 5, 6, 7}
			c    = []byte{8}
			base = []byte{100, 101, 102}
		)

		Convey("Sub-sequence sub-slices should work", func() {
			expectedOutput := []byte{1, 2, 3, 4, 5, 6, 7, 8}
			actualOutput, err := composeBytesDiffs(base, []bytesDiff{
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

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, expectedOutput)
		})

		Convey("Sub-sequence sub-slices should work out of order", func() {
			expectedOutput := []byte{1, 2, 3, 4, 5, 6, 7, 8}
			actualOutput, err := composeBytesDiffs(base, []bytesDiff{
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

			So(err, ShouldBeNil)
			So(actualOutput, ShouldResemble, expectedOutput)
		})
	})

	Convey("Given multiple text diffs and base text", t, func() {
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

		Convey("Full text replacement should function as expected", func() {
			expectedOutput := textOutputStr1
			actualOutputBytes, err := composeBytesDiffs(textInputBytes1, []bytesDiff{
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

			So(err, ShouldBeNil)
			So(string(actualOutputBytes[:]), ShouldEqual, expectedOutput)
		})
	})
}
