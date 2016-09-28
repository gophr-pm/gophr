package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/config"
	"github.com/skeswa/gophr/common/depot"
	"github.com/skeswa/gophr/common/github"
	"github.com/skeswa/gophr/common/io"
	"github.com/skeswa/gophr/common/models"
	"github.com/skeswa/gophr/common/verdeps"
)

const (
	formKeyGoGet           = "go-get"
	formValueGoGet         = "1"
	contentTypeHTML        = "text/html"
	httpLocationHeader     = "Location"
	gitInfoRefsSubPath     = "/info/refs"
	depotRepoURLTemplate   = "https://%s/depot/%s.git%s"
	gitUploadPackSubPath   = "/git-upload-pack"
	httpContentTypeHeader  = "Content-Type"
	basePackageURLTemplate = "https://%s%s"
	packagePageURLTemplate = "https://%s/#/packages/%s/%s"
)

// PackageRequest is stuct that standardizes the output of all the scenarios
// through which a package may be requested. PackageRequest is essentially a
// helper struct to move data between the sub-functions of
// RespondToPackageRequest and RespondToPackageRequest itself.
type packageRequest struct {
	req             *http.Request
	parts           *packageRequestParts
	matchedSHA      string
	matchedSHALabel string
}

// newPackageRequestArgs is the arguments struct for newPackageRequest.
type newPackageRequestArgs struct {
	req          *http.Request
	downloadRefs refsDownloader
	fetchFullSHA fullSHAFetcher
	doHTTPHead   github.HTTPHeadReq
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
	)

	if isGoGetRequest(args.req) || isGitRequest(parts) {
		// Check if we have a SHA selector.
		if parts.hasSHASelector() {
			// If we have a short SHA selector convert it to a full SHA.
			if parts.hasShortSHASelector {
				matchedSHA, err = args.fetchFullSHA(
					github.FetchFullSHAArgs{
						Author:     parts.author,
						Repo:       parts.repo,
						ShortSHA:   parts.shaSelector,
						DoHTTPHead: args.doHTTPHead,
					},
				)
				if err != nil {
					return nil, err
				}
			}

			// If we have a full SHA selector set the matchedSHA.
			if parts.hasFullSHASelector {
				matchedSHA = parts.shaSelector
			}
		}

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

			return &packageRequest{
				req:             args.req,
				parts:           parts,
				matchedSHA:      matchedSHA,
				matchedSHALabel: matchedSHALabel,
			}, nil
		}

		// Set the default matched sha in case there is no semver selector.
		matchedSHA = refs.MasterRefHash
	}

	return &packageRequest{
		req:             args.req,
		parts:           parts,
		matchedSHA:      matchedSHA,
		matchedSHALabel: matchedSHALabel,
	}, nil
}

// respondToPackageRequestArgs is the arguments struct for
// packageRequest#respond.
type respondToPackageRequestArgs struct {
	io                    io.IO
	db                    *gocql.Session
	res                   http.ResponseWriter
	conf                  *config.Config
	creds                 *config.Credentials
	ghSvc                 github.RequestService
	versionPackage        packageVersioner
	isPackageArchived     packageArchivalChecker
	recordPackageArchival packageArchivalRecorder
	recordPackageDownload packageDownloadRecorder
}

