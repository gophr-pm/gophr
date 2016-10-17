package verdeps

import (
	"errors"
	"testing"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchSHA(t *testing.T) {
	Convey("Given a github import path", t, func() {
		var (
			importPath         = `"github.com/x/y"`
			packageSHA         = "1234123412341234123412341234123412341234"
			packageRepo        = "b"
			packageAuthor      = "a"
			packageVersionDate = time.Date(2016, time.April, 8, 14, 12, 0, 0, time.UTC)
		)

		Convey("When the importPath is a subpath of this package, then simply copy the parent package SHA", func() {
			var (
				actualOutputSHA        string
				actualOutputImportPath string

				mockGhSvc                = github.NewMockRequestService()
				outputChan               = make(chan *importPathSHA, 1)
				expectedOutputSHA        = packageSHA
				pendingSHARequests       = newSyncedInt()
				expectedOutputImportPath = `"github.com/a/b/c"`
			)

			// Call fetchSHA synchronously for simplicity's sake. This shouldn't lock
			// despite the fact that fetch sha writes to an output channel, since the
			// channel is buffered (for test purposes).
			pendingSHARequests.increment()
			fetchSHA(fetchSHAArgs{
				ghSvc:              mockGhSvc,
				outputChan:         outputChan,
				importPath:         expectedOutputImportPath,
				packageSHA:         packageSHA,
				packageRepo:        packageRepo,
				packageAuthor:      packageAuthor,
				pendingSHARequests: pendingSHARequests,
				packageVersionDate: packageVersionDate,
			})

			// There should be exactly one output sha, so break after reading it.
			for ips := range outputChan {
				actualOutputSHA = ips.sha
				actualOutputImportPath = ips.importPath

				// Close the channel in order to break the loop.
				close(outputChan)
			}

			// Ensure that the github API was not hit with a request.
			mockGhSvc.AssertNotCalled(t, "FetchCommitSHA", "a", "b")

			So(pendingSHARequests.value(), ShouldEqual, 0)
			So(actualOutputSHA, ShouldEqual, expectedOutputSHA)
			So(actualOutputImportPath, ShouldEqual, expectedOutputImportPath)
		})

		Convey("When the SHA request fails, no SHA should be enqueued", func() {
			var (
				mockGhSvc          = github.NewMockRequestService()
				outputChan         = make(chan *importPathSHA, 1)
				pendingSHARequests = newSyncedInt()
			)

			// Expect that fetch commit sha is called.
			mockGhSvc.On(
				"FetchCommitSHA",
				"x",
				"y",
				packageVersionDate,
			).Return("", errors.New("this is an error"))

			// Call fetchSHA synchronously for simplicity's sake. This shouldn't lock
			// despite the fact that fetch sha writes to an output channel, since the
			// channel is buffered (for test purposes).
			pendingSHARequests.increment()
			fetchSHA(fetchSHAArgs{
				ghSvc:              mockGhSvc,
				outputChan:         outputChan,
				importPath:         importPath,
				packageSHA:         packageSHA,
				packageRepo:        packageRepo,
				packageAuthor:      packageAuthor,
				pendingSHARequests: pendingSHARequests,
				packageVersionDate: packageVersionDate,
			})

			// Make sure the output channel gets closed.
			defer close(outputChan)

			So(pendingSHARequests.value(), ShouldEqual, 0)
			// There should be no SHA in the output chan.
			So(len(outputChan), ShouldEqual, 0)
		})

		Convey("When the SHA request returns with an empty SHA, no SHA should be enqueued", func() {
			var (
				mockGhSvc          = github.NewMockRequestService()
				outputChan         = make(chan *importPathSHA, 1)
				pendingSHARequests = newSyncedInt()
			)

			// Expect that fetch commit sha is called.
			mockGhSvc.On(
				"FetchCommitSHA",
				"x",
				"y",
				packageVersionDate,
			).Return("", nil)

			// Call fetchSHA synchronously for simplicity's sake. This shouldn't lock
			// despite the fact that fetch sha writes to an output channel, since the
			// channel is buffered (for test purposes).
			pendingSHARequests.increment()
			fetchSHA(fetchSHAArgs{
				ghSvc:              mockGhSvc,
				outputChan:         outputChan,
				importPath:         importPath,
				packageSHA:         packageSHA,
				packageRepo:        packageRepo,
				packageAuthor:      packageAuthor,
				pendingSHARequests: pendingSHARequests,
				packageVersionDate: packageVersionDate,
			})

			// Make sure the output channel gets closed.
			defer close(outputChan)

			So(pendingSHARequests.value(), ShouldEqual, 0)
			// There should be no SHA in the output chan.
			So(len(outputChan), ShouldEqual, 0)
		})

		Convey("When the SHA request suceeds, the SHA should be enqueued", func() {
			var (
				actualOutputSHA        string
				actualOutputImportPath string

				mockGhSvc                = github.NewMockRequestService()
				outputChan               = make(chan *importPathSHA, 1)
				expectedOutputSHA        = "thisistheoutputshathisistheoutputsha!!!!"
				pendingSHARequests       = newSyncedInt()
				expectedOutputImportPath = importPath
			)

			// Expect that fetch commit sha is called.
			mockGhSvc.On(
				"FetchCommitSHA",
				"x",
				"y",
				packageVersionDate,
			).Return(expectedOutputSHA, nil)

			// Call fetchSHA synchronously for simplicity's sake. This shouldn't lock
			// despite the fact that fetch sha writes to an output channel, since the
			// channel is buffered (for test purposes).
			pendingSHARequests.increment()
			fetchSHA(fetchSHAArgs{
				ghSvc:              mockGhSvc,
				outputChan:         outputChan,
				importPath:         importPath,
				packageSHA:         packageSHA,
				packageRepo:        packageRepo,
				packageAuthor:      packageAuthor,
				pendingSHARequests: pendingSHARequests,
				packageVersionDate: packageVersionDate,
			})

			// There should be exactly one output sha, so break after reading it.
			for ips := range outputChan {
				actualOutputSHA = ips.sha
				actualOutputImportPath = ips.importPath

				// Close the channel in order to break the loop.
				close(outputChan)
			}

			So(pendingSHARequests.value(), ShouldEqual, 0)
			So(actualOutputSHA, ShouldEqual, expectedOutputSHA)
			So(actualOutputImportPath, ShouldEqual, expectedOutputImportPath)
		})
	})
}
