package verdeps

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

var (
	fileDataWithComment = []byte(`
/**
* Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam ex orci,
* cursus et vehicula eget, condimentum in mauris. Nam lacinia, turpis eget
* volutpat pellentesque, tortor dolor gravida nisl, vel pellentesque felis
* purus quis est. Class aptent taciti sociosqu ad litora
* torquent per conubia nostra, per inceptos himenaeos. */

package thingy // import "github.com/a/thingy"

this()
is()
some()
other()
stuff()
`)
	fileDataWithoutComment = []byte(`
/**
* Lorem ipsum dolor sit amet, consectetur adipiscing elit. Nullam ex orci,
* cursus et vehicula eget, condimentum in mauris. Nam lacinia, turpis eget
* volutpat pellentesque, tortor dolor gravida nisl, vel pellentesque felis
* purus quis est. Class aptent taciti sociosqu ad litora
* torquent per conubia nostra, per inceptos himenaeos. */

package thingy

this()
is()
some()
other()
stuff()
`)
)

func TestFindPackageImportComment(t *testing.T) {
	Convey("Given file data and a start index", t, func() {
		Convey("When there is a package import comment, its indices should be returned", func() {
			var (
				fileData           = fileDataWithComment
				expectedEndIndex   = 392
				expectedStartIndex = 360
			)

			actualStartIndex, actualEndIndex := findPackageImportComment(
				fileData,
				351)

			So(actualEndIndex, ShouldEqual, expectedEndIndex)
			So(actualStartIndex, ShouldEqual, expectedStartIndex)

			Convey("If the package start index is after the package import comment, then -1 should be returned", func() {
				expectedEndIndex = -1
				expectedStartIndex = -1
				actualStartIndex, actualEndIndex = findPackageImportComment(fileData, len(fileData)-1)

				So(actualEndIndex, ShouldEqual, expectedEndIndex)
				So(actualStartIndex, ShouldEqual, expectedStartIndex)
			})
		})

		Convey("When there is not a package import comment, -1 should be returned", func() {
			var (
				fileData           = fileDataWithoutComment
				expectedEndIndex   = -1
				expectedStartIndex = -1
			)

			actualStartIndex, actualEndIndex := findPackageImportComment(
				fileData,
				351)

			So(actualEndIndex, ShouldEqual, expectedEndIndex)
			So(actualStartIndex, ShouldEqual, expectedStartIndex)
		})
	})
}
