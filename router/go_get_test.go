package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGenerateGophrURL(t *testing.T) {
	assert.Equal(t, "a.b/c/d", generateGophrURL("a.b", "c", "d", ""))
	assert.Equal(t, "a.b/c/d@e", generateGophrURL("a.b", "c", "d", "e"))
}

func TestGenerateGithubTreeURLTemplate(t *testing.T) {
	assert.Equal(t, "https://github.com/a/b/tree/c{/dir}", generateGithubTreeURLTemplate("a", "b", "c"))
	assert.Equal(t, "https://github.com/a/b/tree/master{/dir}", generateGithubTreeURLTemplate("a", "b", ""))
}

func TestGenerateDepotBlobURLTemplate(t *testing.T) {
	assert.Equal(t, "https://a/blob/b/c/d{/dir}/{file}#L{line}", generateDepotBlobURLTemplate("a", "b", "c", "d"))
	assert.Equal(t, "https://a/blob/b/c/master{/dir}/{file}#L{line}", generateDepotBlobURLTemplate("a", "b", "c", ""))
}

func TestGenerateGoGetMetadata(t *testing.T) {
	assert.Equal(t, `
<html>
<head>
<meta name="go-import" content=" git https://">
<meta name="go-source" content=" _  ">
</head>
<body>
go get 
</body>
</html>
`, generateGoGetMetadata(generateGoGetMetadataArgs{}))

	assert.Equal(t, `
<html>
<head>
<meta name="go-import" content="a git https://a">
<meta name="go-source" content="a _ b c">
</head>
<body>
go get a
</body>
</html>
`, generateGoGetMetadata(generateGoGetMetadataArgs{
		gophrURL:        "a",
		treeURLTemplate: "b",
		blobURLTemplate: "c"}))
}
