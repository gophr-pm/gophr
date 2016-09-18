package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
)

const (
	formKeyGoGet                = "go-get"
	formValueGoGet              = "1"
	contentTypeHTML             = "text/html"
	someFakeGitTagRef           = "refs/tags/thisisnotathinginanyrepowehopenothatitmatters"
	gitRefsInfoSubPath          = "/info/refs"
	httpLocationHeader          = "Location"
	gitUploadPackSubPath        = "/git-upload-pack"
	httpContentTypeHeader       = "Content-Type"
	packagePageURLTemplate      = "https://%s/#/packages/%s/%s"
	contentTypeGitUploadPack    = "application/x-git-upload-pack-advertisement"
	githubUploadPackURLTemplate = "https://github.com/%s/%s/git-upload-pack"
)

// TODO(skeswa): IMPORTANT! When we merge in depot, go-gets will no longer
// require ref fetches for some go-gets.

// PackageRequest is stuct that standardizes the output of all the scenarios
// through which a package may be requested. PackageRequest is essentially a
// helper struct to move data between the sub-functions of
// RespondToPackageRequest and RespondToPackageRequest itself.
type packageRequest struct {
	req             *http.Request
	parts           *packageRequestParts
	refsData        []byte
	matchedSHA      string
	matchedSHALabel string
}

// newPackageRequestArgs is the arguments struct for newPackageRequest.
type newPackageRequestArgs struct {
	req          *http.Request
	downloadRefs refsDownloader
}

// newPackageRequest parses and simplifies the information in a package version
// request in order to make serializing a response easier.
func newPackageRequest(args newPackageRequestArgs) (*packageRequest, error) {
	// Read the parts of the package request.
	parts, err := readPackageRequestParts(args.req)
	if err != nil {
		return nil, err
	}

	var (
		refs            common.Refs
		matchedSHA      string
		matchedSHALabel string
		packageRefsData []byte
	)

	// Only go out to fetch refs if they're going to get used.
	if isGoGetRequest(args.req) || isInfoRefsRequest(parts) {
		// Get and process all of the refs for this package.
		if refs, err = args.downloadRefs(
			parts.author,
			parts.repo); err != nil {
			return nil, err
		}

		// Set the default matched sha.
		matchedSHA = refs.MasterRefHash

		// Figure out what the best candidate is.
		if parts.hasSemverSelector() {
			// If there are no candidates, return in failure.
			if refs.Candidates == nil || len(refs.Candidates) < 1 {
				return nil, NewNoSuchPackageVersionError(
					parts.author,
					parts.repo,
					parts.semverSelector.String())
			}

			// Find the best candidate.
			bestCandidate := refs.Candidates.Best(parts.semverSelector)
			if bestCandidate == nil {
				return nil, NewNoSuchPackageVersionError(
					parts.author,
					parts.repo,
					parts.semverSelector.String())
			}

			// Re-serialize the refs data with said candidate.
			matchedSHA = bestCandidate.GitRefHash
			matchedSHALabel = bestCandidate.String()
			packageRefsData = refs.Reserialize(
				bestCandidate.GitRefName,
				bestCandidate.GitRefHash)
		} else if parts.hasSHASelector() {
			// Re-serialize the refs data with the sha.
			matchedSHA = parts.shaSelector
			packageRefsData = refs.Reserialize(
				someFakeGitTagRef,
				parts.shaSelector)
			// TODO(skeswa): investigate validating the ref to see if it actually
			// exists.
		} else {
			// Since there was no selector, we are fine with the fact that we didn't
			// find a match. Now, return the original refs that we downloaded from
			// github that point to master by default.
			packageRefsData = refs.Data
		}
	}

	return &packageRequest{
		req:             args.req,
		parts:           parts,
		refsData:        packageRefsData,
		matchedSHA:      matchedSHA,
		matchedSHALabel: matchedSHALabel,
	}, nil
}

// respondToPackageRequestArgs is the arguments struct for
// packageRequest#respond.
type respondToPackageRequestArgs struct {
	db                    *gocql.Session
	res                   http.ResponseWriter
	conf                  *config.Config
	creds                 *config.Credentials
	versionPackage        packageVersioner
	isPackageArchived     packageArchivalChecker
	recordPackageArchival packageArchivalRecorder
	recordPackageDownload packageDownloadRecorder
}

