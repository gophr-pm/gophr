package verdeps

import (
	"errors"
	"go/ast"
	"go/token"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/gophr-pm/gophr/lib/io"
	. "github.com/smartystreets/goconvey/convey"
	"github.com/stretchr/testify/mock"
)

func TestReviseDeps(t *testing.T) {
	Convey("Given a steady stream of revisions", t, func() {
		Convey("If there are no I/O issues, all revisions should be applied correctly", func() {
			var (
				io                = io.NewMockIO()
				inputChan         = make(chan *revision)
				importCounts      = newSyncedImportCounts()
				accumulatedErrors = newSyncedErrors()
				revisionWaitGroup = &sync.WaitGroup{}
			)

			// Setup the IO mock.
			io.On("ReadFile", "filepath1").Once().Return([]byte(goFile1), nil)
			io.On("ReadFile", "filepath2").Once().Return([]byte(goFile2), nil)
			io.On(
				"WriteFile",
				"filepath1",
				mock.MatchedBy(equalsRevisedGoFile1),
				os.FileMode(0644)).
				Once().
				Return(nil)
			io.On(
				"WriteFile",
				"filepath2",
				mock.MatchedBy(equalsRevisedGoFile2),
				os.FileMode(0644)).
				Once().
				Return(nil)

			// Start revise deps in the background.
			revisionWaitGroup.Add(1)
			go reviseDeps(reviseDepsArgs{
				io:                 io,
				inputChan:          inputChan,
				composeBytesDiffs:  composeBytesDiffs,
				revisionWaitGroup:  revisionWaitGroup,
				accumulatedErrors:  accumulatedErrors,
				syncedImportCounts: importCounts,
			})

			// Enqueue all the revisions.
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath1", 2)
			inputChan <- newTestImportRevision(
				goFile1,
				"filepath1",
				`"github.com/a/b"`,
				`"gophr.pm/a/b@somesha"`)
			introduceRandomLag(0.4, 15)
			inputChan <- newTestPackageRevision(
				goFile1,
				"filepath1")
			introduceRandomLag(0.4, 15)
			inputChan <- newTestImportRevision(
				goFile1,
				"filepath1",
				`"github.com/c/d/e"`,
				`"gophr.pm/c/d@somesha/e"`)
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath2", 1)
			inputChan <- newTestImportRevision(
				goFile2,
				"filepath2",
				`"github.com/j/k/l/m"`,
				`"gophr.pm/j/k@somesha/l/m"`)
			introduceRandomLag(0.4, 15)
			inputChan <- newTestPackageRevision(
				goFile1,
				"filepath2")
			close(inputChan)

			// Wait until reviseDeps exits.
			revisionWaitGroup.Wait()

			// Assert that there were no issues, and that files have been updated
			// correctly.
			So(accumulatedErrors.len(), ShouldEqual, 0)
			io.AssertExpectations(t)
		})

		Convey("Missing revisions should be reported, and not prevent other revisions in the same file", func() {
			var (
				io                = io.NewMockIO()
				inputChan         = make(chan *revision)
				importCounts      = newSyncedImportCounts()
				accumulatedErrors = newSyncedErrors()
				revisionWaitGroup = &sync.WaitGroup{}
			)

			// Setup the IO mock.
			io.On("ReadFile", "filepath2").Once().Return([]byte(goFile2), nil)
			io.On(
				"WriteFile",
				"filepath2",
				mock.MatchedBy(equalsRevisedGoFile2),
				os.FileMode(0644)).
				Once().
				Return(nil)

			// Start revise deps in the background.
			revisionWaitGroup.Add(1)
			go reviseDeps(reviseDepsArgs{
				io:                 io,
				inputChan:          inputChan,
				composeBytesDiffs:  composeBytesDiffs,
				revisionWaitGroup:  revisionWaitGroup,
				accumulatedErrors:  accumulatedErrors,
				syncedImportCounts: importCounts,
			})

			// Enqueue all the revisions.
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath2", 2)
			inputChan <- newTestImportRevision(
				goFile2,
				"filepath2",
				`"github.com/j/k/l/m"`,
				`"gophr.pm/j/k@somesha/l/m"`)
			close(inputChan)

			// Wait until reviseDeps exits.
			revisionWaitGroup.Wait()

			// Assert that there were no issues, and that files have been updated
			// correctly.
			So(accumulatedErrors.len(), ShouldEqual, 0)
			io.AssertExpectations(t)
		})

		Convey("Errors while reading files should bubble up", func() {
			var (
				io                = io.NewMockIO()
				inputChan         = make(chan *revision)
				importCounts      = newSyncedImportCounts()
				accumulatedErrors = newSyncedErrors()
				revisionWaitGroup = &sync.WaitGroup{}
			)

			// Setup the IO mock.
			io.On("ReadFile", "filepath2").Once().Return(
				[]byte{},
				errors.New("this is an error"))

			// Start revise deps in the background.
			revisionWaitGroup.Add(1)
			go reviseDeps(reviseDepsArgs{
				io:                 io,
				inputChan:          inputChan,
				composeBytesDiffs:  composeBytesDiffs,
				revisionWaitGroup:  revisionWaitGroup,
				accumulatedErrors:  accumulatedErrors,
				syncedImportCounts: importCounts,
			})

			// Enqueue all the revisions.
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath2", 1)
			inputChan <- newTestImportRevision(
				goFile2,
				"filepath2",
				`"github.com/j/k/l/m"`,
				`"gophr.pm/j/k@somesha/l/m"`)
			introduceRandomLag(0.4, 15)
			inputChan <- newTestPackageRevision(
				goFile1,
				"filepath2")
			close(inputChan)

			// Wait until reviseDeps exits.
			revisionWaitGroup.Wait()

			// Assert that there were no issues, and that files have been updated
			// correctly.
			So(accumulatedErrors.len(), ShouldEqual, 1)
		})

		Convey("Errors while writing files should bubble up", func() {
			var (
				io                = io.NewMockIO()
				inputChan         = make(chan *revision)
				importCounts      = newSyncedImportCounts()
				accumulatedErrors = newSyncedErrors()
				revisionWaitGroup = &sync.WaitGroup{}
			)

			// Setup the IO mock.
			io.On("ReadFile", "filepath2").Once().Return([]byte(goFile2), nil)
			io.On(
				"WriteFile",
				"filepath2",
				mock.MatchedBy(equalsRevisedGoFile2),
				os.FileMode(0644)).
				Once().
				Return(errors.New("this is an error"))

			// Start revise deps in the background.
			revisionWaitGroup.Add(1)
			go reviseDeps(reviseDepsArgs{
				io:                 io,
				inputChan:          inputChan,
				composeBytesDiffs:  composeBytesDiffs,
				revisionWaitGroup:  revisionWaitGroup,
				accumulatedErrors:  accumulatedErrors,
				syncedImportCounts: importCounts,
			})

			// Enqueue all the revisions.
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath2", 1)
			inputChan <- newTestImportRevision(
				goFile2,
				"filepath2",
				`"github.com/j/k/l/m"`,
				`"gophr.pm/j/k@somesha/l/m"`)
			introduceRandomLag(0.4, 15)
			inputChan <- newTestPackageRevision(
				goFile1,
				"filepath2")
			close(inputChan)

			// Wait until reviseDeps exits.
			revisionWaitGroup.Wait()

			// Assert that there were no issues, and that files have been updated
			// correctly.
			So(accumulatedErrors.len(), ShouldEqual, 1)
		})

		Convey("Errors while finding import path indicies should bubble up", func() {
			var (
				io                = io.NewMockIO()
				inputChan         = make(chan *revision)
				importCounts      = newSyncedImportCounts()
				accumulatedErrors = newSyncedErrors()
				revisionWaitGroup = &sync.WaitGroup{}
			)

			// Setup the IO mock.
			io.On("ReadFile", "filepath2").Once().Return([]byte(goFile2), nil)
			io.On(
				"WriteFile",
				"filepath2",
				mock.MatchedBy(equalsRevisedGoFile2),
				os.FileMode(0644)).
				Once().
				Return(nil)

			// Start revise deps in the background.
			revisionWaitGroup.Add(1)
			go reviseDeps(reviseDepsArgs{
				io:                 io,
				inputChan:          inputChan,
				composeBytesDiffs:  composeBytesDiffs,
				revisionWaitGroup:  revisionWaitGroup,
				accumulatedErrors:  accumulatedErrors,
				syncedImportCounts: importCounts,
			})

			// Enqueue all the revisions.
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath2", 1)
			inputChan <- newTestImportRevisionWithImpossiblePos(
				"filepath2",
				`"github.com/j/k/l/m"`,
				`"gophr.pm/j/k@somesha/l/m"`)
			introduceRandomLag(0.4, 15)
			inputChan <- newTestPackageRevision(
				goFile1,
				"filepath2")
			close(inputChan)

			// Wait until reviseDeps exits.
			revisionWaitGroup.Wait()

			// Assert that there were no issues, and that files have been updated
			// correctly.
			So(accumulatedErrors.len(), ShouldEqual, 1)
		})

		Convey("Errors while composing bytes diffs should bubble up", func() {
			var (
				io                = io.NewMockIO()
				inputChan         = make(chan *revision)
				importCounts      = newSyncedImportCounts()
				accumulatedErrors = newSyncedErrors()
				revisionWaitGroup = &sync.WaitGroup{}
				composeBytesDiffs = func(
					bytes []byte,
					diffs []bytesDiff,
				) ([]byte, error) {
					return nil, errors.New("this is an error")
				}
			)

			// Setup the IO mock.
			io.On("ReadFile", "filepath2").Once().Return([]byte(goFile2), nil)

			// Start revise deps in the background.
			revisionWaitGroup.Add(1)
			go reviseDeps(reviseDepsArgs{
				io:                 io,
				inputChan:          inputChan,
				composeBytesDiffs:  composeBytesDiffs,
				revisionWaitGroup:  revisionWaitGroup,
				accumulatedErrors:  accumulatedErrors,
				syncedImportCounts: importCounts,
			})

			// Enqueue all the revisions.
			introduceRandomLag(0.4, 15)
			importCounts.setImportCount("filepath2", 1)
			inputChan <- newTestImportRevision(
				goFile2,
				"filepath2",
				`"github.com/j/k/l/m"`,
				`"gophr.pm/j/k@somesha/l/m"`)
			introduceRandomLag(0.4, 15)
			inputChan <- newTestPackageRevision(
				goFile1,
				"filepath2")
			close(inputChan)

			// Wait until reviseDeps exits.
			revisionWaitGroup.Wait()

			// Assert that there were no issues, and that files have been updated
			// correctly.
			So(accumulatedErrors.len(), ShouldEqual, 1)
		})
	})
}

