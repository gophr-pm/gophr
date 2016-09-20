package main

import (
	"bytes"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"sync"
	"testing"

	"github.com/jinzhu/copier"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/depot"
	"github.com/skeswa/gophr/common/semver"
	"github.com/skeswa/gophr/common/verdeps"
	"github.com/stretchr/testify/assert"
)

var (
	baseFakeRefs, _ = common.NewRefs([]byte(reflines(
		"00000000000000000000000000000000000hash5 HEAD",
		"00000000000000000000000000000000000hash5 refs/heads/master",
		"00000000000000000000000000000000000hash3 refs/tags/v1",
		"00000000000000000000000000000000000hash4 refs/tags/v1",
		"00000000000000000000000000000000000hash5 refs/tags/v2")[:]))
)

func fakeRefs(masterRefHash string, candidates []semver.SemverCandidate) common.Refs {
	newFakeRefs := common.Refs{}
	copier.Copy(&newFakeRefs, &baseFakeRefs)
	if len(masterRefHash) > 0 {
		newFakeRefs.MasterRefHash = masterRefHash
	}
	if candidates != nil {
		newFakeRefs.Candidates = candidates
	}

	return newFakeRefs
}

func fakeHTTPRequest(host string, path string, goGet bool) *http.Request {
	form := url.Values{}
	if goGet {
		form.Add("go-get", "1")
	}

	return &http.Request{URL: &url.URL{Path: path, Host: host}, Form: form}
}

func fakeRefsDownloader(refs common.Refs, err error) refsDownloader {
	return func(author, repo string) (common.Refs, error) {
		return refs, err
	}
}

func reflines(lines ...string) string {
	var buf bytes.Buffer
	buf.WriteString("001e# service=git-upload-pack\n0000")
	for _, l := range lines {
		buf.WriteString(fmt.Sprintf("%04x%s\n", len(l)+5, l))
	}
	buf.WriteString("0000")
	return buf.String()
}

func TestNewPackageRequest(t *testing.T) {
	pr, err := newPackageRequest(newPackageRequestArgs{
		req:          fakeHTTPRequest("testalicious.af", "////", false),
		downloadRefs: fakeRefsDownloader(common.Refs{}, nil),
	})
	assert.Nil(t, pr)
	assert.NotNil(t, err)

	req := fakeHTTPRequest("testalicious.af", "/myauthor/myrepo/mysubpath", false)
	pr, err = newPackageRequest(newPackageRequestArgs{
		req:          req,
		downloadRefs: fakeRefsDownloader(fakeRefs("mymasterhash", nil), nil),
	})
	assert.NotNil(t, pr)
	assert.Nil(t, err)
	assert.Equal(t, []byte(nil), pr.refsData)
	assert.Equal(t, "", pr.matchedSHA)
	assert.Equal(t, "", pr.matchedSHALabel)

	req = fakeHTTPRequest("testalicious.af", "/myauthor/myrepo@1.x/mysubpath", true)
	pr, err = newPackageRequest(newPackageRequestArgs{
		req:          req,
		downloadRefs: fakeRefsDownloader(fakeRefs("mymasterhash", []semver.SemverCandidate{}), nil),
	})
	assert.Nil(t, pr)
	assert.NotNil(t, err)

	req = fakeHTTPRequest("testalicious.af", "/myauthor/myrepo", true)
	pr, err = newPackageRequest(newPackageRequestArgs{
		req:          req,
		downloadRefs: fakeRefsDownloader(fakeRefs("mymasterhash", nil), nil),
	})
	assert.NotNil(t, pr)
	assert.Nil(t, err)
	assert.Equal(t, baseFakeRefs.Data, pr.refsData)
	assert.Equal(t, "mymasterhash", pr.matchedSHA)
	assert.Equal(t, "", pr.matchedSHALabel)

	req = fakeHTTPRequest("testalicious.af", "/myauthor/myrepo@1.x/mysubpath", true)
	pr, err = newPackageRequest(newPackageRequestArgs{
		req: req,
		downloadRefs: fakeRefsDownloader(fakeRefs(
			"mymasterhash",
			[]semver.SemverCandidate{
				semver.SemverCandidate{
					GitRefHash:   "GitRefHash1GitRefHash1GitRefHash1GitRefH",
					GitRefName:   "GitRefName1",
					GitRefLabel:  "GitRefLabel1",
					MajorVersion: 1,
					MinorVersion: 2,
					PatchVersion: 0,
				},
				semver.SemverCandidate{
					GitRefHash:   "GitRefHash2GitRefHash2GitRefHash2GitRefH",
					GitRefName:   "GitRefName2",
					GitRefLabel:  "GitRefLabel2",
					MajorVersion: 2,
					MinorVersion: 4,
					PatchVersion: 1,
				},
			}), nil),
	})
	assert.NotNil(t, pr)
	assert.Nil(t, err)
	assert.Equal(
		t,
		"001e# service=git-upload-pack\n00000032GitRefHash1GitRefHash1GitRefHash1GitRefH HEAD\n003fGitRefHash1GitRefHash1GitRefHash1GitRefH refs/heads/master\n003a00000000000000000000000000000000000hash3 refs/tags/v1\n003a00000000000000000000000000000000000hash4 refs/tags/v1\n003a00000000000000000000000000000000000hash5 refs/tags/v2\n0000",
		string(pr.refsData[:]))
	assert.Equal(t, "GitRefHash1GitRefHash1GitRefHash1GitRefH", pr.matchedSHA)
	assert.Equal(t, "1.2.0", pr.matchedSHALabel)

	req = fakeHTTPRequest("testalicious.af", "/myauthor/myrepo@1234567890123456789012345678901234567890", true)
	pr, err = newPackageRequest(newPackageRequestArgs{
		req:          req,
		downloadRefs: fakeRefsDownloader(fakeRefs("somemasterhash", nil), nil),
	})
	assert.NotNil(t, pr)
	assert.Nil(t, err)
	assert.Equal(
		t,
		"001e# service=git-upload-pack\n000000321234567890123456789012345678901234567890 HEAD\n003f1234567890123456789012345678901234567890 refs/heads/master\n003a00000000000000000000000000000000000hash3 refs/tags/v1\n003a00000000000000000000000000000000000hash4 refs/tags/v1\n003a00000000000000000000000000000000000hash5 refs/tags/v2\n0000",
		string(pr.refsData[:]))
	assert.Equal(t, "1234567890123456789012345678901234567890", pr.matchedSHA)
	assert.Equal(t, "", pr.matchedSHALabel)
}

