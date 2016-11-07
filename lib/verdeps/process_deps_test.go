package verdeps

import (
	"errors"
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
					time.Local)
			)

			// Create fakes of the worker functions passed into processDeps.
			fetchSHA = func(args fetchSHAArgs) {
				// Record this invocation of fetchSHA in a synchronized map.
				shaRequestsLock.Lock()
				shaRequests[args.importPath] = &args
				shaRequestsLock.Unlock()

				// After a pause, signal that this SHA was fetched.
				introduceRandomLag(0.6, 30)

				args.outputChan <- newFetchSHASuccess(
					args.importPath,
					"thisistheshafor"+args.importPath)
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
				io:                      io,
				ghSvc:                   ghSvc,
				fetchSHA:                fetchSHA,
				reviseDeps:              reviseDeps,
				packageSHA:              packageSHA,
				packagePath:             packagePath,
				packageRepo:             packageRepo,
				packageAuthor:           packageAuthor,
				readPackageDir:          readPackageDir,
				packageVersionDate:      packageVersionDate,
				newSpecWaitingList:      newSpecWaitingList,
				newSyncedStringMap:      newSyncedStringMap,
				newSyncedWaitingListMap: newSyncedWaitingListMap,
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
				So(args.packageVersionDate, ShouldResemble, packageVersionDate)
			}

			// Now make the asserts for revise deps.
			So(actualReviseDepsArgs.accumulatedErrors, ShouldNotBeNil)
			So(actualReviseDepsArgs.composeBytesDiffs, ShouldNotBeNil)
			So(actualReviseDepsArgs.inputChan, ShouldNotBeNil)
			So(actualReviseDepsArgs.io, ShouldNotBeNil)
			So(actualReviseDepsArgs.revisionWaitGroup, ShouldNotBeNil)
			So(actualReviseDepsArgs.syncedImportCounts, ShouldNotBeNil)
			So(importRevStrings[`filepath1:101:119:"gophr.pm/c/d@thisistheshafor"github.com/c/d/e"/e"`], ShouldBeTrue)
			So(importRevStrings[`filepath2:102:118:"gophr.pm/f/g@thisistheshafor"github.com/f/g""`], ShouldBeTrue)
			So(importRevStrings[`filepath1:103:119:"gophr.pm/b/c@thisistheshafor"github.com/b/c""`], ShouldBeTrue)
			So(importRevStrings[`filepath1:104:120:"gophr.pm/a/b@thisistheshafor"github.com/a/b""`], ShouldBeTrue)
			So(importRevStrings[`filepath2:105:125:"gophr.pm/h/i@thisistheshafor"github.com/h/i/j/k"/j/k"`], ShouldBeTrue)
			So(importRevStrings[`filepath3:106:122:"gophr.pm/b/c@thisistheshafor"github.com/b/c""`], ShouldBeTrue)
			So(packageRevStrings[`filepath1:1:0:`], ShouldBeTrue)
			So(packageRevStrings[`filepath2:2:0:`], ShouldBeTrue)
			So(packageRevStrings[`filepath3:3:0:`], ShouldBeTrue)

			// Time for the readPackageDir asserts.
			So(actualReadPackageDirArgs.errors, ShouldNotBeNil)
			So(actualReadPackageDirArgs.generatedInternalDirName, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.importCounts, ShouldNotBeNil)
			So(actualReadPackageDirArgs.importSpecChan, ShouldNotBeNil)
			So(actualReadPackageDirArgs.io, ShouldNotBeNil)
			So(actualReadPackageDirArgs.packageDirPath, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.packageSpecChan, ShouldNotBeNil)
			So(actualReadPackageDirArgs.traversePackageDir, ShouldNotBeNil)

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
					time.Local)
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

				args.outputChan <- newFetchSHASuccess(
					args.importPath,
					"thisistheshafor"+args.importPath)
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
				args.importCounts.setImportCount("filepath2", 3)
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
				io:                      io,
				ghSvc:                   ghSvc,
				fetchSHA:                fetchSHA,
				reviseDeps:              reviseDeps,
				packageSHA:              packageSHA,
				packagePath:             packagePath,
				packageRepo:             packageRepo,
				packageAuthor:           packageAuthor,
				readPackageDir:          readPackageDir,
				packageVersionDate:      packageVersionDate,
				newSpecWaitingList:      newSpecWaitingList,
				newSyncedStringMap:      newSyncedStringMap,
				newSyncedWaitingListMap: newSyncedWaitingListMap,
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
				So(args.packageVersionDate, ShouldResemble, packageVersionDate)
			}

			// Now make the asserts for revise deps.
			So(actualReviseDepsArgs.accumulatedErrors, ShouldNotBeNil)
			So(actualReviseDepsArgs.composeBytesDiffs, ShouldNotBeNil)
			So(actualReviseDepsArgs.inputChan, ShouldNotBeNil)
			So(actualReviseDepsArgs.io, ShouldNotBeNil)
			So(actualReviseDepsArgs.revisionWaitGroup, ShouldNotBeNil)
			So(actualReviseDepsArgs.syncedImportCounts, ShouldNotBeNil)
			So(importRevStrings[`filepath1:101:117:"gophr.pm/a/b@thisistheshafor"github.com/a/b""`], ShouldBeTrue)
			So(importRevStrings[`filepath2:102:120:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/c"`], ShouldBeTrue)
			So(importRevStrings[`filepath1:103:123:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/d/e"`], ShouldBeTrue)
			So(importRevStrings[`filepath1:104:124:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/d/e"`], ShouldBeTrue)
			So(importRevStrings[`filepath2:106:128:"gophr.pm/a/b@thisistheshafor"github.com/a/b"/x/y/z"`], ShouldBeTrue)
			So(importRevStrings[`filepath2:105:125:"gophr.pm/h/i@thisistheshafor"github.com/h/i/j/k"/j/k"`], ShouldBeTrue)
			So(packageRevStrings[`filepath1:1:0:`], ShouldBeTrue)
			So(packageRevStrings[`filepath2:2:0:`], ShouldBeTrue)

			// Time for the readPackageDir asserts.
			So(actualReadPackageDirArgs.errors, ShouldNotBeNil)
			So(actualReadPackageDirArgs.generatedInternalDirName, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.importCounts, ShouldNotBeNil)
			So(actualReadPackageDirArgs.importSpecChan, ShouldNotBeNil)
			So(actualReadPackageDirArgs.io, ShouldNotBeNil)
			So(actualReadPackageDirArgs.packageDirPath, ShouldNotBeEmpty)
			So(actualReadPackageDirArgs.packageSpecChan, ShouldNotBeNil)
			So(actualReadPackageDirArgs.traversePackageDir, ShouldNotBeNil)

			// Finally, perform asserts for general outputs.
			So(err, ShouldBeNil)
		})

		Convey("An error should be raised if an import spec cannot be enqueued to receive a SHA", func() {
			var (
				reviseDeps         depsReviser
				readPackageDir     packageDirReader
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
					time.Local)
			)

			reviseDeps = func(args reviseDepsArgs) {
				// Signal that revise deps is over when the function exits.
				defer args.revisionWaitGroup.Done()

				// Read the input revisions, and record what comes through.
				for input := range args.inputChan {
					if input.revisesImport {
						importRevStrings[stringifyRevision(input)] = true
					} else if input.revisesPackage {
						packageRevStrings[stringifyRevision(input)] = true
					}
				}
			}
			readPackageDir = func(args readPackageDirArgs) {
				args.importCounts.setImportCount("filepath1", 1)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(
					101,
					"filepath1",
					`"github.com/a/b"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath1", 1)

				// Close both channels once we're done.
				close(args.importSpecChan)
				close(args.packageSpecChan)
			}

			var (
				waitingSpecs    = newFakeSyncedWaitingListMap()
				fetchSHAResults = newFakeSyncedStringMap()
				specWaitingList = newFakeSpecWaitingList(fetchSHAResults, "")
			)

			// Make sure the SHA "doesn't exist yet".
			fetchSHAResults.overrideGet(
				importPathHashOf(`"github.com/a/b"`),
				"",
				false)
			// Ensure that there is already a waiting list.
			waitingSpecs.overrideGet(
				importPathHashOf(`"github.com/a/b"`),
				specWaitingList,
				true)
			// The add must fail.
			specWaitingList.overwriteAdd(generateTestImportSpec(
				"filepath1",
				`"github.com/a/b"`))

			// Execute synchronously to make life easier.
			err := processDeps(processDepsArgs{
				io:                      nil,
				ghSvc:                   nil,
				fetchSHA:                nil,
				reviseDeps:              reviseDeps,
				packageSHA:              "",
				packagePath:             "",
				packageRepo:             "",
				packageAuthor:           "",
				readPackageDir:          readPackageDir,
				packageVersionDate:      packageVersionDate,
				newSpecWaitingList:      specWaitingList.creator(),
				newSyncedStringMap:      fetchSHAResults.creator(),
				newSyncedWaitingListMap: waitingSpecs.creator(),
			})

			// The error should bubble up.
			So(err, ShouldNotBeNil)
			// Should generate an error instead of enqueuing the revision.
			So(
				importRevStrings[`filepath1:101:117:"gophr.pm/a/b@thisistheshafor"github.com/a/b""`],
				ShouldBeFalse)
		})

		Convey("The SHA should still be bound if the SHA of an import is received mid-waiting-list-enqueue", func() {
			var (
				reviseDeps         depsReviser
				readPackageDir     packageDirReader
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
					time.Local)
			)

			reviseDeps = func(args reviseDepsArgs) {
				// Signal that revise deps is over when the function exits.
				defer args.revisionWaitGroup.Done()

				// Read the input revisions, and record what comes through.
				for input := range args.inputChan {
					if input.revisesImport {
						importRevStrings[stringifyRevision(input)] = true
					} else if input.revisesPackage {
						packageRevStrings[stringifyRevision(input)] = true
					}
				}
			}
			readPackageDir = func(args readPackageDirArgs) {
				args.importCounts.setImportCount("filepath1", 1)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(
					101,
					"filepath1",
					`"github.com/a/b"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath1", 1)

				// Close both channels once we're done.
				close(args.importSpecChan)
				close(args.packageSpecChan)
			}

			var (
				waitingSpecs    = newFakeSyncedWaitingListMap()
				fetchSHAResults = newFakeSyncedStringMap()
				specWaitingList = newFakeSpecWaitingList(
					fetchSHAResults,
					"someshathatdoesnotmatter")
			)

			// Make sure the SHA "doesn't exist yet".
			fetchSHAResults.overrideGet(
				importPathHashOf(`"github.com/a/b"`),
				"",
				false)
			// Ensure that there is already a waiting list.
			waitingSpecs.overrideGet(
				importPathHashOf(`"github.com/a/b"`),
				specWaitingList,
				true)
			// The add must fail.
			specWaitingList.overwriteAdd(generateTestImportSpec(
				"filepath1",
				`"github.com/a/b"`))

			// Execute synchronously to make life easier.
			err := processDeps(processDepsArgs{
				io:                      nil,
				ghSvc:                   nil,
				fetchSHA:                nil,
				reviseDeps:              reviseDeps,
				packageSHA:              "",
				packagePath:             "",
				packageRepo:             "",
				packageAuthor:           "",
				readPackageDir:          readPackageDir,
				packageVersionDate:      packageVersionDate,
				newSpecWaitingList:      specWaitingList.creator(),
				newSyncedStringMap:      fetchSHAResults.creator(),
				newSyncedWaitingListMap: waitingSpecs.creator(),
			})

			// There is no error since the SHA ends up paired with the import.
			So(err, ShouldBeNil)
			// Should not generate an error, and instead enqueue the revision.
			So(
				importRevStrings[`filepath1:101:117:"gophr.pm/a/b@someshathatdoesnotmatter"`],
				ShouldBeTrue)
		})
		Convey("An error not should be raised if a SHA request fails", func() {
			var (
				fetchSHA           shaFetcher
				reviseDeps         depsReviser
				readPackageDir     packageDirReader
				packageVersionDate = time.Date(
					2016,
					time.April,
					8,
					14,
					12,
					0,
					0,
					time.Local)
			)

			// Create fakes of the worker functions passed into processDeps.
			fetchSHA = func(args fetchSHAArgs) {
				introduceRandomLag(0.5, 30)
				args.outputChan <- newFetchSHAFailure(errors.New("this is an error"))
			}
			reviseDeps = func(args reviseDepsArgs) {
				// Read the input revisions, and record what comes through.
				for range args.inputChan {
					// No-op (we don't care about what comes through.)
				}

				// Don't expect anything to come through - so, just exit.
				args.revisionWaitGroup.Done()
			}
			readPackageDir = func(args readPackageDirArgs) {
				args.importCounts.setImportCount("filepath1", 1)
				introduceRandomLag(0.4, 15)
				args.importSpecChan <- generateTestImportSpecWithPos(
					101,
					"filepath1",
					`"github.com/a/b"`)
				introduceRandomLag(0.4, 15)
				args.packageSpecChan <- generateTestPackageSpec("filepath1", 1)

				// Close both channels once we're done.
				close(args.importSpecChan)
				close(args.packageSpecChan)
			}

			// Execute synchronously to make life easier.
			err := processDeps(processDepsArgs{
				io:                      nil,
				ghSvc:                   nil,
				fetchSHA:                fetchSHA,
				reviseDeps:              reviseDeps,
				packageSHA:              "",
				packagePath:             "",
				packageRepo:             "",
				packageAuthor:           "",
				readPackageDir:          readPackageDir,
				packageVersionDate:      packageVersionDate,
				newSpecWaitingList:      newSpecWaitingList,
				newSyncedStringMap:      newSyncedStringMap,
				newSyncedWaitingListMap: newSyncedWaitingListMap,
			})

			// The error should bubble up.
			So(err, ShouldBeNil)
		})
	})
}

