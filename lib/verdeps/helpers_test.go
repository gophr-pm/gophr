package verdeps

import (
	"os"
	"regexp"
	"testing"

	"github.com/gophr-pm/gophr/lib/io"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

var validInternalDirnameRegex = regexp.MustCompile(`\b[0-9a-f]{16}\b`)
var generatedInternalDirName = generateInternalDirName()

const author = "raymondChandler"
const repo = "theLongGoodBye"
const subpath = "/goodRead"
const sixCharSha = "abcdef"
const packageDirPath = "gophr"

func TestIsSubPackage(t *testing.T) {
	isSub := isSubPackage("skeswa", "skeswa", "gophr", "gophr")
	assert.True(t, isSub)

	isSub = isSubPackage("skeswa", "shikkic", "gophr", "gophr")
	assert.False(t, isSub)

	isSub = isSubPackage("skeswa", "skeswa", "gophr", "gophr-pm")
	assert.False(t, isSub)

	isSub = isSubPackage("skeswa", "shikkic", "gophr", "gophr-pm")
	assert.False(t, isSub)
}

func TestParseImportPath_validAuthor(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := author, "", ""
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + author + "/")
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestParseImportPath_invalidAuthor(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := "IdontEndInAForwardSlashSoIWillBeMissingACha", "", ""
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + "IdontEndInAForwardSlashSoIWillBeMissingAChar")
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestParseImportPath_validRepo(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := author, repo, ""
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + author + "/" + repo)
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestParseImportPath_validSubpath(t *testing.T) {
	t.Parallel()
	expectedAuthor, expectedRepo, expectedSubpath := author, repo, subpath
	actualAuthor, actualRepo, actualSubpath := parseImportPath(githubPrefix + author + "/" + repo + subpath)
	assert.Equal(t, expectedAuthor, actualAuthor)
	assert.Equal(t, expectedRepo, actualRepo)
	assert.Equal(t, expectedSubpath, actualSubpath)
}

func TestGenerateInternalDirName_hasProperLengthAndAcceptedCharacters(t *testing.T) {
	t.Parallel()
	for i := 0; i < 10; i++ {
		internalDirName := generateInternalDirName()
		assert.True(t, validInternalDirnameRegex.MatchString(internalDirName))
	}
}

func TestComposeNewImportPath_noSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, "", generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestComposeNewImportPath_withInternalSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + "/" + generatedInternalDirName + "/" + "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, internalSubPathPart, generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestComposeNewImportPath_withInternalSuffixSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + "/" + generatedInternalDirName + "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, internalSubPathSuffix, generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestComposeNewImportPath_withRegularSubpath(t *testing.T) {
	t.Parallel()
	expectedComposedPath := gophrPrefix + author + "/" + repo + "@" + sixCharSha + subpath + "\""
	actualComposedPath := string(composeNewImportPath(author, repo, sixCharSha, subpath, generatedInternalDirName)[:])
	assert.Equal(t, expectedComposedPath, actualComposedPath)
}

func TestGetPackageDirPaths_onlyVendorDir(t *testing.T) {
	t.Parallel()
	var mockFiles []os.FileInfo
	var expectedSubDirs []string

	mockVendorFile := io.MockFileInfo{NameProp: "vendor", IsDirProp: true}
	mockIo := io.NewMockIO()
	mockIo.On("Stat", mock.AnythingOfType("string")).Return(mockVendorFile, nil)

	verdepsHelperArgs := verdepsHelperArgs{io: mockIo}
	mockFiles = append(mockFiles, mockVendorFile)

	vendorDirPath, subDirNames, goFilePaths := verdepsHelperArgs.getPackageDirPaths(mockFiles, packageDirPath)
	assert.Equal(t, packageDirPath+"/vendor/src", vendorDirPath)
	assert.Equal(t, expectedSubDirs, subDirNames)
	assert.Equal(t, expectedSubDirs, goFilePaths)
}

func TestGetPackageDirPaths_vendorDirWithSrcDir(t *testing.T) {
	t.Parallel()
	var mockFiles []os.FileInfo
	var expectedSubDirs []string

	mockVendorFile := io.MockFileInfo{NameProp: "vendor", IsDirProp: true}
	mockSrcFile := io.MockFileInfo{NameProp: "uselessfile", IsDirProp: false}
	mockIo := io.NewMockIO()
	mockIo.On("Stat", mock.AnythingOfType("string")).Return(mockSrcFile, nil)

	verdepsHelperArgs := verdepsHelperArgs{io: mockIo}
	mockFiles = append(mockFiles, mockVendorFile, mockSrcFile)

	vendorDirPath, subDirNames, goFilePaths := verdepsHelperArgs.getPackageDirPaths(mockFiles, packageDirPath)
	assert.Equal(t, packageDirPath+"/vendor", vendorDirPath)
	assert.Equal(t, expectedSubDirs, subDirNames)
	assert.Equal(t, expectedSubDirs, goFilePaths)
}

func TestGetPackageDirPaths_vendorDirAndGoFiles(t *testing.T) {
	t.Parallel()
	var mockFiles []os.FileInfo
	var expectedSubDirs []string
	var expectedGoFileNames []string

	expectedGoFiles := makeRandomMockFiles(5, ".go", false)
	for _, file := range expectedGoFiles {
		expectedGoFileNames = append(expectedGoFileNames, packageDirPath+"/"+file.Name())
		mockFiles = append(mockFiles, file)
	}

	mockVendorFile := io.MockFileInfo{NameProp: "vendor", IsDirProp: true}
	mockIo := io.NewMockIO()
	mockIo.On("Stat", mock.AnythingOfType("string")).Return(mockVendorFile, nil)

	verdepsHelperArgs := verdepsHelperArgs{io: mockIo}
	mockFiles = append(mockFiles, mockVendorFile)

	vendorDirPath, subDirNames, goFilePaths := verdepsHelperArgs.getPackageDirPaths(mockFiles, packageDirPath)

	assert.Equal(t, packageDirPath+"/vendor/src", vendorDirPath)
	assert.Equal(t, expectedSubDirs, subDirNames)
	assert.Equal(t, expectedGoFileNames, goFilePaths)
}

func TestGetPackageDirPaths_vendorDirAndGoFilesAndSubdirs(t *testing.T) {
	t.Parallel()
	var mockFiles []os.FileInfo
	var expectedSubDirNames []string
	var expectedGoFileNames []string

	expectedGoFiles := makeRandomMockFiles(5, ".go", false)
	for _, file := range expectedGoFiles {
		expectedGoFileNames = append(expectedGoFileNames, packageDirPath+"/"+file.Name())
		mockFiles = append(mockFiles, file)
	}

	expectedSubDirs := makeRandomMockFiles(5, "", true)
	for _, file := range expectedSubDirs {
		expectedSubDirNames = append(expectedSubDirNames, file.Name())
		mockFiles = append(mockFiles, file)
	}

	mockVendorFile := io.MockFileInfo{NameProp: "vendor", IsDirProp: true}
	mockIo := io.NewMockIO()
	mockIo.On("Stat", mock.AnythingOfType("string")).Return(mockVendorFile, nil)

	verdepsHelperArgs := verdepsHelperArgs{io: mockIo}
	mockFiles = append(mockFiles, mockVendorFile)

	vendorDirPath, subDirNames, goFilePaths := verdepsHelperArgs.getPackageDirPaths(mockFiles, packageDirPath)

	assert.Equal(t, packageDirPath+"/vendor/src", vendorDirPath)
	assert.Equal(t, expectedSubDirNames, subDirNames)
	assert.Equal(t, expectedGoFileNames, goFilePaths)
}

func makeRandomMockFiles(numberOfFiles int, extension string, isDir bool) []os.FileInfo {
	var mockFileInfos []os.FileInfo
	for x := 0; x < numberOfFiles; x++ {
		// conveniently, generateInternalDirName is a random string generator
		mockFileInfos = append(mockFileInfos, io.MockFileInfo{NameProp: generateInternalDirName() + extension, IsDirProp: isDir})
	}
	return mockFileInfos
}
