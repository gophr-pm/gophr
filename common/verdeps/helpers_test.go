package verdeps

import (
	"testing"
	"regexp"
	"github.com/stretchr/testify/assert"
)
/*
 Untested functions: 
 isSubPackage (too simple)
 importPathHashOf (too simple and I don't understand why it's called that)
 subDirExists (unused) 
 TODO: getPackageDirPaths have to mock OS.stat 
*/
var validInternalDirnameRegex = regexp.MustCompile(`\b[0-9a-f]{16}\b`)
var generatedInternalDirName string = generateInternalDirName();

const author = "raymondChandler"
const repo = "theLongGoodBye"
const subpath = "/goodRead"
const sixCharSha = "abcdef"


func TestParseImportPath_validAuthor(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := author, "", "";
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + author + "/")
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestParseImportPath_invalidAuthor(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := "IdontEndInAForwardSlashSoIWillBeMissingACha", "", "";
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + "IdontEndInAForwardSlashSoIWillBeMissingAChar")
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestParseImportPath_validRepo(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := author, repo, "";
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + author + "/" + repo)
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestParseImportPath_validSubpath(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := author, repo, subpath;
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + author + "/" + repo + subpath)
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestGenerateInternalDirName_hasProperLengthAndAcceptedCharacters(t *testing.T) {
	t.Parallel()
	for i := 0; i < 10; i++ {
		internalDirName := generateInternalDirName()
		assert.True(t, validInternalDirnameRegex.MatchString(internalDirName));
	}
}

func TestComposeNewImportPath_noSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha +  "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, "", generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestComposeNewImportPath_withInternalSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + "/" + generatedInternalDirName + "/" +  "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, internalSubPathPart, generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestComposeNewImportPath_withInternalSuffixSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + "/" + generatedInternalDirName +  "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, internalSubPathSuffix, generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestComposeNewImportPath_withRegularSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + subpath +  "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha,  subpath, generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}
