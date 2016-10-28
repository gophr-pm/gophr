package verdeps

import (
	"errors"
	"fmt"
	stdio "io"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"testing"

	"github.com/gophr-pm/gophr/lib/io"
	. "github.com/smartystreets/goconvey/convey"
)

func TestTraversePackageDir(t *testing.T) {
	Convey("Given a package dir to traverse", t, func() {
		Convey("All go files should be parsed in order", func() {
			var (
				fio                      = &fakeIO{root: generateTestFSTree()}
				errs                     = newSyncedErrors()
				parse                    goFileASTParser
				dirPath                  = "root"
				waitGroup                = &sync.WaitGroup{}
				subDirPath               = ""
				inVendorDir              = false
				importCounts             = newSyncedImportCounts()
				vendorContext            = newVendorContext()
				parsedGoFiles            = make(map[string]bool)
				importSpecChan           = make(chan *importSpec)
				packageSpecChan          = make(chan *packageSpec)
				parseGoFileArgsList      []parseGoFileArgs
				parseGoFileArgsListLock  sync.Mutex
				generatedInternalDirName = "SOME_GENERATED_DIR_NAME"
			)

			parse = func(args parseGoFileArgs) {
				parseGoFileArgsListLock.Lock()
				parsedGoFiles[args.filePath] = true
				parseGoFileArgsList = append(parseGoFileArgsList, args)
				parseGoFileArgsListLock.Unlock()

				args.waitGroup.Done()
			}

			waitGroup.Add(1)
			traversePackageDir(traversePackageDirArgs{
				io:                       fio,
				errors:                   errs,
				dirPath:                  dirPath,
				waitGroup:                waitGroup,
				subDirPath:               subDirPath,
				parseGoFile:              parse,
				inVendorDir:              inVendorDir,
				importCounts:             importCounts,
				vendorContext:            vendorContext,
				importSpecChan:           importSpecChan,
				packageSpecChan:          packageSpecChan,
				generatedInternalDirName: generatedInternalDirName,
			})
			waitGroup.Wait()

			for _, arg := range parseGoFileArgsList {
				So(arg.errors, ShouldNotBeNil)
				So(arg.filePath, ShouldNotBeEmpty)
				So(arg.importSpecChan, ShouldNotBeNil)
				So(arg.packageSpecChan, ShouldNotBeNil)
				So(arg.parse, ShouldNotBeNil)
				So(arg.vendorContext, ShouldNotBeNil)
				So(arg.waitGroup, ShouldNotBeNil)
			}

			So(errs.len(), ShouldEqual, 0)
			So(parsedGoFiles[`root/vendor/golang.org/x/snapper/tastes_good/fish.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/cli/command_c.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/cli/command_a.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/cli/command_b.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/cli/SOME_GENERATED_DIR_NAME/helpers/the_help.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/cli/main.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/vendor/github.com/x/y/y.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/vendor/github.com/c/d/e/e.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/vendor/github.com/a/b/b.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/bar.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/foo.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/SOME_GENERATED_DIR_NAME/debug.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/SOME_GENERATED_DIR_NAME/errors.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/SOME_GENERATED_DIR_NAME/helpers/delete.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/SOME_GENERATED_DIR_NAME/helpers/read.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/SOME_GENERATED_DIR_NAME/helpers/create.go`], ShouldBeTrue)
			So(parsedGoFiles[`root/lib/SOME_GENERATED_DIR_NAME/helpers/update.go`], ShouldBeTrue)
		})

		Convey("If a rename operation fails, it should bubble up", func() {
			var (
				fio = &fakeIO{
					root:               generateTestFSTree(),
					forceRenameFailure: true,
				}

				errs                     = newSyncedErrors()
				parse                    goFileASTParser
				dirPath                  = "root"
				waitGroup                = &sync.WaitGroup{}
				subDirPath               = ""
				inVendorDir              = false
				importCounts             = newSyncedImportCounts()
				vendorContext            = newVendorContext()
				importSpecChan           = make(chan *importSpec)
				packageSpecChan          = make(chan *packageSpec)
				generatedInternalDirName = "SOME_GENERATED_DIR_NAME"
			)

			parse = func(args parseGoFileArgs) {
				args.waitGroup.Done()
			}

			waitGroup.Add(1)
			traversePackageDir(traversePackageDirArgs{
				io:                       fio,
				errors:                   errs,
				dirPath:                  dirPath,
				waitGroup:                waitGroup,
				subDirPath:               subDirPath,
				parseGoFile:              parse,
				inVendorDir:              inVendorDir,
				importCounts:             importCounts,
				vendorContext:            vendorContext,
				importSpecChan:           importSpecChan,
				packageSpecChan:          packageSpecChan,
				generatedInternalDirName: generatedInternalDirName,
			})
			waitGroup.Wait()

			// There are two renames that should fail.
			So(errs.len(), ShouldEqual, 2)
		})

		Convey("If a directory read operation fails, it should bubble up", func() {
			var (
				fio = &fakeIO{
					root:                generateTestFSTree(),
					forceReadDirFailure: true,
				}

				errs                     = newSyncedErrors()
				parse                    goFileASTParser
				dirPath                  = "root"
				waitGroup                = &sync.WaitGroup{}
				subDirPath               = ""
				inVendorDir              = false
				importCounts             = newSyncedImportCounts()
				vendorContext            = newVendorContext()
				importSpecChan           = make(chan *importSpec)
				packageSpecChan          = make(chan *packageSpec)
				generatedInternalDirName = "SOME_GENERATED_DIR_NAME"
			)

			parse = func(args parseGoFileArgs) {
				args.waitGroup.Done()
			}

			waitGroup.Add(1)
			traversePackageDir(traversePackageDirArgs{
				io:                       fio,
				errors:                   errs,
				dirPath:                  dirPath,
				waitGroup:                waitGroup,
				subDirPath:               subDirPath,
				parseGoFile:              parse,
				inVendorDir:              inVendorDir,
				importCounts:             importCounts,
				vendorContext:            vendorContext,
				importSpecChan:           importSpecChan,
				packageSpecChan:          packageSpecChan,
				generatedInternalDirName: generatedInternalDirName,
			})
			waitGroup.Wait()

			So(errs.len(), ShouldEqual, 1)
		})
	})
}

