package verdeps

import (
	"fmt"
	"go/ast"
	"go/token"
	"math/rand"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessDeps(t *testing.T) {
	Convey("Given go package with dependencies", t, func() {
		Convey("In the canonical case, processDeps should always succeed", func() {
			var (
				fetchSHA                 shaFetcher
				reviseDeps               depsReviser
				readPackageDir           packageDirReader
				shaRequestsLock          sync.RWMutex
				actualReviseDepsArgs     reviseDepsArgs
				actualReadPackageDirArgs readPackageDirArgs

				io                 = io.NewMockIO()
				ghSvc              = github.NewMockRequestService()
				packageSHA         = "thisisthepackageshathisisthepackagesha!!"
				packagePath        = "sompackagepath"
				packageRepo        = "somerepo"
				shaRequests        = make(map[string]*fetchSHAArgs)
				packageAuthor      = "someauthor"
				importRevStrings   = make(map[string]bool)
				packageRevStrings  = make(map[string]bool)
				packageVersionDate = time.Date(
					2016,
					time.April,
					8,
					14,
					12,
					0,
					0,
					time.UTC)
			)

			// Create fakes of the worker functions passed into processDeps.
			fetchSHA = func(args fetchSHAArgs) {
				// Record this invocation of fetchSHA in a synchronized map.
				shaRequestsLock.Lock()
				shaRequests[args.importPath] = &args
				shaRequestsLock.Unlock()

				// After a pause, signal that this SHA was fetched.
				introduceRandomLag(0.6, 30)
				args.pendingSHARequests.decrement()

				args.outputChan <- &importPathSHA{
					sha:        "thisistheshafor" + args.importPath,
					importPath: args.importPath,
				}
			}
			reviseDeps = func(args reviseDepsArgs) {
				// Stow the received args for later assertions.
				actualReviseDepsArgs = args

				// Read the input revisions, and record what comes through.
				for input := range args.inputChan {
					if input.revisesImport {
						importRevStrings[stringifyRevision(input)] = true
					} else if input.revisesPackage {
						packageRevStrings[stringifyRevision(input)] = true
					}
				}

				// Signal that revise deps is over.
				args.revisionWaitGroup.Done()
			}
			readPackageDir = func(args readPackageDirArgs) {
				// Stow the received args for later assertions.
				actualReadPackageDirArgs = args

				args.importCounts.setImportCount("filepath1", 3)
				introduceRandomLag(0.4, 15)
				args.importCounts.setImportCount("filepath2", 2)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(101, "filepath1", `"github.com/c/d/e"`)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(102, "filepath2", `"github.com/f/g"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath2", 2)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(103, "filepath1", `"github.com/b/c"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath1", 1)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(104, "filepath1", `"github.com/a/b"`)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(105, "filepath2", `"github.com/h/i/j/k"`)
				introduceRandomLag(0.4, 15)
				args.importCounts.setImportCount("filepath3", 1)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(106, "filepath3", `"github.com/b/c"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath3", 3)

				// Close both channels once we're done.
				close(args.importSpecChan)
				close(args.packageSpecChan)
			}

			// Execute synchronously to make life easier.
			err := processDeps(processDepsArgs{
				io:                 io,
				ghSvc:              ghSvc,
				fetchSHA:           fetchSHA,
				reviseDeps:         reviseDeps,
				packageSHA:         packageSHA,
				packagePath:        packagePath,
				packageRepo:        packageRepo,
				packageAuthor:      packageAuthor,
				readPackageDir:     readPackageDir,
				packageVersionDate: packageVersionDate,
				newSpecWaitingList: newSpecWaitingList,
				newSyncedStringMap: newSyncedStringMap,
			})

			// Assert up a storm starting with fetchSHA.
			So(len(shaRequests), ShouldEqual, 5)
			So(shaRequests[`"github.com/c/d/e"`], ShouldNotBeNil)
			So(shaRequests[`"github.com/f/g"`], ShouldNotBeNil)
			So(shaRequests[`"github.com/b/c"`], ShouldNotBeNil)
			So(shaRequests[`"github.com/a/b"`], ShouldNotBeNil)
			So(shaRequests[`"github.com/h/i/j/k"`], ShouldNotBeNil)
			for _, args := range shaRequests {
				So(args.ghSvc, ShouldNotBeNil)
				So(args.outputChan, ShouldNotBeNil)
				So(args.importPath, ShouldStartWith, `"github.com/`)
				So(args.packageSHA, ShouldEqual, packageSHA)
				So(args.packageRepo, ShouldEqual, packageRepo)
				So(args.packageAuthor, ShouldEqual, packageAuthor)
				So(args.pendingSHARequests, ShouldNotBeNil)
				So(args.packageVersionDate, ShouldResemble, packageVersionDate)
			}

			// Now make the asserts for revise deps.
			So(actualReviseDepsArgs.accumulatedErrors, ShouldNotBeNil)
			So(actualReviseDepsArgs.inputChan, ShouldNotBeNil)
			So(actualReviseDepsArgs.revisionWaitGroup, ShouldNotBeNil)
			So(actualReviseDepsArgs.syncedImportCounts, ShouldNotBeNil)
			So(importRevStrings[`filepath1:101:119:"gophr.pm/c/d@thisistheshafor"github.com/c/d/e"/e"`], ShouldBeTrue)
			So(importRevStrings[`filepath2:102:118:"gophr.pm/f/g@thisistheshafor"github.com/f/g""`], ShouldBeTrue)
			So(importRevStrings[`filepath1:103:119:"gophr.pm/b/c@thisistheshafor"github.com/b/c""`], ShouldBeTrue)
			So(importRevStrings[`filepath1:104:120:"gophr.pm/a/b@thisistheshafor"github.com/a/b""`], ShouldBeTrue)
			So(importRevStrings[`filepath2:105:125:"gophr.pm/h/i@thisistheshafor"github.com/h/i/j/k"/j/k"`], ShouldBeTrue)
			So(importRevStrings[`filepath3:106:122:"gophr.pm/b/c@thisistheshafor"github.com/b/c""`], ShouldBeTrue)

			// Time for the readPackageDir asserts.
			So(actualReadPackageDirArgs.errors, ShouldNotBeNil)
			So(actualReadPackageDirArgs.generatedInternalDirName, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.importCounts, ShouldNotBeNil)
			So(actualReadPackageDirArgs.importSpecChan, ShouldNotBeNil)
			So(actualReadPackageDirArgs.io, ShouldNotBeNil)
			So(actualReadPackageDirArgs.packageDirPath, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.packageSpecChan, ShouldNotBeNil)

			// Finally, perform asserts for general outputs.
			So(err, ShouldBeNil)
		})

		Convey("Multiple imports requiring the same SHA should only issue one request", func() {
			var (
				fetchSHA                 shaFetcher
				reviseDeps               depsReviser
				readPackageDir           packageDirReader
				allFetchSHAArgs          []fetchSHAArgs
				allFetchSHAArgsLock      sync.RWMutex
				actualReviseDepsArgs     reviseDepsArgs
				actualReadPackageDirArgs readPackageDirArgs

				io                 = io.NewMockIO()
				ghSvc              = github.NewMockRequestService()
				packageSHA         = "thisisthepackageshathisisthepackagesha!!"
				packagePath        = "sompackagepath"
				packageRepo        = "somerepo"
				packageAuthor      = "someauthor"
				importRevStrings   = make(map[string]bool)
				packageRevStrings  = make(map[string]bool)
				packageVersionDate = time.Date(
					2016,
					time.April,
					8,
					14,
					12,
					0,
					0,
					time.UTC)
			)

			// Create fakes of the worker functions passed into processDeps.
			fetchSHA = func(args fetchSHAArgs) {
				// Record this invocation of fetchSHA in a synchronized map.
				allFetchSHAArgsLock.Lock()
				allFetchSHAArgs = append(allFetchSHAArgs, args)
				allFetchSHAArgsLock.Unlock()

				// After a pause, signal that this SHA was fetched.
				if strings.HasPrefix(args.importPath, `"github.com/a/b"`) {
					introduceRandomLag(1.0, 120)
				} else {
					introduceRandomLag(0.5, 30)
				}
				args.pendingSHARequests.decrement()

				args.outputChan <- &importPathSHA{
					sha:        "thisistheshafor" + args.importPath,
					importPath: args.importPath,
				}
			}
			reviseDeps = func(args reviseDepsArgs) {
				// Stow the received args for later assertions.
				actualReviseDepsArgs = args

				// Read the input revisions, and record what comes through.
				for input := range args.inputChan {
					if input.revisesImport {
						importRevStrings[stringifyRevision(input)] = true
					} else if input.revisesPackage {
						packageRevStrings[stringifyRevision(input)] = true
					}
				}

				// Signal that revise deps is over.
				args.revisionWaitGroup.Done()
			}
			readPackageDir = func(args readPackageDirArgs) {
				// Stow the received args for later assertions.
				actualReadPackageDirArgs = args

				args.importCounts.setImportCount("filepath1", 3)
				introduceRandomLag(0.4, 15)
				args.importCounts.setImportCount("filepath2", 2)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(101, "filepath1", `"github.com/a/b"`)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(102, "filepath2", `"github.com/a/b/c"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath2", 2)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(103, "filepath1", `"github.com/a/b/d/e"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath1", 1)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(104, "filepath1", `"github.com/a/b/d/e"`)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(105, "filepath2", `"github.com/h/i/j/k"`)
				// Throw in an extra big delay here to test that the SHA caching isn't
				// transient.
				introduceRandomLag(1.0, 120)
				args.importSpecChan <- generateTestImportSpecWithPos(106, "filepath2", `"github.com/a/b/x/y/z"`)

				// Close both channels once we're done.
				close(args.importSpecChan)
				close(args.packageSpecChan)
			}

			// Execute synchronously to make life easier.
			err := processDeps(processDepsArgs{
				io:                 io,
				ghSvc:              ghSvc,
				fetchSHA:           fetchSHA,
				reviseDeps:         reviseDeps,
				packageSHA:         packageSHA,
				packagePath:        packagePath,
				packageRepo:        packageRepo,
				packageAuthor:      packageAuthor,
				readPackageDir:     readPackageDir,
				packageVersionDate: packageVersionDate,
				newSpecWaitingList: newSpecWaitingList,
				newSyncedStringMap: newSyncedStringMap,
			})

			// Assert up a storm starting with fetchSHA.
			So(len(allFetchSHAArgs), ShouldEqual, 2)
			for _, args := range allFetchSHAArgs {
				So(args.ghSvc, ShouldNotBeNil)
				So(args.outputChan, ShouldNotBeNil)
				So(args.importPath, ShouldStartWith, `"github.com/`)
				So(args.packageSHA, ShouldEqual, packageSHA)
				So(args.packageRepo, ShouldEqual, packageRepo)
				So(args.packageAuthor, ShouldEqual, packageAuthor)
				So(args.pendingSHARequests, ShouldNotBeNil)
				So(args.packageVersionDate, ShouldResemble, packageVersionDate)
			}

			// Now make the asserts for revise deps.
			So(actualReviseDepsArgs.accumulatedErrors, ShouldNotBeNil)
			So(actualReviseDepsArgs.inputChan, ShouldNotBeNil)
			So(actualReviseDepsArgs.revisionWaitGroup, ShouldNotBeNil)
			So(actualReviseDepsArgs.syncedImportCounts, ShouldNotBeNil)
			So(importRevStrings[`filepath1:101:117:"gophr.pm/a/b@thisistheshafor"github.com/a/b""`], ShouldBeTrue)
			So(importRevStrings[`filepath2:102:120:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/c"`], ShouldBeTrue)
			So(importRevStrings[`filepath1:103:123:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/d/e"`], ShouldBeTrue)
			So(importRevStrings[`filepath1:104:124:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/d/e"`], ShouldBeTrue)
			So(importRevStrings[`filepath2:106:128:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/x/y/z"`], ShouldBeTrue)
			So(importRevStrings[`filepath2:105:125:"gophr.pm/h/i@thisistheshafor"github.com/h/i/j/k"/j/k"`], ShouldBeTrue)

			// Time for the readPackageDir asserts.
			So(actualReadPackageDirArgs.errors, ShouldNotBeNil)
			So(actualReadPackageDirArgs.generatedInternalDirName, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.importCounts, ShouldNotBeNil)
			So(actualReadPackageDirArgs.importSpecChan, ShouldNotBeNil)
			So(actualReadPackageDirArgs.io, ShouldNotBeNil)
			So(actualReadPackageDirArgs.packageDirPath, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.packageSpecChan, ShouldNotBeNil)

			// Finally, perform asserts for general outputs.
			So(err, ShouldBeNil)
		})
	})
}

// stringifyRevision turns revisions into strings for easy comparison.
func stringifyRevision(rev *revision) string {
	return fmt.Sprintf(
		`%v:%v:%v:%v`,
		rev.path,
		rev.fromIndex,
		rev.toIndex,
		string(rev.gophrURL[:]))
}

// introduceRandomLag conditionally pauses briefly. The goal here is to throw
// some fuzz into every test to catch race conditions.
func introduceRandomLag(chanceOfLag float32, maxMS int) {
	if chanceOfLag >= rand.Float32() {
		time.Sleep(time.Duration(float32(maxMS)*rand.Float32()) * time.Millisecond)
	}
}

// generateTestImportSpecWithPos generates an import spec that has position
// metadata included.
func generateTestImportSpecWithPos(pos int, filePath, importPath string) *importSpec {
	return &importSpec{
		imports: &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:     token.STRING,
				Value:    importPath,
				ValuePos: token.Pos(pos),
			},
		},
		filePath: filePath,
	}
}

type fakeSpecWaitingList struct {
	lock    *sync.Mutex
	specs   []*importSpec
	cleared bool
}

func newFakeSpecWaitingList(initialSpecs ...*importSpec) specWaitingList {
	return &fakeSpecWaitingList{
		lock:    &sync.Mutex{},
		specs:   initialSpecs,
		cleared: false,
	}
}

// add adds a spec to the waiting list and returns true if it was successful.
func (swl *fakeSpecWaitingList) add(spec *importSpec) bool {
	swl.lock.Lock()
	defer swl.lock.Unlock()

	if swl.cleared {
		return false
	}

	swl.specs = append(swl.specs, spec)
	return true
}

// clear returns every spec on the waiting list and prevents all future adds from
// succeeding.
func (swl *fakeSpecWaitingList) clear() []*importSpec {
	swl.lock.Lock()
	defer swl.lock.Unlock()

	specs := swl.specs
	swl.specs = nil
	swl.cleared = true
	return specs
}