func TestFindImportPathBoundaries(t *testing.T) {
	Convey("Given an index in file data", t, func() {
		Convey("Should find end quote of import path", func() {
			from, to, err := findImportPathBoundaries(
				[]byte(`sdfk fkjlsdfg kjsdf"github.com/a/b"skd lkjsdd ksldjk`),
				23,
				31)
			So(from, ShouldEqual, 19)
			So(to, ShouldEqual, 35)
			So(err, ShouldBeNil)
		})

		Convey("Should return an error when there is no end quote", func() {
			_, _, err := findImportPathBoundaries(
				[]byte(`sdfk fkjlsdfg kjsdf"github.com/a/bskd lkjsdd ksldjk`),
				23,
				31)
			So(err, ShouldNotBeNil)
		})
	})
}

/********************************* TEST DATA **********************************/

const (
	goFile1 = `
package test // import "github.com/x/y"

import (
  "github.com/a/b"
  e "github.com/c/d/e"
  "github.com/f/g"
)

func main() int {
  b.B()
  e.E()
  g.G()
  test()
  return 0
}
`
	revisedGoFile1 = `
package test

import (
  "gophr.pm/a/b@somesha"
  e "gophr.pm/c/d@somesha/e"
  "github.com/f/g"
)

func main() int {
  b.B()
  e.E()
  g.G()
  test()
  return 0
}
`
	goFile2 = `
package test

import (
  "fmt"

  "github.com/h/i"
  "github.com/j/k/l/m"
)

func test() {
  i.I()
  m.M()
  fmt.Println("this is a test")
}
`
	revisedGoFile2 = `
package test

import (
  "fmt"

  "github.com/h/i"
  "gophr.pm/j/k@somesha/l/m"
)

func test() {
  i.I()
  m.M()
  fmt.Println("this is a test")
}
`
)