// respond crafts an appropriate response for a package request, serializes the
// aforesaid response and sends it back to the original client.
func (pr *packageRequest) respond(args respondToPackageRequestArgs) error {
	// Git requests must be redirected to depot.
	if isGitRequest(pr.parts) {
		// Join the repo URL with the subpath of this request.
		redirectLocation := fmt.Sprintf(
			depotRepoURLTemplate,
			getRequestDomain(pr.req),
			depot.BuildHashedRepoName(
				pr.parts.author,
				pr.parts.repo,
				pr.matchedSHA),
			pr.parts.subpath)

		// Issue a permanent redirect.
		http.Redirect(
			args.res,
			pr.req,
			redirectLocation,
			http.StatusMovedPermanently)

		// Without blocking, count a packfile request as a package download.
		if pr.parts.subpath == gitUploadPackSubPath {
			go args.recordPackageDownload(packageDownloadRecorderArgs{
				db:     args.db,
				sha:    pr.matchedSHA,
				repo:   pr.parts.repo,
				author: pr.parts.author,
				// It is ok for the matched sha label to be left blank.
				version: pr.matchedSHALabel,
			})
		}

		return nil
	}

	// This means that go-get is requesting package/repository metadata.
	if isGoGetRequest(pr.req) {
		// Check whether this package has already been archived.
		packageArchived, err := args.isPackageArchived(packageArchivalCheckerArgs{
			db:                    args.db,
			sha:                   pr.matchedSHA,
			repo:                  pr.parts.repo,
			author:                pr.parts.author,
			packageExistsInDepot:  packageExistsInDepot,
			recordPackageArchival: args.recordPackageArchival,
			isPackageArchivedInDB: models.IsPackageArchived,
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
				io:                     args.io,
				db:                     args.db,
				sha:                    pr.matchedSHA,
				repo:                   pr.parts.repo,
				conf:                   args.conf,
				creds:                  args.creds,
				ghSvc:                  args.ghSvc,
				author:                 pr.parts.author,
				pushToDepot:            pushToDepot,
				versionDeps:            verdeps.VersionDeps,
				downloadPackage:        downloadPackage,
				createDepotRepo:        createRepoInDepot,
				destroyDepotRepo:       deleteRepoInDepot,
				isPackageArchived:      isPackageArchived,
				constructionZonePath:   args.conf.ConstructionZonePath,
				recordPackageArchival:  args.recordPackageArchival,
				attemptWorkDirDeletion: deleteFolder,
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

		// Compile the go-get metadata accordingly.
		var (
			domain         = getRequestDomain(pr.req)
			basePackageURL = domain + pr.parts.getBasePackagePath()
			metaData       = []byte(generateGoGetMetadata(generateGoGetMetadataArgs{
				gophrURL: basePackageURL,
				treeURLTemplate: generateGithubTreeURLTemplate(
					pr.parts.author,
					pr.parts.repo,
					pr.matchedSHA),
				blobURLTemplate: generateDepotBlobURLTemplate(
					domain,
					pr.parts.author,
					pr.parts.repo,
					pr.matchedSHA),
			}))
		)

		// Return the go-get metadata.
		args.res.Header().Set(httpContentTypeHeader, contentTypeHTML)
		args.res.Write(metaData)
		return nil
	}

	// If none of the other cases matched, then redirect to the package page.
	// TODO(skeswa): make this redirect specific to the version of the package.
	http.Redirect(
		args.res,
		pr.req,
		fmt.Sprintf(
			packagePageURLTemplate,
			getRequestDomain(pr.req),
			pr.parts.author,
			pr.parts.repo),
		http.StatusMovedPermanently)
	return nil
}

// isGoGetRequest returns true if the request was made by go get.
func isGoGetRequest(req *http.Request) bool {
	return req.FormValue(formKeyGoGet) == formValueGoGet
}

// isGitRequest returns true if the request was made by git (in a clone setting).
func isGitRequest(parts *packageRequestParts) bool {
	return parts.subpath == gitInfoRefsSubPath ||
		parts.subpath == gitUploadPackSubPath
}

// getRequestDomain isolates and returns the domain of the specified request.
func getRequestDomain(req *http.Request) string {
	// If there is no request, don't make a fuss: just return empty.
	if req == nil {
		return ""
	}

	// Resolve the domain of the request.
	var domain string
	if len(req.Host) > 0 {
		domain = req.Host
	} else if len(req.URL.Host) > 0 {
		domain = req.URL.Host
	} else {
		// This is a last resort.
		domain = "gophr.pm"
	}

	return domain
}
