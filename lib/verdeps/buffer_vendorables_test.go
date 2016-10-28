package verdeps

import (
	"go/ast"
	"go/token"
	"strconv"
	"sync"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestBufferVendorables(t *testing.T) {
	Convey("Given a bunch of import and package specs", t, func() {
		var (
			waitGroup             = &sync.WaitGroup{}
			importCounts          = newSyncedImportCounts()
			inputImportSpecChan   = make(chan *importSpec)
			outputImportSpecChan  = make(chan *importSpec)
			currentVendorContext  = newVendorContext()
			inputPackageSpecChan  = make(chan *packageSpec)
			outputPackageSpecChan = make(chan *packageSpec)
		)

		// Add vendored packages.
		currentVendorContext.add(`github.com/a/b`)
		currentVendorContext.add(`github.com/x/y`)

		// Start buffering in the background.
		waitGroup.Add(1)
		go bufferVendorables(bufferVendorablesArgs{
			waitGroup:             waitGroup,
			importCounts:          importCounts,
			inputImportSpecChan:   inputImportSpecChan,
			outputImportSpecChan:  outputImportSpecChan,
			currentVendorContext:  currentVendorContext,
			inputPackageSpecChan:  inputPackageSpecChan,
			outputPackageSpecChan: outputPackageSpecChan,
		})

		// Add a bunch of specs into the mix.
		inputImportSpecChan <- generateTestImportSpec("filepath1", `"github.com/d/e"`)
		inputPackageSpecChan <- generateTestPackageSpec("filepath1", 34)
		inputImportSpecChan <- generateTestImportSpec("filepath1", `"github.com/f/g"`)
		inputImportSpecChan <- generateTestImportSpec("filepath1", `"github.com/g/h"`)
		inputPackageSpecChan <- generateTestPackageSpec("filepath2", 17821)
		inputImportSpecChan <- generateTestImportSpec("filepath2", `"github.com/a/b"`)
		inputImportSpecChan <- generateTestImportSpec("filepath2", `"github.com/b/c"`)
		inputImportSpecChan <- generateTestImportSpec("filepath2", `"github.com/c/d"`)
		inputPackageSpecChan <- generateTestPackageSpec("filepath3", 1)
		inputImportSpecChan <- generateTestImportSpec("filepath3", `"ignore.me/r/s"`)
		inputImportSpecChan <- generateTestImportSpec("filepath3", `"github.com/x/y"`)
		// Prepare to read outputs.
		close(inputImportSpecChan)
		close(inputPackageSpecChan)

		var (
			numberOfImportSpecs         = 8
			numberOfPackageSpecs        = 3
			numberOfVendoredImportSpecs = 2

			outputImportSpecStrings  = make(map[string]bool)
			outputPackageSpecStrings = make(map[string]bool)
		)

		// Make sure these channels get closed before the test completes.
		defer close(outputImportSpecChan)
		defer close(outputPackageSpecChan)

		i := 0
		for spec := range outputImportSpecChan {
			key := spec.filePath + ":" + spec.imports.Path.Value
			outputImportSpecStrings[key] = true

			// Break after all the specs we put in come out again.
			if i++; i >= (numberOfImportSpecs - numberOfVendoredImportSpecs) {
				break
			}
		}

		i = 0
		for spec := range outputPackageSpecChan {
			key := spec.filePath + ":" + strconv.Itoa(spec.startIndex)
			outputPackageSpecStrings[key] = true

			// Break after all the specs we put in come out again.
			if i++; i >= numberOfPackageSpecs {
				break
			}
		}

		// Make sure everything lines up as expected.
		Convey("Function shouldn't deadlock", func() {
			waitGroup.Wait()
		})
		Convey("Unvendored github imports should come through", func() {
			So(outputImportSpecStrings[`filepath1:"github.com/d/e"`], ShouldBeTrue)
			So(outputImportSpecStrings[`filepath1:"github.com/f/g"`], ShouldBeTrue)
			So(outputImportSpecStrings[`filepath1:"github.com/g/h"`], ShouldBeTrue)
			So(outputImportSpecStrings[`filepath2:"github.com/b/c"`], ShouldBeTrue)
			So(outputImportSpecStrings[`filepath2:"github.com/c/d"`], ShouldBeTrue)
		})
		Convey("Vendored imports should be ignored", func() {
			So(outputImportSpecStrings[`filepath2:"github.com/a/b"`], ShouldBeFalse)
			So(outputImportSpecStrings[`filepath3:"github.com/x/y"`], ShouldBeFalse)
		})
		Convey("Non-github imports shouldn't be ignored", func() {
			So(outputImportSpecStrings[`filepath3:"ignore.me/r/s"`], ShouldBeTrue)
		})
		Convey("The import counts should be correct", func() {
			So(importCounts.importCountOf("filepath1"), ShouldEqual, 3)
			So(importCounts.importCountOf("filepath2"), ShouldEqual, 2)
			So(importCounts.importCountOf("filepath3"), ShouldEqual, 1)
		})
	})
}

func generateTestPackageSpec(filePath string, startIndex int) *packageSpec {
	return &packageSpec{
		startIndex: startIndex,
		filePath:   filePath,
	}
}

func generateTestImportSpec(filePath, importPath string) *importSpec {
	return &importSpec{
		imports: &ast.ImportSpec{
			Path: &ast.BasicLit{
				Kind:  token.STRING,
				Value: importPath,
			},
		},
		filePath: filePath,
	}
}
