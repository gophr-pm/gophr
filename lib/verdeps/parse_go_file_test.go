package verdeps

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strconv"
	"sync"
	"testing"

	"github.com/gophr-pm/gophr/lib/errors"
	. "github.com/smartystreets/goconvey/convey"
)

func TestParseGoFile(t *testing.T) {
	Convey("Given a filepath to a go source file", t, func() {
		fakeGoFilePath := "/go/src/github.com/a/b/lib.go"

		Convey("If the go file parser returns an error, it should be added to the errors list", func() {
			var (
				errs            = errors.NewSyncedErrors()
				parse           goFileParser
				filePath        = fakeGoFilePath
				waitGroup       = sync.WaitGroup{}
				importCounts    = newSyncedImportCounts()
				expectedError   = errors.New("This is an error")
				vendorContext   = newVendorContext()
				importSpecChan  = make(chan *importSpec)
				packageSpecChan = make(chan *packageSpec)

				parseArgSrc      interface{}
				parseArgMode     parser.Mode
				parseArgFset     *token.FileSet
				parseArgFilename string
			)

			// Create a mock go file parser.
			parse = func(
				fset *token.FileSet,
				filename string,
				src interface{},
				mode parser.Mode,
			) (f *ast.File, err error) {
				parseArgSrc = src
				parseArgMode = mode
				parseArgFset = fset
				parseArgFilename = filename

				// Return an error for the purposes of this test.
				return nil, expectedError
			}

			waitGroup.Add(1)
			go parseGoFile(parseGoFileArgs{
				parse:           parse,
				errors:          errs,
				filePath:        filePath,
				waitGroup:       &waitGroup,
				importCounts:    importCounts,
				vendorContext:   vendorContext,
				importSpecChan:  importSpecChan,
				packageSpecChan: packageSpecChan,
			})

			// Wait for the function to exit before asserting anything.
			waitGroup.Wait()

			// Make sure that channels get closed.
			close(importSpecChan)
			close(packageSpecChan)

			So(errs.Len(), ShouldEqual, 1)
			So(errs.Get()[0].Error(), ShouldEqual, expectedError.Error())
			So(parseArgSrc, ShouldBeNil)
			So(parseArgMode, ShouldEqual, parser.ImportsOnly)
			So(parseArgFset, ShouldNotBeNil)
			So(parseArgFilename, ShouldEqual, filePath)
		})

		Convey("When import specs aren't vendored, then none should be ignored", func() {
			var (
				errs            = errors.NewSyncedErrors()
				parse           goFileParser
				filePath        = fakeGoFilePath
				waitGroup       = sync.WaitGroup{}
				importCounts    = newSyncedImportCounts()
				vendorContext   = newVendorContext()
				importSpecChan  = make(chan *importSpec, 3)
				packageSpecChan = make(chan *packageSpec, 1)

				parseArgSrc      interface{}
				parseArgMode     parser.Mode
				parseArgFset     *token.FileSet
				parseArgFilename string
			)

			// Create a mock go file parser.
			parse = func(
				fset *token.FileSet,
				filename string,
				src interface{},
				mode parser.Mode,
			) (f *ast.File, err error) {
				parseArgSrc = src
				parseArgMode = mode
				parseArgFset = fset
				parseArgFilename = filename

				// Return an error for the purposes of this test.
				return &ast.File{
					Imports: []*ast.ImportSpec{
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								ValuePos: token.Pos(1),
								Kind:     token.STRING,
								Value:    `"github.com/e/f/g"`,
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								ValuePos: token.Pos(2),
								Kind:     token.STRING,
								Value:    `"github.com/h/i"`,
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								ValuePos: token.Pos(3),
								Kind:     token.STRING,
								Value:    `"github.com/j/k/l/m"`,
							},
						},
					},
					Package: token.Pos(42),
				}, nil
			}

			waitGroup.Add(1)
			go parseGoFile(parseGoFileArgs{
				parse:           parse,
				errors:          errs,
				filePath:        filePath,
				waitGroup:       &waitGroup,
				importCounts:    importCounts,
				vendorContext:   vendorContext,
				importSpecChan:  importSpecChan,
				packageSpecChan: packageSpecChan,
			})

			// Wait for the function to exit before asserting anything.
			waitGroup.Wait()

			// Read all the import specs into a map for easy assertion.
			importSpecStrings := make(map[string]bool)
			for spec := range importSpecChan {
				importSpecStrings[spec.filePath+`:`+spec.imports.Path.Value] = true

				// Close the channel when out of data.
				if len(importSpecChan) == 0 {
					close(importSpecChan)
				}
			}

			// Read all the package specs for easy assertion.
			packageSpecStrings := make(map[string]bool)
			for spec := range packageSpecChan {
				packageSpecStrings[spec.filePath+`:`+strconv.Itoa(spec.startIndex)] = true

				// Close the channel when out of data.
				if len(packageSpecChan) == 0 {
					close(packageSpecChan)
				}
			}

			So(errs.Len(), ShouldEqual, 0)
			So(parseArgSrc, ShouldBeNil)
			So(parseArgMode, ShouldEqual, parser.ImportsOnly)
			So(parseArgFset, ShouldNotBeNil)
			So(parseArgFilename, ShouldEqual, filePath)
			So(packageSpecStrings[filePath+`:42`], ShouldBeTrue)
			So(importCounts.importCountOf(filePath), ShouldEqual, 3)
			So(importSpecStrings[filePath+`:"github.com/e/f/g"`], ShouldBeTrue)
			So(importSpecStrings[filePath+`:"github.com/h/i"`], ShouldBeTrue)
			So(importSpecStrings[filePath+`:"github.com/j/k/l/m"`], ShouldBeTrue)
		})

		Convey("When import specs are vendored, they should be ignored", func() {
			var (
				errs            = errors.NewSyncedErrors()
				parse           goFileParser
				filePath        = fakeGoFilePath
				waitGroup       = sync.WaitGroup{}
				importCounts    = newSyncedImportCounts()
				vendorContext   = newVendorContext()
				importSpecChan  = make(chan *importSpec, 3)
				packageSpecChan = make(chan *packageSpec, 1)

				parseArgSrc      interface{}
				parseArgMode     parser.Mode
				parseArgFset     *token.FileSet
				parseArgFilename string
			)

			// Add to the vendor context.
			vendorContext.add("github.com/e/f/g")
			vendorContext.add("github.com/j/k/l/m")

			// Create a mock go file parser.
			parse = func(
				fset *token.FileSet,
				filename string,
				src interface{},
				mode parser.Mode,
			) (f *ast.File, err error) {
				parseArgSrc = src
				parseArgMode = mode
				parseArgFset = fset
				parseArgFilename = filename

				// Return an error for the purposes of this test.
				return &ast.File{
					Imports: []*ast.ImportSpec{
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								ValuePos: token.Pos(1),
								Kind:     token.STRING,
								Value:    `"github.com/e/f/g"`,
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								ValuePos: token.Pos(2),
								Kind:     token.STRING,
								Value:    `"github.com/h/i"`,
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								ValuePos: token.Pos(3),
								Kind:     token.STRING,
								Value:    `"github.com/j/k/l/m"`,
							},
						},
					},
					Package: token.Pos(42),
				}, nil
			}

			waitGroup.Add(1)
			go parseGoFile(parseGoFileArgs{
				parse:           parse,
				errors:          errs,
				filePath:        filePath,
				waitGroup:       &waitGroup,
				importCounts:    importCounts,
				vendorContext:   vendorContext,
				importSpecChan:  importSpecChan,
				packageSpecChan: packageSpecChan,
			})

			// Wait for the function to exit before asserting anything.
			waitGroup.Wait()

			// Read all the import specs into a map for easy assertion.
			importSpecStrings := make(map[string]bool)
			for spec := range importSpecChan {
				importSpecStrings[spec.filePath+`:`+spec.imports.Path.Value] = true

				// Close the channel when out of data.
				if len(importSpecChan) == 0 {
					close(importSpecChan)
				}
			}

			// Read all the package specs for easy assertion.
			packageSpecStrings := make(map[string]bool)
			for spec := range packageSpecChan {
				packageSpecStrings[spec.filePath+`:`+strconv.Itoa(spec.startIndex)] = true

				// Close the channel when out of data.
				if len(packageSpecChan) == 0 {
					close(packageSpecChan)
				}
			}

			So(errs.Len(), ShouldEqual, 0)
			So(parseArgSrc, ShouldBeNil)
			So(parseArgMode, ShouldEqual, parser.ImportsOnly)
			So(parseArgFset, ShouldNotBeNil)
			So(parseArgFilename, ShouldEqual, filePath)
			So(packageSpecStrings[filePath+`:42`], ShouldBeTrue)
			So(importCounts.importCountOf(filePath), ShouldEqual, 1)
			So(importSpecStrings[filePath+`:"github.com/e/f/g"`], ShouldBeFalse)
			So(importSpecStrings[filePath+`:"github.com/h/i"`], ShouldBeTrue)
			So(importSpecStrings[filePath+`:"github.com/j/k/l/m"`], ShouldBeFalse)
		})
	})
}
