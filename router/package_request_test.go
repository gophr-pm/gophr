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
	"github.com/skeswa/gophr/common/semver"
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
		packageVersionerArgs{
			db:                   nil,
			sha:                  "thisshouldbeashathisshouldbeashathisshou",
			repo:                 "myrepo",
			conf:                 conf,
			creds:                creds,
			author:               "myauthor",
			constructionZonePath: "/a/b/c",
		},
		actualPVAArgs)

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
		packageVersionerArgs{
			db:                   nil,
			sha:                  "thisshouldbeashathisshouldbeashathisshou",
			repo:                 "myrepo",
			conf:                 conf,
			creds:                creds,
			author:               "myauthor",
			constructionZonePath: "/a/b/c",
		},
		actualPVAArgs)
	assert.Equal(t, 200, w.Code)
	assert.Equal(t, "text/html", w.Header().Get("Content-Type"))
	assert.Equal(
		t,
		`
<html>
<head>
<meta name="go-import" content="besthost.ever/a/s/n git https://besthost.ever/a/s/n">
<meta name="go-source" content="besthost.ever/a/s/n _ https://github.com/gophr-packages/8myauthor6myrepo/tree/thisshouldbeashathisshouldbeashathissho{/dir} https://github.com/gophr-packages/8myauthor6myrepo/blob/thisshouldbeashathisshouldbeashathissho{/dir}/{file}#L{line}">
</head>
<body>
go get besthost.ever/a/s/n
</body>
</html>
`, w.Body.String())

	// TODO one more test for if none of the other cases got matched.
}

