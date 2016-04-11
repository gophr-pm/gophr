package main

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"regexp"

	"github.com/skeswa/gophr/common"
)

const (
	packageRequestRegexIndexUser                    = 1
	packageRequestRegexIndexRepo                    = 2
	packageRequestRegexIndexSemverPrefix            = 3
	packageRequestRegexIndexSemverMajorVersion      = 4
	packageRequestRegexIndexSemverMinorVersion      = 5
	packageRequestRegexIndexSemverPatchVersion      = 6
	packageRequestRegexIndexSemverPrereleaseLabel   = 7
	packageRequestRegexIndexSemverPrereleaseVersion = 8
	packageRequestRegexIndexSemverSuffix            = 9
	packageRequestRegexIndexSubpath                 = 10
)

const (
	errorPackageRequestParsePathDoesntMatch = "The URL of the request was not a package request URL"
	errorPackageRequestParseNoSuchVersion   = "Could not find a version of \"%s/%s%s\" that matches %s"
)

const (
	formKeyGoGet                = "go-get"
	formValueGoGet              = "1"
	contentTypeHTML             = "text/html"
	masterGitRefLabel           = "master"
	gitRefsInfoSubPath          = "/info/refs"
	githubRootTemplate          = "github.com/%s/%s"
	httpLocationHeader          = "Location"
	gitUploadPackSubPath        = "/git-upload-pack"
	httpContentTypeHeader       = "Content-Type"
	packagePageURLTemplate      = "https://%s/#/packages/%s/%s"
	contentTypeGitUploadPack    = "application/x-git-upload-pack-advertisement"
	githubUploadPackURLTemplate = "https://github.com/%s/%s/git-upload-pack"
)

var (
	versionSelectorRegexStr = fmt.Sprintf(
		`([\%c\%c]?)([0-9]+)(?:\.([0-9]+|%c))?(?:\.([0-9]+|%c))?(?:\-([a-zA-Z0-9\-_]+[a-zA-Z0-9])(?:\.([0-9]+|%c))?)?([\%c\%c]?)`,
		common.SemverSelectorTildeChar,
		common.SemverSelectorCaratChar,
		common.SemverSelectorWildcardChar,
		common.SemverSelectorWildcardChar,
		common.SemverSelectorWildcardChar,
		common.SemverSelectorLessThanChar,
		common.SemverSelectorGreaterThanChar,
	)

	packageRequestRegexStr = fmt.Sprintf(
		`^\/([a-zA-Z0-9][a-zA-Z0-9\-]*[a-zA-Z0-9])\/([a-zA-Z0-9\-_]+)(?:@%s)?((?:\/[a-zA-Z0-9][-.a-zA-Z0-9]*)*)`,
		versionSelectorRegexStr,
	)

	packageRequestRegex = regexp.MustCompile(packageRequestRegexStr)

	goGetTemplate = template.Must(template.New("").Parse(
		`{{$gophrRoot := .GophrRoot}}{{$githubRoot := .GithubRoot}}{{$githubTree := .GithubTree}}<html><head><meta name="go-import" content="{{$gophrRoot}} git {{.Protocol}}://{{$gophrRoot}}"><meta name="go-source" content="{{.GophrRoot}} _ https://{{$githubRoot}}/tree/{{$githubTree}}{/dir} https://{{$githubRoot}}/blob/{{$githubTree}}{/dir}/{file}#L{line}"></head><body>go get {{.GophrPath}}</body></html>`,
	))
)

// GoGetTemplateDataSource has all the fields necessary to compile the template
// used to respond to go get metadata requests.
type GoGetTemplateDataSource struct {
	Protocol   string
	GophrRoot  string
	GophrPath  string
	GithubRoot string
	GithubTree string
}

// NewGoGetTemplateDataSource creates a new instance of GoGetTemplateDataSource.
func NewGoGetTemplateDataSource(
	user string,
	repo string,
	semverSelectorExists bool,
	semverSelector common.SemverSelector,
	subpath string,
	hasMatchedCandidate bool,
	matchedCandidate common.SemverCandidate,
) GoGetTemplateDataSource {
	var (
		buffer bytes.Buffer

		config     = getConfig()
		templateDS = GoGetTemplateDataSource{}
	)

	if config.dev {
		templateDS.Protocol = "http"
	} else {
		templateDS.Protocol = "https"
	}

	buffer.WriteString(config.domain)
	buffer.WriteByte('/')
	buffer.WriteString(user)
	buffer.WriteByte('/')
	buffer.WriteString(repo)
	if semverSelectorExists {
		buffer.WriteByte('@')
		buffer.WriteString(semverSelector.String())
	}
	templateDS.GophrRoot = buffer.String()

	buffer.WriteString(subpath)
	templateDS.GophrPath = buffer.String()

	buffer.Reset()
	buffer.WriteString("github.com")
	buffer.WriteByte('/')
	buffer.WriteString(user)
	buffer.WriteByte('/')
	buffer.WriteString(repo)
	templateDS.GithubRoot = buffer.String()

	if hasMatchedCandidate {
		templateDS.GithubTree = matchedCandidate.GitRefLabel
	} else {
		templateDS.GithubTree = masterGitRefLabel
	}

	return templateDS
}

