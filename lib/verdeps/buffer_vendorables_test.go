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
		inputImportSpecChan <- generateTestImportSpec("filepath3", `"github.com/x/y"`)

		// Prepare to read outputs.
		close(inputImportSpecChan)
		close(inputPackageSpecChan)

		var (
			outputImportSpecStrings  = make(map[string]bool)
			outputPackageSpecStrings = make(map[string]bool)
		)

		// Read the import specs.
		for spec := range outputImportSpecChan {
			outputImportSpecStrings[spec.filePath+":"+spec.imports.Path.Value] = true
		}
		for spec := range outputPackageSpecChan {
			outputPackageSpecStrings[spec.filePath+":"+strconv.Itoa(spec.startIndex)] = true
		}

		// Make sure everything lines up as expected.
		// TODO(skeswa): continue here.
		Convey("All the right import specs should come through", func() {
			So(outputImportSpecStrings[`filepath1:"github.com/d/e"`], ShouldBeTrue)
			So(outputImportSpecStrings[`filepath1:"github.com/f/g"`], ShouldBeTrue)
			So(outputImportSpecStrings[`filepath1:"github.com/f/g"`], ShouldBeTrue)
		})
		Convey("The import counts should be correct", func() {

		})
		Convey("The import counts should be correct", func() {

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
