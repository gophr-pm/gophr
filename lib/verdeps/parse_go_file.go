package verdeps

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"sync"

	"github.com/gophr-pm/gophr/lib/errors"
)

// parseGoFileArgs is the arguments struct for parseGoFileArgs.
type parseGoFileArgs struct {
	errors          *errors.SyncedErrors
	filePath        string
	waitGroup       *sync.WaitGroup
	importCounts    *syncedImportCounts
	vendorContext   *vendorContext
	importSpecChan  chan *importSpec
	packageSpecChan chan *packageSpec
}

// parseGoFile uses the golang compiler's AST parser to parse out the import
// specs and package spec of a file written in Go. The specs are then returned
// via the import spec and package spec channels.
func parseGoFile(args parseGoFileArgs) {
	var (
		f     *ast.File
		err   error
		specs []*importSpec
	)

	// Signal that the wait group should continue when this function exits.
	defer args.waitGroup.Done()

	// Parse the imports of the file.
	if f, err = parser.ParseFile(
		token.NewFileSet(),
		args.filePath,
		nil,
		parser.ImportsOnly,
	); err != nil {
		args.errors.Add(err)
		return
	}

	// Filter the specs.
	for _, spec := range f.Imports {
		// Ignore the surrounding quotes.
		importString := strings.Trim(spec.Path.Value, "\"")

		// Only pursue a dependency if it has a github prefix and is not vendored.
		if !args.vendorContext.contains(importString) &&
			strings.HasPrefix(importString, "github.com/") {
			// Both conditions were met, so add this import spec to the list.
			specs = append(specs, &importSpec{
				imports:  spec,
				filePath: args.filePath,
			})
		}
	}

	// Set the import count for this file path.
	if args.importCounts != nil {
		args.importCounts.setImportCount(args.filePath, len(specs))
	}

	// Throw all the newly discovered specs into the mix.
	args.packageSpecChan <- &packageSpec{
		filePath:   args.filePath,
		startIndex: int(f.Package),
	}
	for _, spec := range specs {
		args.importSpecChan <- spec
	}

	return
}