/*********************************** HELPERS **********************************/

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
func generateTestImportSpecWithPos(
	pos int,
	filePath string,
	importPath string,
) *importSpec {
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

/************************ FAKE SYNCED WAITING LIST MAP ************************/

type fakeSyncedWaitingListMap struct {
	swlm              syncedWaitingListMapImpl
	getOverrideKey    string
	getOverrideLock   sync.Mutex
	getOverrideValue  specWaitingList
	getOverrideExists bool
}

func newFakeSyncedWaitingListMap() *fakeSyncedWaitingListMap {
	return &fakeSyncedWaitingListMap{
		swlm: syncedWaitingListMapImpl{
			values: make(map[string]specWaitingList),
			lock:   &sync.RWMutex{},
		},
	}
}

func (m *fakeSyncedWaitingListMap) creator() syncedWaitingListMapCreator {
	return func() syncedWaitingListMap {
		return m
	}
}
func (m *fakeSyncedWaitingListMap) overrideGet(
	key string,
	value specWaitingList,
	exists bool,
) {
	m.getOverrideLock.Lock()
	m.getOverrideKey = key
	m.getOverrideValue = value
	m.getOverrideExists = exists
	m.getOverrideLock.Unlock()
}
func (m *fakeSyncedWaitingListMap) get(key string) (specWaitingList, bool) {
	m.getOverrideLock.Lock()
	if len(m.getOverrideKey) > 0 && m.getOverrideKey == key {
		m.getOverrideLock.Unlock()
		return m.getOverrideValue, m.getOverrideExists
	}
	m.getOverrideLock.Unlock()

	return m.swlm.get(key)
}
func (m *fakeSyncedWaitingListMap) setIfAbsent(
	key string,
	value specWaitingList,
) {
	m.swlm.setIfAbsent(key, value)
}
func (m *fakeSyncedWaitingListMap) clear() {
	m.swlm.clear()
}
func (m *fakeSyncedWaitingListMap) each(
	fn func(string, specWaitingList),
) syncedWaitingListMap {
	m.swlm.each(fn)
	return m
}

/*************************** FAKE SPEC WAITING LIST ***************************/

type fakeSpecWaitingList struct {
	wl          specWaitingListImpl
	sha         string
	fssm        *fakeSyncedStringMap
	addOverride *importSpec
}

func newFakeSpecWaitingList(
	fssm *fakeSyncedStringMap,
	sha string,
) *fakeSpecWaitingList {
	return &fakeSpecWaitingList{
		wl: specWaitingListImpl{
			lock:    &sync.Mutex{},
			cleared: false,
		},
		sha:  sha,
		fssm: fssm,
	}
}

func (swl *fakeSpecWaitingList) creator() specWaitingListCreator {
	return func(initialSpecs ...*importSpec) specWaitingList {
		swl.wl.specs = initialSpecs
		return swl
	}
}
func (swl *fakeSpecWaitingList) overwriteAdd(key *importSpec) {
	swl.addOverride = key
}
func (swl *fakeSpecWaitingList) add(spec *importSpec) bool {
	// If the catalyst gets hit. Return false (supposed to fail).
	if swl.addOverride != nil &&
		spec.filePath == swl.addOverride.filePath &&
		spec.imports.Path.Value == swl.addOverride.imports.Path.Value {
		// If there is an fssm, then override its next get.
		if len(swl.sha) <= 0 {
			swl.fssm.overrideGet(
				importPathHashOf(spec.imports.Path.Value),
				"",
				false)
		} else {
			swl.fssm.overrideGet(
				importPathHashOf(spec.imports.Path.Value),
				swl.sha,
				true)
		}

		return false
	}

	return swl.wl.add(spec)
}
func (swl *fakeSpecWaitingList) clear() []*importSpec {
	return swl.wl.clear()
}

/*************************** FAKE SYNCED STRING MAP ***************************/

type fakeSyncedStringMap struct {
	ssm               syncedStringMapImpl
	getOverrideKey    string
	getOverrideLock   sync.Mutex
	getOverrideValue  string
	getOverrideExists bool
}

func newFakeSyncedStringMap() *fakeSyncedStringMap {
	return &fakeSyncedStringMap{
		ssm: syncedStringMapImpl{
			values: make(map[string]string),
			lock:   &sync.RWMutex{},
		},
	}
}

func (m *fakeSyncedStringMap) creator() syncedStringMapCreator {
	return func() syncedStringMap {
		return m
	}
}
func (m *fakeSyncedStringMap) overrideGet(key, value string, exists bool) {
	m.getOverrideLock.Lock()
	m.getOverrideKey = key
	m.getOverrideValue = value
	m.getOverrideExists = exists
	m.getOverrideLock.Unlock()
}
func (m *fakeSyncedStringMap) get(key string) (string, bool) {
	m.getOverrideLock.Lock()
	if len(m.getOverrideKey) > 0 && key == m.getOverrideKey {
		val, exists := m.getOverrideValue, m.getOverrideExists
		m.getOverrideLock.Unlock()
		return val, exists
	}
	m.getOverrideLock.Unlock()

	return m.ssm.get(key)
}
func (m *fakeSyncedStringMap) set(key, value string) {
	m.ssm.set(key, value)
}