/************************ TEST STRUCT CREATION HELPERS ************************/

func newTestImportRevisionWithImpossiblePos(
	filePath string,
	oldImportPath string,
	newImportPath string,
) *revision {
	return newImportRevision(
		&importSpec{
			imports: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    oldImportPath,
					ValuePos: -23,
				},
			},
			filePath: filePath,
		},
		[]byte(newImportPath))
}

func newTestImportRevision(
	fileData string,
	filePath string,
	oldImportPath string,
	newImportPath string,
) *revision {
	return newImportRevision(
		&importSpec{
			imports: &ast.ImportSpec{
				Path: &ast.BasicLit{
					Kind:     token.STRING,
					Value:    oldImportPath,
					ValuePos: token.Pos(strings.Index(fileData, oldImportPath) + 1),
				},
			},
			filePath: filePath,
		},
		[]byte(newImportPath))
}

func newTestPackageRevision(fileData, filePath string) *revision {
	return newPackageRevision(
		&packageSpec{
			filePath:   filePath,
			startIndex: strings.Index(fileData, `package test`),
		})
}

/************************ BYTE SLICE EQUALITY HELPERS *************************/

func equalsRevisedGoFile1(bytes []byte) bool {
	equal := byteSlicesEqual(bytes, []byte(revisedGoFile1))
	return equal
}

func equalsRevisedGoFile2(bytes []byte) bool {
	equal := byteSlicesEqual(bytes, []byte(revisedGoFile2))
	return equal
}

func byteSlicesEqual(a, b []byte) bool {
	if a == nil && b == nil {
		return true
	}

	if a == nil || b == nil {
		return false
	}

	if len(a) != len(b) {
		return false
	}

	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}

	return true
}