func generateTestFSTree() *fakeFSNode {
	return &fakeFSNode{
		name: "root",
		children: []*fakeFSNode{
			&fakeFSNode{
				name: "lib",
				children: []*fakeFSNode{
					&fakeFSNode{
						name: "internal",
						children: []*fakeFSNode{
							&fakeFSNode{
								name: "helpers",
								children: []*fakeFSNode{
									&fakeFSNode{name: "create.go"},
									&fakeFSNode{name: "read.go"},
									&fakeFSNode{name: "update.go"},
									&fakeFSNode{name: "delete.go"},
								},
							},
							&fakeFSNode{name: "errors.go"},
							&fakeFSNode{name: "debug.go"},
						},
					},
					&fakeFSNode{name: "foo.go"},
					&fakeFSNode{name: "bar.go"},
					&fakeFSNode{
						name: "vendor",
						children: []*fakeFSNode{
							&fakeFSNode{
								name: "github.com",
								children: []*fakeFSNode{
									&fakeFSNode{
										name: "a",
										children: []*fakeFSNode{
											&fakeFSNode{
												name: "b",
												children: []*fakeFSNode{
													&fakeFSNode{name: "b.go"},
												},
											},
										},
									},
									&fakeFSNode{
										name: "c",
										children: []*fakeFSNode{
											&fakeFSNode{
												name: "d",
												children: []*fakeFSNode{
													&fakeFSNode{
														name: "e",
														children: []*fakeFSNode{
															&fakeFSNode{name: "e.go"},
														},
													},
												},
											},
										},
									},
									&fakeFSNode{
										name: "x",
										children: []*fakeFSNode{
											&fakeFSNode{
												name: "y",
												children: []*fakeFSNode{
													&fakeFSNode{name: "y.go"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			{
				name: "cli",
				children: []*fakeFSNode{
					{name: "main.go"},
					{name: "command_a.go"},
					{name: "command_b.go"},
					{name: "command_c.go"},
					{
						name: "internal",
						children: []*fakeFSNode{
							{
								name: "helpers",
								children: []*fakeFSNode{
									{name: "the_help.go"},
								},
							},
						},
					},
				},
			},
			{
				name: "vendor",
				children: []*fakeFSNode{
					{
						name: "golang.org",
						children: []*fakeFSNode{
							{
								name: "x",
								children: []*fakeFSNode{
									{
										name: "snapper",
										children: []*fakeFSNode{
											{
												name: "tastes_good",
												children: []*fakeFSNode{
													{name: "fish.go"},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

type fakeFSNode struct {
	name     string
	children []*fakeFSNode
}

func (fsn fakeFSNode) Info() os.FileInfo {
	return io.NewFakeFileInfo(
		fsn.name,
		42,
		!strings.HasSuffix(fsn.name, ".go"))
}

type fakeIO struct {
	root                *fakeFSNode
	lock                sync.RWMutex
	forceRenameFailure  bool
	forceReadDirFailure bool
}

func (fio *fakeIO) getNodeAtPath(path string) (*fakeFSNode, *fakeFSNode, bool) {
	var (
		parts    = strings.Split(path, fmt.Sprintf("%c", os.PathSeparator))
		prevNode *fakeFSNode
		currNode = fio.root
	)

	// Remove "root" from the beginning.
	if len(parts) > 0 && parts[0] == "root" {
		parts = parts[1:]
	}

	// Try to match each part of the path.
	for _, part := range parts {
		nextNodeFound := false

		// Look for the next path part amongst the children.
		for _, child := range currNode.children {
			if child.name == part {
				prevNode = currNode
				currNode = child
				nextNodeFound = true
				break
			}
		}

		// If nothing matched, exit.
		if !nextNodeFound {
			return nil, nil, false
		}
	}

	return prevNode, currNode, true
}
func (fio *fakeIO) Mkdir(name string, perm os.FileMode) error {
	return nil
}
func (fio *fakeIO) Create(name string) (*os.File, error) {
	return nil, nil
}
func (fio *fakeIO) Copy(dst stdio.Writer, src stdio.Reader) (int64, error) {
	return int64(0), nil
}
func (fio *fakeIO) ReadDir(dirname string) ([]os.FileInfo, error) {
	fio.lock.RLock()
	defer fio.lock.RUnlock()

	if fio.forceReadDirFailure {
		return nil, errors.New("this is supposed to fail")
	}

	_, node, exists := fio.getNodeAtPath(dirname)
	if !exists {
		return nil, errors.New("no such path: " + dirname)
	}

	var children []os.FileInfo
	for _, child := range node.children {
		children = append(children, child.Info())
	}

	return children, nil
}
func (fio *fakeIO) Stat(dirname string) (os.FileInfo, error) {
	fio.lock.RLock()
	defer fio.lock.RUnlock()

	_, node, exists := fio.getNodeAtPath(dirname)
	if !exists {
		return nil, errors.New("no such path: " + dirname)
	}

	return node.Info(), nil
}
func (fio *fakeIO) ReadFile(filename string) ([]byte, error) {
	return nil, nil
}
func (fio *fakeIO) WriteFile(
	filename string,
	data []byte,
	perm os.FileMode,
) error {
	return nil
}
func (fio *fakeIO) Rename(oldPath, newPath string) error {
	fio.lock.Lock()
	defer fio.lock.Unlock()

	if fio.forceRenameFailure {
		return errors.New("this is supposed to fail")
	}

	oldPathDir, _ := filepath.Split(oldPath)
	newPathDir, newPathName := filepath.Split(newPath)

	oldParentNode, oldNode, exists := fio.getNodeAtPath(oldPath)
	if !exists {
		return fmt.Errorf(`oldPath "%s" does not exist`, oldPath)
	}

	newParentNode, _, exists := fio.getNodeAtPath(newPath)
	if exists {
		return fmt.Errorf(`newPath "%s" (for oldPath "%s") already exists`, newPath, oldPath)
	}

	if oldPathDir == newPathDir {
		oldNode.name = newPathName
	} else {
		newOldParentNodeChildren := []*fakeFSNode{}
		for _, childNode := range oldParentNode.children {
			if childNode != oldNode {
				newOldParentNodeChildren = append(newOldParentNodeChildren, childNode)
			}
		}
		oldParentNode.children = newOldParentNodeChildren
		newParentNode.children = append(newParentNode.children, oldNode)
	}

	return nil
}
