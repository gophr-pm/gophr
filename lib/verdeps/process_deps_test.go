package verdeps

import (
	"fmt"
	"go/ast"
	"go/token"
	"math/rand"
	"sync"
	"testing"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
	"github.com/gophr-pm/gophr/lib/io"
	. "github.com/smartystreets/goconvey/convey"
)

func TestProcessDeps(t *testing.T) {
	Convey("Given go package with dependencies", t, func() {
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

			fmt.Println("\nfetchSHA for", args.importPath)

			args.outputChan <- &importPathSHA{
				sha:        "thisistheshafor" + args.importPath,
				importPath: args.importPath,
			}
		}
		reviseDeps = func(args reviseDepsArgs) {
			// Stow the received args for later assertions.
			actualReviseDepsArgs = args

			fmt.Println("\nreviseDeps begin")
			// Read the input revisions, and record what comes through.
			for input := range args.inputChan {
				if input.revisesImport {
					fmt.Println("\nreviseDeps revisesImport", stringifyRevision(input))
					importRevStrings[stringifyRevision(input)] = true
				} else if input.revisesPackage {
					fmt.Println("\nreviseDeps revisesPackage", stringifyRevision(input))
					packageRevStrings[stringifyRevision(input)] = true
				}
			}

			fmt.Println("\nreviseDeps end")
			// Signal that revise deps is over.
			args.revisionWaitGroup.Done()
		}
		readPackageDir = func(args readPackageDirArgs) {
			// Stow the received args for later assertions.
			actualReadPackageDirArgs = args

			fmt.Println("\nreadPackageDir begin")
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
			fmt.Println("\nreadPackageDir end")
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