// package main
//
// import (
// 	"net/http"
// 	"net/url"
// 	"testing"
// )
//
// var (
// 	packageRequestTestTuples = []packageRequestTestTuple{
// 		packageRequestTestTuple{
// 			path:      "/kajldlkjadshflkasjhdfl",
// 			willError: true,
// 		},
// 		packageRequestTestTuple{
// 			path:      "/codegangsta/cli@~1.x.x+",
// 			willError: true,
// 		},
// 		packageRequestTestTuple{
// 			path:      "/asdasdasdasda/thisdoesnotexistfam@1/info/refs",
// 			willError: true,
// 		},
// 		packageRequestTestTuple{
// 			path:      "/skeswa/onedark.vim@1/info/refs",
// 			willError: true,
// 		},
// 		packageRequestTestTuple{
// 			path:           "/skeswa/onedark.vim",
// 			expectedStatus: 301,
// 		},
// 		packageRequestTestTuple{
// 			path:           "/skeswa/onedark.vim/git-upload-pack",
// 			expectedStatus: 301,
// 		},
// 		packageRequestTestTuple{
// 			path:           "/skeswa/onedark.vim/info/refs",
// 			expectedStatus: 200,
// 			expectedResponse: `001e# service=git-upload-pack
// 000000e85c461cf49a467a44720fad81c7880baa66e133ef HEAD` + "\x00" + `multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done symref=HEAD:refs/heads/master agent=git/2:2.6.5+github-1394-g163a735
// 004bd8593c0c26404d92bf62c01c0bfe1eae9d950126 refs/heads/jhbabon-git-colors
// 003f5c461cf49a467a44720fad81c7880baa66e133ef refs/heads/master
// 0000`,
// 		},
// 		packageRequestTestTuple{
// 			path:               "/codegangsta/cli@1.1.x",
// 			expectedStatus:     200,
// 			expectedResponse:   `<html><head><meta name="go-import" content="gophr.dev/codegangsta/cli@1.1.x git http://gophr.dev/codegangsta/cli@1.1.x"><meta name="go-source" content="gophr.dev/codegangsta/cli@1.1.x _ https://github.com/codegangsta/cli/tree/v1.2.0{/dir} https://github.com/codegangsta/cli/blob/v1.2.0{/dir}/{file}#L{line}"></head><body>go get gophr.dev/codegangsta/cli@1.1.x</body></html>`,
// 			isGoGetMetaRequest: true,
// 		},
// 		packageRequestTestTuple{
// 			path:               "/codegangsta/cli@1+",
// 			expectedStatus:     200,
// 			expectedResponse:   `<html><head><meta name="go-import" content="gophr.dev/codegangsta/cli@1+ git http://gophr.dev/codegangsta/cli@1+"><meta name="go-source" content="gophr.dev/codegangsta/cli@1+ _ https://github.com/codegangsta/cli/tree/1.0.0{/dir} https://github.com/codegangsta/cli/blob/1.0.0{/dir}/{file}#L{line}"></head><body>go get gophr.dev/codegangsta/cli@1+</body></html>`,
// 			isGoGetMetaRequest: true,
// 		},
// 		packageRequestTestTuple{
// 			// TODO(skeswa): come up with a repo we control for this test
// 			// We can continue to expect this test will pass because its been deprecated for 2 years
// 			path:           "/keenlabs/KeenClient-Node@0.1.1/info/refs",
// 			expectedStatus: 200,
// 			expectedResponse: `001e# service=git-upload-pack
// 000000e8f521a5122c5787590c479a58c188254899dfbeb5 HEAD` + "\x00" + `multi_ack thin-pack side-band side-band-64k ofs-delta shallow no-progress include-tag multi_ack_detailed no-done oldref=HEAD:refs/heads/master agent=git/2:2.6.5+github-1394-g163a735
// 003ff521a5122c5787590c479a58c188254899dfbeb5 refs/heads/master
// 0043d541bc6b7c143b65f988653d3bb0f781daeb787d refs/heads/0.0.8-test
// 0041b67efdf74cbe422b97f87d170eee80d443b40b3d refs/heads/new-keys
// 00447a67d25cc86d0edef8350f75b5b837274a59e7e9 refs/heads/scoped-keys
// 003fdc2cd855458478e25f57293be6551b597adf4004 refs/pull/10/head
// 0040c6eda2e1a84f83256e28c5bbbf4e7bcb50f6d78a refs/pull/10/merge
// 003ff7c8ad1cdb139a504cf8c48222f81c6a6bf385df refs/pull/12/head
// 003f088b589497f1246b991cedc390043b06fab851a0 refs/pull/13/head
// 003f1d0661bae1df844c61d9585d2f72de3ca156b3ad refs/pull/16/head
// 003f1eab970f853eb069ecb5c331aa03ed82f5a3f015 refs/pull/17/head
// 0040d70118a201ddeefe6e61cb2688d5fcc1dbf405c4 refs/pull/17/merge
// 003f5bc44e96687175db34f27148d6b8698b908513c5 refs/pull/18/head
// 003fe975e99e9fd71d1e0c585ddd47fe8731c99b6d91 refs/pull/19/head
// 004010102b8b42b43ec20532eade09019b23dbccc83a refs/pull/19/merge
// 003f2ae5a3070705960872c2e952156c707586a9fa74 refs/pull/21/head
// 003f2ccfa5f319af92a655f35b11c2152be14de5b0b0 refs/pull/23/head
// 003f20ed311f890bbf5a5845bbf2f2e81dc69a218862 refs/pull/27/head
// 003f8cb3b29e990d05e5043efd1ac742bbd1451e2ed4 refs/pull/29/head
// 0040099d7a306602bfb6390522b1a71d61d350882e93 refs/pull/29/merge
// 003e7de5a0280918b3c3b0791cf418c6bbfa0e0b2de8 refs/pull/3/head
// 003fd7ccd933ade26f6471dfdeb5bb036c553518c548 refs/pull/3/merge
// 003fc0f80e70f32602cc44cb5a6db9ca0d299fe330d6 refs/pull/31/head
// 0040a92399c59d605403b49efb68df0f95cc53f9dfbc refs/pull/31/merge
// 003f4c32a5596fb18b477e829935e5c890c87c4734f3 refs/pull/33/head
// 0040a172c321abeb69735d953733ff9c8abfd2a999b5 refs/pull/33/merge
// 003f4c799069734c0cac49921007e4ae2b8832b7588f refs/pull/35/head
// 0040e59f51b45a318e4d4208250bae35194e1bc86b17 refs/pull/35/merge
// 003e9b3da5365e347fbe9dae7f2cbf650817b909fe92 refs/pull/4/head
// 003e70527cf86b738f0d0e024739811ee6aa0a786c6c refs/pull/5/head
// 003e6205da9ef5fda6778da4a0c77b5e5ea918b91861 refs/pull/6/head
// 003ed2533a7c3b9a940398e28e4d4e61c5a4c462ca51 refs/pull/7/head
// 003e4b50f64df375b9282299b3f48da183d5f42b6d38 refs/pull/8/head
// 003e48588b9835a40b8f2b1f8508cc50a0868bb0b7c2 refs/pull/9/head
// 003ebfaad0d71eb8580970b4f5b6954b2aec345b1d29 refs/tags/v0.1.0
// 00415bc44e96687175db34f27148d6b8698b908513c5 refs/tags/v0.1.0^{}
// 003ef521a5122c5787590c479a58c188254899dfbeb5 refs/tags/v0.1.1
// 003eeca61423af2a5dc1971e9d2813a1e216d37e7132 refs/tags/v0.1.2
// 0041843bd40db1776b4d65b11752c36dbf274deff2fe refs/tags/v0.1.2^{}
// 0000`,
// 		},
// 	}
// )
//
// type packageRequestTestTuple struct {
// 	path               string
// 	willError          bool
// 	expectedStatus     int
// 	expectedResponse   string
// 	isGoGetMetaRequest bool
// }
//
// type fakeResponseWriter struct {
// 	body       []byte
// 	statusCode int
// }
//
// func (w *fakeResponseWriter) Header() http.Header {
// 	return http.Header{}
// }
//
// func (w *fakeResponseWriter) Write(body []byte) (int, error) {
// 	w.body = append(w.body[:], body[:]...)
// 	return len(body), nil
// }
//
// func (w *fakeResponseWriter) WriteHeader(statusCode int) {
// 	w.statusCode = statusCode
// }
//
// func generateRequestFor(tuple packageRequestTestTuple) *http.Request {
// 	var req = http.Request{URL: &url.URL{Path: tuple.path}}
//
// 	if tuple.isGoGetMetaRequest {
// 		req.Form = url.Values{
// 			"go-get": []string{"1"},
// 		}
// 	}
//
// 	return &req
// }
//
// func TestRespondToPackageRequest(t *testing.T) {
// 	// getConfig().dev = true
// 	// getConfig().domain = "gophr.dev"
// 	// for _, tuple := range packageRequestTestTuples {
// 	// req := generateRequestFor(tuple)
// 	// res := &fakeResponseWriter{statusCode: 200}
//
// 	// TODO(skeswa): mock the database session.
// 	// err := RespondToPackageRequest(
// 	// 	nil,
// 	// 	nil,
// 	// 	nil,
// 	// 	NewRequestContext(nil),
// 	// 	req,
// 	// 	res,
// 	// )
// 	// if tuple.willError {
// 	// 	assert.NotNil(t, err, "There should be an error for "+tuple.path)
// 	// } else {
// 	// 	assert.Nil(t, err, "There should be no error")
// 	// 	assert.Equal(t, tuple.expectedStatus, res.statusCode, "The status code should match its expected value")
// 	// 	if len(tuple.expectedResponse) > 0 {
// 	// 		var bodyStr string
// 	// 		if res.body != nil && len(res.body) > 0 {
// 	// 			bodyStr = string(res.body[:len(res.body)])
// 	// 		}
// 	// 		assert.Equal(t, tuple.expectedResponse, bodyStr, "The response body should match its expected value")
// 	// 	}
// 	// }
// 	// }
// }
