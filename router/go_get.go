package main

import (
	"bytes"
	"fmt"
)

const (
	goGetMetadataFormat = `
<html>
<head>
<meta name="go-import" content="%s git %s">
<meta name="go-source" content="%s _ %s %s">
</head>
<body>
go get %s
</body>
</html>
`
	githubTreeURLTemplate = "https://github.com/%s/%s/tree/%s{/dir}"
	depotBlobURLTemplate  = "https://%s/api/blob/%s/%s/%s{/dir}/{file}#L{line}"
)

type generateGoGetMetadataArgs struct {
	gophrURL        string // e.g. "gophr.pm/abc/wxyz@1.1+"
	depotURL        string // e.g. "https://gophr.pm/depot/3abc4wxyz-97e17db9944a97a72765fcc18a237aaa0bb200a3.git"
	treeURLTemplate string // e.g. "https://github.com/urfave/cli/tree/v1.18.1{/dir}"
	blobURLTemplate string // e.g. "https://github.com/urfave/cli/blob/v1.18.1{/dir}/{file}#L{line}"
}

// TODO(skeswa): write the formatter and comapre against gopkg.

func generateGophrURL(domain, author, repo, selector string) string {
	var buffer bytes.Buffer
	buffer.WriteString(domain)
	buffer.WriteByte('/')
	buffer.WriteString(author)
	buffer.WriteByte('/')
	buffer.WriteString(repo)

	if len(selector) > 0 {
		buffer.WriteByte('@')
		buffer.WriteString(selector)
	}

	return buffer.String()
}

// generateGithubTreeURLTemplate generates a github tree url.
func generateGithubTreeURLTemplate(author, repo, ref string) string {
	if len(ref) < 1 {
		ref = "master"
	}

	return fmt.Sprintf(githubTreeURLTemplate, author, repo, ref)
}

// generateDepotBlobURLTemplate generates a depot blob url.
func generateDepotBlobURLTemplate(domain, author, repo, sha string) string {
	return fmt.Sprintf(depotBlobURLTemplate, domain, author, repo, sha)
}

// generateGoGetMetadata generates metadata in the format that go-get likes it.
func generateGoGetMetadata(args generateGoGetMetadataArgs) string {
	return fmt.Sprintf(
		goGetMetadataFormat,
		args.gophrURL,
		args.depotURL,
		args.gophrURL,
		args.treeURLTemplate,
		args.blobURLTemplate,
		args.gophrURL)
}
