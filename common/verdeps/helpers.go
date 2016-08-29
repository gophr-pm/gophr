package verdeps

import "bytes"

func isSubPackage(depAuthor, depRepo, packageAuthor, packageRepo string) bool {
	return depAuthor == packageAuthor && depRepo == packageRepo
}

func parseImportPath(importPath string) (author string, repo string, subpath string) {
	var (
		i                 int
		repoStartIndex    int
		subpathStartIndex int

		importPathLength = len(importPath)
		authorStartIndex = len(githubPrefix)
	)

	// Advance to the next slash.
	for i = authorStartIndex; i < importPathLength && importPath[i] != '/'; i++ {
	}

	// Exit if we reached the end of the import path.
	if i == importPathLength {
		return importPath[authorStartIndex : importPathLength-1], "", ""
	}

	author = importPath[authorStartIndex:i]
	repoStartIndex = i + 1

	// Advance past the current slash to the next one (or the end of the string).
	for i = repoStartIndex; i < importPathLength && importPath[i] != '/' && importPath[i] != '"'; i++ {
	}

	repo = importPath[repoStartIndex:i]
	subpathStartIndex = i

	// Advance past the current slash to the end of the string.
	for i = subpathStartIndex; i < importPathLength && importPath[i] != '"'; i++ {
	}

	// Set the subpath if there was one.
	if i > subpathStartIndex {
		subpath = importPath[subpathStartIndex:i]
	}

	return author, repo, subpath
}

func composeNewImportPath(author, repo, sha, subpath string) []byte {
	var buffer bytes.Buffer
	buffer.WriteString(gophrPrefix)
	buffer.WriteString(author)
	buffer.WriteByte('/')
	buffer.WriteString(repo)
	buffer.WriteByte('@')
	buffer.WriteString(sha)

	if len(subpath) > 0 {
		buffer.WriteString(subpath)
	}

	buffer.WriteByte('"')

	return buffer.Bytes()
}

func importPathHashOf(importPath string) string {
	author, repo, _ := parseImportPath(importPath)

	buffer := bytes.Buffer{}
	buffer.WriteString(author)
	buffer.WriteByte('/')
	buffer.WriteString(repo)

	return buffer.String()
}
