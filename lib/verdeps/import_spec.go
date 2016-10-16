package verdeps

import (
	"go/ast"
)

type importSpec struct {
	imports  *ast.ImportSpec
	filePath string
}
