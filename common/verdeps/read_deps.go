package verdeps

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type readDepsArgs struct {
	outputChan         chan *importSpec
	packagePath        string
	accumulatedErrors  *syncedErrors
	syncedImportCounts *syncedImportCounts
}

func readDeps(args readDepsArgs) {
	var (
		err             error
		fullPackagePath string
	)

	if fullPackagePath, err = filepath.Abs(args.packagePath); err != nil {
		args.accumulatedErrors.add(err)
		return
	}

	// Make sure that the output chan is closed when there are no more files to
	// walk.
	defer close(args.outputChan)

	if err = filepath.Walk(fullPackagePath, func(path string, info os.FileInfo, err error) error {
		var f *ast.File

		if err != nil {
			return err
		}

		// Check that the current visited node on the walk is a go source file.
		if !info.IsDir() && strings.HasSuffix(info.Name(), goFileSuffix) {
			if f, err = parser.ParseFile(token.NewFileSet(), path, nil, parser.ImportsOnly); err != nil {
				return err
			}

			// Filter the deps.
			var filteredSpecs []*importSpec
			for _, spec := range f.Imports {
				// Only pursue a dependency if it has a github prefix.
				if strings.HasPrefix(spec.Path.Value, githubPrefix) {
					filteredSpecs = append(filteredSpecs, &importSpec{
						imports:  spec,
						filePath: path,
					})
					fmt.Printf("\n<DEP FIND> [%s %d:%d]:\n\t%s\n\n", path, spec.Path.Pos(), spec.Path.End(), spec.Path.Value)
				}
			}

			// Provided that we have specs, ship them to the next stage.
			if len(filteredSpecs) > 0 {
				// Set the import count before enqueing deps.
				args.syncedImportCounts.setImportCount(path, len(filteredSpecs))

				// Enqueue the deps.
				for _, spec := range filteredSpecs {
					args.outputChan <- spec
				}
			}
		}

		return nil
	}); err != nil {
		args.accumulatedErrors.add(err)
		return
	}
}
