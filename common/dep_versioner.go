package common

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"strings"
)

const goFileSuffix = ".go"

func readDeps(packagePath string, deps chan *ast.ImportSpec) error {
	var (
		err             error
		fullPackagePath string
	)

	if fullPackagePath, err = filepath.Abs(packagePath); err != nil {
		return err
	}

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

			for _, i := range f.Imports {
				deps <- i
			}
		}

		return nil
	}); err != nil {
		return err
	}

	return nil
}