// respond crafts an appropriate response for a package request, serializes the
// aforesaid response and sends it back to the original client.
func (pr *packageRequest) respond(args respondToPackageRequestArgs) error {
	// Take care of the cases that deoend inf variations in the subpath.
	switch pr.parts.subpath {
	case gitUploadPackSubPath:
		// Send a 301 stipulating the repository can be found on github.
		args.res.Header().Set(
			httpLocationHeader,
			fmt.Sprintf(
				githubUploadPackURLTemplate,
				pr.parts.author,
				pr.parts.repo))
		args.res.WriteHeader(http.StatusMovedPermanently)
		return nil
	case gitRefsInfoSubPath:
		// Return the adjusted refs data when refs info is requested.
		args.res.Header().Set(httpContentTypeHeader, contentTypeGitUploadPack)
		args.res.Write(pr.refsData)
		return nil
	}

	// This means that go-get is requesting package/repository metadata.
	if isGoGetRequest(pr.req) {
		// Without blocking, count go-get surveying this package for installation as
		// a download in the database.
		go args.recordPackageDownload(packageDownloadRecorderArgs{
			db:     args.db,
			sha:    pr.matchedSHA,
			repo:   pr.parts.repo,
			author: pr.parts.author,
			// It is ok for the matched sha label to be left blank.
			version: pr.matchedSHALabel,
		})

		// Check whether this package has already been archived.
		packageArchived, err := args.isPackageArchived(packageArchivalArgs{
			db:                    args.db,
			sha:                   pr.matchedSHA,
			repo:                  pr.parts.repo,
			author:                pr.parts.author,
			recordPackageArchival: args.recordPackageArchival,
		})
		// If we cannot check whether a package has been archived, return
		// unsuccessfully.
		if err != nil {
			return err
		}

		// Only run the sub-versioning for this package if we haven't before.
		if !packageArchived {
			// Indicate in the logs that the package was archived.
			log.Printf(
				"Package %s/%s@%s has not yet been archived.\n",
				pr.parts.author,
				pr.parts.repo,
				pr.matchedSHA)

			// Perform sub-versioning.
			if err := args.versionPackage(packageVersionerArgs{
				db:                    args.db,
				sha:                   pr.matchedSHA,
				repo:                  pr.parts.repo,
				conf:                  args.conf,
				creds:                 args.creds,
				author:                pr.parts.author,
				pushToDepot:           pushToDepot,
				downloadPackage:       downloadPackage,
				constructionZonePath:  args.conf.ConstructionZonePath,
				recordPackageArchival: args.recordPackageArchival,
			}); err != nil {
				// Report the sub-versioning failure to the logs.
				log.Printf(
					"Sub-versioning failed for package %s/%s@%s: %v\n",
					pr.parts.author,
					pr.parts.repo,
					pr.matchedSHA,
					err)

				return err
			}
		}

		// Resolve the domain of the request.
		var domain string
		if len(pr.req.Host) > 0 {
			domain = pr.req.Host
		} else if len(pr.req.URL.Host) > 0 {
			domain = pr.req.URL.Host
		} else {
			// This is a last resort.
			domain = "gophr.pm"
		}

		// Compile the go-get metadata accordingly.)
		var (
			metaData = []byte(generateGoGetMetadata(generateGoGetMetadataArgs{
				gophrURL:        domain + pr.req.URL.Path,
				treeURLTemplate: generateGithubTreeURLTemplate(pr.parts.author, pr.parts.repo, pr.matchedSHA),
				blobURLTemplate: generateDepotBlobURLTemplate(domain, pr.parts.author, pr.parts.repo, pr.matchedSHA),
			}))
		)

		// Return the go-get metadata.
		args.res.Header().Set(httpContentTypeHeader, contentTypeHTML)
		args.res.Write(metaData)
		return nil
	}

	// If none of the other cases matched, then redirect to the package page.
	// TODO(skeswa): make this redirect specific to the version of the package.
	args.res.Header().Set(
		httpLocationHeader,
		fmt.Sprintf(
			packagePageURLTemplate,
			pr.req.URL.Host,
			pr.parts.author,
			pr.parts.repo))
	args.res.WriteHeader(http.StatusMovedPermanently)
	return nil
}

// isGoGetRequest returns true if the request was made by go get.
func isGoGetRequest(req *http.Request) bool {
	return req.FormValue(formKeyGoGet) == formValueGoGet
}

// isInfoRefsRequest returns true if the request parts reflect that the request
// is a git refs info request.
func isInfoRefsRequest(parts *packageRequestParts) bool {
	return parts.subpath == gitRefsInfoSubPath
}
