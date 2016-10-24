package verdeps

import (
	"errors"
	"testing"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
	. "github.com/smartystreets/goconvey/convey"
)

func TestVersionDeps(t *testing.T) {
	Convey("Given package metadata and a package directory", t, func() {
		Convey("If an invalid IO is provided, an error should be returned", func() {
			err := VersionDeps(VersionDepsArgs{
				IO:     nil,
				SHA:    "1234123412341234123412341234123412341234",
				Repo:   "myrepo",
				Path:   "/a/b/c",
				Author: "myauthor",
				processDeps: func(args processDepsArgs) error {
					return nil
				},
				GithubService: github.NewMockRequestService(),
			})

			So(err, ShouldNotBeNil)
		})

		Convey("If an invalid SHA is provided, an error should be returned", func() {
			err := VersionDeps(VersionDepsArgs{
				IO:     io.NewMockIO(),
				SHA:    "",
				Repo:   "myrepo",
				Path:   "/a/b/c",
				Author: "myauthor",
				processDeps: func(args processDepsArgs) error {
					return nil
				},
				GithubService: github.NewMockRequestService(),
			})

			So(err, ShouldNotBeNil)
		})

		Convey("If an invalid path is provided, an error should be returned", func() {
			err := VersionDeps(VersionDepsArgs{
				IO:     io.NewMockIO(),
				SHA:    "1234123412341234123412341234123412341234",
				Repo:   "myrepo",
				Path:   "",
				Author: "myauthor",
				processDeps: func(args processDepsArgs) error {
					return nil
				},
				GithubService: github.NewMockRequestService(),
			})

			So(err, ShouldNotBeNil)
		})

		Convey("If an invalid repo is provided, an error should be returned", func() {
			err := VersionDeps(VersionDepsArgs{
				IO:     io.NewMockIO(),
				SHA:    "1234123412341234123412341234123412341234",
				Repo:   "",
				Path:   "/a/b/c",
				Author: "myauthor",
				processDeps: func(args processDepsArgs) error {
					return nil
				},
				GithubService: github.NewMockRequestService(),
			})

			So(err, ShouldNotBeNil)
		})

		Convey("If an invalid author is provided, an error should be returned", func() {
			err := VersionDeps(VersionDepsArgs{
				IO:     io.NewMockIO(),
				SHA:    "1234123412341234123412341234123412341234",
				Repo:   "myrepo",
				Path:   "/a/b/c",
				Author: "",
				processDeps: func(args processDepsArgs) error {
					return nil
				},
				GithubService: github.NewMockRequestService(),
			})

			So(err, ShouldNotBeNil)
		})

		Convey("If an invalid github service is provided, an error should be returned", func() {
			err := VersionDeps(VersionDepsArgs{
				IO:     io.NewMockIO(),
				SHA:    "1234123412341234123412341234123412341234",
				Repo:   "myrepo",
				Path:   "/a/b/c",
				Author: "myauthor",
				processDeps: func(args processDepsArgs) error {
					return nil
				},
				GithubService: nil,
			})

			So(err, ShouldNotBeNil)
		})

		Convey("If the commit timestamp cannot be fetched, an error should be returned", func() {
			ghSvc := github.NewMockRequestService()
			ghSvc.On(
				"FetchCommitTimestamp",
				"myauthor",
				"myrepo",
				"1234123412341234123412341234123412341234").
				Return(time.Now(), errors.New("this is an error"))

			err := VersionDeps(VersionDepsArgs{
				IO:            io.NewMockIO(),
				SHA:           "1234123412341234123412341234123412341234",
				Repo:          "myrepo",
				Path:          "/a/b/c",
				Author:        "myauthor",
				processDeps:   nil,
				GithubService: ghSvc,
			})

			So(err, ShouldNotBeNil)
			ghSvc.AssertExpectations(t)
		})

		Convey("If the commit timestamp is fetched, processDeps should be invoked", func() {
			var (
				ghSvc                 = github.NewMockRequestService()
				testTime              = time.Now()
				actualProcessDepsArgs processDepsArgs
			)

			ghSvc.On(
				"FetchCommitTimestamp",
				"myauthor",
				"myrepo",
				"1234123412341234123412341234123412341234").
				Return(testTime, nil)

			err := VersionDeps(VersionDepsArgs{
				IO:     io.NewMockIO(),
				SHA:    "1234123412341234123412341234123412341234",
				Repo:   "myrepo",
				Path:   "/a/b/c",
				Author: "myauthor",
				processDeps: func(args processDepsArgs) error {
					actualProcessDepsArgs = args
					return nil
				},
				GithubService: ghSvc,
			})

			So(err, ShouldBeNil)
			So(actualProcessDepsArgs.fetchSHA, ShouldNotBeNil)
			So(actualProcessDepsArgs.ghSvc, ShouldNotBeNil)
			So(actualProcessDepsArgs.io, ShouldNotBeNil)
			So(actualProcessDepsArgs.newSpecWaitingList, ShouldNotBeNil)
			So(actualProcessDepsArgs.newSyncedStringMap, ShouldNotBeNil)
			So(actualProcessDepsArgs.newSyncedWaitingListMap, ShouldNotBeNil)
			So(actualProcessDepsArgs.packageAuthor, ShouldEqual, "myauthor")
			So(actualProcessDepsArgs.packagePath, ShouldEqual, "/a/b/c")
			So(actualProcessDepsArgs.packageRepo, ShouldEqual, "myrepo")
			So(
				actualProcessDepsArgs.packageSHA,
				ShouldEqual,
				"1234123412341234123412341234123412341234")
			So(actualProcessDepsArgs.packageVersionDate, ShouldResemble, testTime)
			So(actualProcessDepsArgs.readPackageDir, ShouldNotBeNil)
			So(actualProcessDepsArgs.reviseDeps, ShouldNotBeNil)
			ghSvc.AssertExpectations(t)
		})
	})
}