// RespondToPackageRequest processes an incoming request, evaluates whether is a
// correctly formatted request for package-related data, and either responds
// appropriately or returns an error indicating what went wrong.
func RespondToPackageRequest(
	requestID string,
	req *http.Request,
	res http.ResponseWriter,
) error {
	matches := packageRequestRegex.FindStringSubmatch(req.URL.Path)
	if matches == nil {
		log.Printf("[%s] Failed to process request as a package request because the URL format didn't match the regular expression\n", requestID)
		return errors.New(errorPackageRequestParsePathDoesntMatch)
	}

	var (
		packageRepo          = matches[packageRequestRegexIndexRepo]
		packageCreator       = matches[packageRequestRegexIndexUser]
		packageSubpath       = matches[packageRequestRegexIndexSubpath]
		requesterIsGoGet     = req.FormValue(formKeyGoGet) == formValueGoGet
		hasMatchedCandidate  = false
		semverSelectorExists = false

		semverSelector   common.SemverSelector
		packageRefsData  []byte
		matchedCandidate common.SemverCandidate
	)

	if len(matches[packageRequestRegexIndexSemverMajorVersion]) > 0 {
		selector, err := common.NewSemverSelector(
			matches[packageRequestRegexIndexSemverPrefix],
			matches[packageRequestRegexIndexSemverMajorVersion],
			matches[packageRequestRegexIndexSemverMinorVersion],
			matches[packageRequestRegexIndexSemverPatchVersion],
			matches[packageRequestRegexIndexSemverPrereleaseLabel],
			matches[packageRequestRegexIndexSemverPrereleaseVersion],
			matches[packageRequestRegexIndexSemverSuffix],
		)

		if err != nil {
			return err
		}

		semverSelector = selector
		semverSelectorExists = true

		log.Printf("[%s] Found a version selector in the request URL\n", requestID)
	}

	// Only go out to fetch refs if they're going to get used
	if requesterIsGoGet || packageSubpath == gitRefsInfoSubPath {
		log.Printf("[%s] Fetching Github refs since this request is either from a go get or has an info path\n", requestID)

		refs, err := common.FetchRefs(fmt.Sprintf(
			githubRootTemplate,
			packageCreator,
			packageRepo,
		))

		if err != nil {
			log.Printf("[%s] Github ref fetch failed\n", requestID)
			return err
		}

		if semverSelectorExists &&
			refs.Candidates != nil &&
			len(refs.Candidates) > 0 {
			// Get the list of candidates that match the selector
			matchedCandidates := refs.Candidates.Match(semverSelector)
			log.Printf("[%s] Matched candidates to the version selector\n", requestID)
			// Only proceed if there is at least one matched candidate
			if matchedCandidates != nil && len(matchedCandidates) > 0 {
				if len(matchedCandidates) == 1 {
					matchedCandidate = matchedCandidates[0]
					hasMatchedCandidate = true
				} else {
					selectorHasLessThan :=
						semverSelector.Suffix == common.SemverSelectorSuffixLessThan
					selectorHasWildcards :=
						semverSelector.MinorVersion.Type == common.SemverSegmentTypeWildcard ||
							semverSelector.PatchVersion.Type == common.SemverSegmentTypeWildcard ||
							semverSelector.PrereleaseVersion.Type == common.SemverSegmentTypeWildcard

					var matchedCandidateReference *common.SemverCandidate
					if selectorHasWildcards || selectorHasLessThan {
						matchedCandidateReference = matchedCandidates.Highest()
					} else {
						matchedCandidateReference = matchedCandidates.Lowest()
					}

					matchedCandidate = *matchedCandidateReference
					hasMatchedCandidate = true
				}

				log.Printf("[%s] There was at least one candidate matched to the version selector\n", requestID)
			}
		}

		if hasMatchedCandidate {
			log.Printf("[%s] Tweaked the refs to focus on the matched candidate\n", requestID)
			refsData, err := refs.Reserialize(matchedCandidate)
			if err != nil {
				return err
			}
			packageRefsData = refsData
		} else {
			if !semverSelectorExists {
				// Since there was no selector, we are fine with the fact that we didn't
				// find a match. Now, return the original refs that we downloaded from
				// github that point to master by default.
				packageRefsData = refs.Data
			} else {
				log.Printf("[%s] Couldn't find any candidates to match to the version selector \"%s\"\n", requestID, semverSelector.String())

				return fmt.Errorf(
					errorPackageRequestParseNoSuchVersion,
					packageCreator,
					packageRepo,
					packageSubpath,
					semverSelector.String(),
				)
			}
		}
	}

	switch packageSubpath {
	case gitUploadPackSubPath:
		log.Printf("[%s] Responding with the Github upload pack permanent redirect\n", requestID)

		res.Header().Set(
			httpLocationHeader,
			fmt.Sprintf(githubUploadPackURLTemplate, packageCreator, packageRepo),
		)
		res.WriteHeader(http.StatusMovedPermanently)
	case gitRefsInfoSubPath:
		log.Printf("[%s] Responding with the git refs pulled from Github\n", requestID)

		res.Header().Set(httpContentTypeHeader, contentTypeGitUploadPack)
		res.Write(packageRefsData)
	default:
		if requesterIsGoGet {
			log.Printf("[%s] Responding with html formatted for go get\n", requestID)

			res.Header().Set(httpContentTypeHeader, contentTypeHTML)
			err := goGetTemplate.Execute(res, NewGoGetTemplateDataSource(
				packageCreator,
				packageRepo,
				semverSelectorExists,
				semverSelector,
				packageSubpath,
				hasMatchedCandidate,
				matchedCandidate,
			))

			if err != nil {
				log.Printf("[%s] Failed to respond with html formatted for go get\n", requestID)
				return err
			}
		} else {
			log.Printf("[%s] Responding with a permanent redirect to the gophr package webpage\n", requestID)

			res.Header().Set(
				httpLocationHeader,
				fmt.Sprintf(
					packagePageURLTemplate,
					getConfig().domain,
					packageCreator,
					packageRepo,
				),
			)
			res.WriteHeader(http.StatusMovedPermanently)
		}
	}

	return nil
}
