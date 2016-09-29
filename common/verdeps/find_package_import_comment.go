package verdeps

import "regexp"

const (
	importAnnotation = `import\s+(?:"[^"]*"|` + "`[^`]*`" + `)`
	importComment    = `(?://\s*` +
		importAnnotation +
		`\s*$|/\*\s*` +
		importAnnotation +
		`\s*\*/)`
	packageImportComment = `(?:package\s+\w+)(\s+` + importComment + `(?:.*))`
)

var packageImportCommentRegex = regexp.MustCompile(packageImportComment)

// findPackageImportComment finds the indices of the package import comment if
// one exists. If not, then it returns -1s.
func findPackageImportComment(
	fileData []byte,
	packageStartIndex int,
) (fromIndex int, toIndex int) {
	// Read until the end of the line or the end of the file.
	packageEndIndex := packageStartIndex
	for packageEndIndex < len(fileData) && fileData[packageEndIndex] != '\n' {
		packageEndIndex++
	}
	// Read backwards until the beginning of the file or the previous line.
	for packageStartIndex >= 0 && fileData[packageStartIndex] != '\n' {
		packageStartIndex--
	}
	// Advance the start index by one to make sure it is securely at the
	// beginning of the line.
	packageStartIndex++

	// Find matches in [packageStartIndex, packageEndIndex).
	line := fileData[packageStartIndex:packageEndIndex]
	match := packageImportCommentRegex.FindSubmatchIndex(line)
	if match != nil {
		// Adjust the indices for the starting index of "line".
		toIndex = packageStartIndex + match[3]
		fromIndex = packageStartIndex + match[2]
		return fromIndex, toIndex
	}

	return -1, -1
}
