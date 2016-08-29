package verdeps

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

type readDepsArgs struct {
	outputChan        chan *importSpec
	packagePath       string
	accumulatedErrors *syncedErrors
}

// TODO(skeswa): instead of saving a file path, we need to calculate the position using the file set. https://golang.org/pkg/go/token/#Position

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

			for _, spec := range f.Imports {
				// Only pursue a dependency if it has a github prefix.
				if strings.HasPrefix(spec.Path.Value, githubPrefix) {
					args.outputChan <- &importSpec{imports: spec, filePath: path}
				}
			}
		}

		return nil
	}); err != nil {
		args.accumulatedErrors.add(err)
		return
	}
}