func TestRespondToPackageRequest(t *testing.T) {
	w := httptest.NewRecorder()
	err := (&packageRequest{
		parts: &packageRequestParts{
			repo:    "xyz",
			author:  "abc",
			subpath: "/git-upload-pack",
		},
	}).respond(respondToPackageRequestArgs{
		res: w,
	})
	assert.Nil(t, err)
	assert.Equal(t, 301, w.Code)
	assert.Equal(t, "https://github.com/abc/xyz/git-upload-pack", w.Header().Get("Location"))

	w = httptest.NewRecorder()
	refsData := []byte{1, 2, 3, 4}
	err = (&packageRequest{
		parts:    &packageRequestParts{subpath: "/info/refs"},
		refsData: refsData,
	}).respond(respondToPackageRequestArgs{
		res: w,
	})
	assert.Nil(t, err)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "application/x-git-upload-pack-advertisement", w.Header().Get("Content-Type"))
	assert.Equal(t, refsData, w.Body.Bytes())

	var (
		actualPAAArgs packageArchivalArgs
		actualRPDArgs packageDownloadRecorderArgs
	)

	wg1 := sync.WaitGroup{}
	wg1.Add(1)
	w = httptest.NewRecorder()
	err = (&packageRequest{
		req: fakeHTTPRequest("besthost.ever", "/a/s/n", true),
		parts: &packageRequestParts{
			repo:   "myrepo",
			author: "myauthor",
		},
		refsData:        baseFakeRefs.Data,
		matchedSHA:      "thisshouldbeashathisshouldbeashathisshou",
		matchedSHALabel: "someshalabel",
	}).respond(respondToPackageRequestArgs{
		res: w,
		recordPackageDownload: func(args packageDownloadRecorderArgs) {
			actualRPDArgs = args
			wg1.Done()
		},
		isPackageArchived: func(args packageArchivalArgs) (bool, error) {
			actualPAAArgs = args
			return false, errors.New("this is an error")
		},
	})
	wg1.Wait()
	assert.NotNil(t, err)
	assert.Equal(
		t,
		packageArchivalArgs{
			db:     nil,
			sha:    "thisshouldbeashathisshouldbeashathisshou",
			repo:   "myrepo",
			author: "myauthor",
		},
		actualPAAArgs)
	assert.Equal(
		t,
		packageDownloadRecorderArgs{
			db:      nil,
			sha:     "thisshouldbeashathisshouldbeashathisshou",
			repo:    "myrepo",
			author:  "myauthor",
			version: "someshalabel",
		},
		actualRPDArgs)

	var (
		conf = &config.Config{
			Port:                 9999,
			ConstructionZonePath: "/a/b/c",
		}
		creds = &config.Credentials{
			GithubPush: config.UserPass{User: "a", Pass: "b"},
		}
		actualPVAArgs packageVersionerArgs
	)

	wg2 := sync.WaitGroup{}
	wg2.Add(1)
	w = httptest.NewRecorder()
	err = (&packageRequest{
		req: fakeHTTPRequest("besthost.ever", "/a/s/n", true),
		parts: &packageRequestParts{
			repo:   "myrepo",
			author: "myauthor",
		},
		refsData:        baseFakeRefs.Data,
		matchedSHA:      "thisshouldbeashathisshouldbeashathisshou",
		matchedSHALabel: "someshalabel",
	}).respond(respondToPackageRequestArgs{
		res:   w,
		conf:  conf,
		creds: creds,
		recordPackageDownload: func(args packageDownloadRecorderArgs) {
			actualRPDArgs = args
			wg2.Done()
		},
		isPackageArchived: func(args packageArchivalArgs) (bool, error) {
			actualPAAArgs = args
			return false, nil
		},
		versionPackage: func(args packageVersionerArgs) error {
			actualPVAArgs = args
			return errors.New("this is an error")
		},
	})
	wg2.Wait()
	assert.NotNil(t, err)
	assert.Equal(
		t,
		packageArchivalArgs{
			db:     nil,
			sha:    "thisshouldbeashathisshouldbeashathisshou",
			repo:   "myrepo",
			author: "myauthor",
		},
		actualPAAArgs)
	assert.Equal(
		t,
		packageDownloadRecorderArgs{
			db:      nil,
			sha:     "thisshouldbeashathisshouldbeashathisshou",
			repo:    "myrepo",
			author:  "myauthor",
			version: "someshalabel",
		},
		actualRPDArgs)
	assert.Equal(
		t,
		fmt.Sprintf("%v", packageVersionerArgs{
			db:                     nil,
			sha:                    "thisshouldbeashathisshouldbeashathisshou",
			repo:                   "myrepo",
			conf:                   conf,
			creds:                  creds,
			author:                 "myauthor",
			pushToDepot:            pushToDepot,
			versionDeps:            verdeps.VersionDeps,
			downloadPackage:        downloadPackage,
			createDepotRepo:        depot.CreateNewRepo,
			destroyDepotRepo:       depot.DestroyRepo,
			isPackageArchived:      isPackageArchived,
			constructionZonePath:   "/a/b/c",
			attemptWorkDirDeletion: deleteFolder,
		}),
		fmt.Sprintf("%v", actualPVAArgs))

	wg3 := sync.WaitGroup{}
	wg3.Add(1)
	w = httptest.NewRecorder()
	err = (&packageRequest{
		req: fakeHTTPRequest("besthost.ever", "/a/s/n", true),
		parts: &packageRequestParts{
			repo:   "myrepo",
			author: "myauthor",
		},
		refsData:        baseFakeRefs.Data,
		matchedSHA:      "thisshouldbeashathisshouldbeashathisshou",
		matchedSHALabel: "someshalabel",
	}).respond(respondToPackageRequestArgs{
		res:   w,
		conf:  conf,
		creds: creds,
		recordPackageDownload: func(args packageDownloadRecorderArgs) {
			actualRPDArgs = args
			wg3.Done()
		},
		isPackageArchived: func(args packageArchivalArgs) (bool, error) {
			actualPAAArgs = args
			return false, nil
		},
		versionPackage: func(args packageVersionerArgs) error {
			actualPVAArgs = args
			return nil
		},
	})
	wg3.Wait()
	assert.Nil(t, err)
	assert.Equal(
		t,
		packageArchivalArgs{
			db:     nil,
			sha:    "thisshouldbeashathisshouldbeashathisshou",
			repo:   "myrepo",
			author: "myauthor",
		},
		actualPAAArgs)
	assert.Equal(
		t,
		packageDownloadRecorderArgs{
			db:      nil,
			sha:     "thisshouldbeashathisshouldbeashathisshou",
			repo:    "myrepo",
			author:  "myauthor",
			version: "someshalabel",
		},
		actualRPDArgs)
	assert.Equal(
		t,
		fmt.Sprintf("%v", packageVersionerArgs{
			db:                     nil,
			sha:                    "thisshouldbeashathisshouldbeashathisshou",
			repo:                   "myrepo",
			conf:                   conf,
			creds:                  creds,
			author:                 "myauthor",
			pushToDepot:            pushToDepot,
			versionDeps:            verdeps.VersionDeps,
			downloadPackage:        downloadPackage,
			createDepotRepo:        depot.CreateNewRepo,
			destroyDepotRepo:       depot.DestroyRepo,
			isPackageArchived:      isPackageArchived,
			constructionZonePath:   "/a/b/c",
			attemptWorkDirDeletion: deleteFolder,
		}),
		fmt.Sprintf("%v", actualPVAArgs))
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
	assert.Equal(
		t,
		`
<html>
<head>
<meta name="go-import" content="besthost.ever/a/s/n git https://besthost.ever/a/s/n">
<meta name="go-source" content="besthost.ever/a/s/n _ https://github.com/myauthor/myrepo/tree/thisshouldbeashathisshouldbeashathisshou{/dir} https://besthost.ever/blob/myauthor/myrepo/thisshouldbeashathisshouldbeashathisshou{/dir}/{file}#L{line}">
</head>
<body>
go get besthost.ever/a/s/n
</body>
</html>
`, w.Body.String())

	w = httptest.NewRecorder()
	err = (&packageRequest{
		req: fakeHTTPRequest("besthost.ever", "/a/s/n", false),
		parts: &packageRequestParts{
			repo:   "re",
			author: "auth",
		},
	}).respond(respondToPackageRequestArgs{
		res: w,
	})
	assert.Nil(t, err)
	assert.Equal(t, 301, w.Code)
	assert.Equal(t, "https://besthost.ever/#/packages/auth/re", w.Header().Get("Location"))
}
