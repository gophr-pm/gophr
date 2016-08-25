package common

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/skeswa/gophr/common/semver"
)

const (
	errorRefsFetchNoSuchRepo       = "Could not find a Github repository at %s"
	errorRefsFetchGithubError      = "Github responded with an error: %v"
	errorRefsFetchGithubParseError = "Cannot read refs from Github: %v"
	errorRefsFetchNetworkFailure   = "Could not reach Github at the moment; Please try again later"
	errorRefsParseSizeFormat       = "Could not parse refs line size: %s"
	errorRefsParseIncompleteRefs   = "Incomplete refs data received from GitHub"
)

const (
	versionRefRegexIndexLabel             = 1
	versionRefRegexIndexMajorVersion      = 2
	versionRefRegexIndexMinorVersion      = 3
	versionRefRegexIndexPatchVersion      = 4
	versionRefRegexIndexPrereleaseLabel   = 5
	versionRefRegexIndexPrereleaseVersion = 6
)

const (
	refsHead                                  = "HEAD"
	refsLineCap                               = "\n\x00"
	refsSpaceChar                             = ' '
	refsHeadPrefix                            = "refs/heads/"
	refsLineFormat                            = "%04x%s"
	refsHeadMaster                            = "refs/heads/master"
	githubRootTemplate                        = "github.com/%s/%s"
	refsMasterLineFormat                      = "%s refs/heads/master\n"
	refsSymRefAssignment                      = "symref="
	refsOldRefAssignment                      = "oldref="
	refsFetchURLTemplate                      = "https://%s.git/info/refs?service=git-upload-pack"
	refsAugmentedHeadLineFormat               = "%s HEAD\n"
	refsAugmentedSymrefHeadLineFormat         = "%s HEAD\x00symref=HEAD:%s\n"
	refsAugmentedHeadLineWithCapsFormat       = "%s HEAD\x00%s\n"
	refsAugmentedSymrefHeadLineWithCapsFormat = "%s HEAD\x00symref=HEAD:%s %s\n"
)

var (
	httpClient      = &http.Client{Timeout: 10 * time.Second}
	versionRefRegex = regexp.MustCompile(`^refs\/(?:tags|heads)\/(v?([0-9]+)(?:\.([0-9]+))?(?:\.([0-9]+))?(?:\-([a-zA-Z0-9\-_]+))?(?:\.([0-9]+))?)(?:\^\{\})?`)
)

// Refs collects information about git references for one specific repository.
type Refs struct {
	Data                 []byte
	DataStr              string
	DataLen              int
	DataStrLen           int
	Candidates           semver.SemverCandidateList
	MasterRefHash        string
	IndexHeadLineEnd     int
	IndexHeadLineStart   int
	IndexMasterLineEnd   int
	IndexMasterLineStart int
}

// NewRefs creates a new Refs instance from raw refs data fetched from Github
// (or elsewhere).
func NewRefs(data []byte) (Refs, error) {
	var (
		dataStr    = string(data)
		dataLen    = len(data)
		dataStrLen = len(dataStr)

		masterRefHash                                 string
		indexHashStart, indexHashEnd                  int
		indexNameStart, indexNameEnd                  int
		indexHeadLineStart, indexHeadLineEnd          int
		indexMasterLineStart, indexMasterLineEnd      int
		versionCandidates, sanitizedVersionCandidates []semver.SemverCandidate
	)

	for i, j := 0, 0; i < dataLen; i = j {
		// Calculate the size by reading and parsing the size string
		size, err := strconv.ParseInt(dataStr[i:i+4], 16, 32)

		// If we can't read the hex, conclude that it was invalid
		if err != nil {
			return Refs{}, fmt.Errorf(errorRefsParseSizeFormat, string(data[i:i+4]))
		}

		// If we found that the size was zero, advance it by 4 since 4 is the
		// acceptable minimum
		if size == 0 {
			size = 4
		}

		// Advance the second cursor so the next token is bounded by the two cursors
		j = i + int(size)

		// If the second cursor exceeds the string boundary, then conclude that the
		// refs data is incomplete
		if j > len(dataStr) {
			return Refs{}, errors.New(errorRefsParseIncompleteRefs)
		}

		// TODO(skeswa): figure out why this line is here
		if dataStr[0] == '#' {
			continue
		}

		// Use the cursors to get the indices of the hash
		indexHashStart = i + 4
		indexHashEnd = strings.IndexByte(
			dataStr[indexHashStart:j],
			refsSpaceChar,
		)

		// Check for invalid hash end
		if indexHashEnd < 0 || indexHashEnd != 40 {
			continue
		}

		// TODO(skeswa): figure out why this line is here
		indexHashEnd += indexHashStart

		// Use the cursors to get the indices of the name
		indexNameStart = indexHashEnd + 1
		indexNameEnd = strings.IndexAny(
			dataStr[indexNameStart:j],
			refsLineCap,
		)

		// Check for invalid name end
		if indexNameEnd < 0 {
			indexNameEnd = j
		} else {
			indexNameEnd += indexNameStart
		}

		// Get the name and hash respectively as strings
		hash := dataStr[indexHashStart:indexHashEnd]
		name := dataStr[indexNameStart:indexNameEnd]

		// Process the name and hash according to whether the name is relevant
		if name == refsHead {
			indexHeadLineStart = i
			indexHeadLineEnd = j
		} else if name == refsHeadMaster {
			indexMasterLineStart = i
			indexMasterLineEnd = j
			masterRefHash = hash
		} else if captureGroups := versionRefRegex.FindStringSubmatch(name); captureGroups != nil {
			var (
				gitRefLabel       = captureGroups[versionRefRegexIndexLabel]
				majorVersion      = captureGroups[versionRefRegexIndexMajorVersion]
				minorVersion      = captureGroups[versionRefRegexIndexMinorVersion]
				patchVersion      = captureGroups[versionRefRegexIndexPatchVersion]
				prereleaseLabel   = captureGroups[versionRefRegexIndexPrereleaseLabel]
				prereleaseVersion = captureGroups[versionRefRegexIndexPrereleaseVersion]
			)

			// Annotated tag is peeled off and overrides the same version just parsed
			if strings.HasSuffix(name, "^{}") {
				name = name[:len(name)-3]
			}

			versionCandidate, err := semver.NewSemverCandidate(
				hash,
				name,
				gitRefLabel,
				majorVersion,
				minorVersion,
				patchVersion,
				prereleaseLabel,
				prereleaseVersion)
			if err == nil {
				versionCandidates = append(versionCandidates, versionCandidate)
			}
		}
	}

	if versionCandidates != nil && len(versionCandidates) > 0 {
		// First attach the sortable type to the slice of candidates.
		versionCandidatesList := semver.SemverCandidateList(versionCandidates)
		// Sort the list of candidates.
		sort.Sort(versionCandidatesList)
		// Remove duplicates by adding them to a new slice altogether.
		var lastInsertedCandidate semver.SemverCandidate
		for i, versionCandidate := range versionCandidatesList {
			if i == 0 || versionCandidate.CompareTo(lastInsertedCandidate) != 0 {
				sanitizedVersionCandidates = append(sanitizedVersionCandidates, versionCandidate)
				lastInsertedCandidate = versionCandidate
			}
		}
	}

	return Refs{
		Data:                 data,
		DataStr:              dataStr,
		DataLen:              dataLen,
		DataStrLen:           dataStrLen,
		Candidates:           sanitizedVersionCandidates,
		MasterRefHash:        masterRefHash,
		IndexHeadLineEnd:     indexHeadLineEnd,
		IndexMasterLineEnd:   indexMasterLineEnd,
		IndexHeadLineStart:   indexHeadLineStart,
		IndexMasterLineStart: indexMasterLineStart,
	}, nil
}

// FetchRefs downloads and processes refs data from Github and ultimately
// contructs a Refs instance with it.
func FetchRefs(author, repo string) (Refs, error) {
	githubRoot := fmt.Sprintf(
		githubRootTemplate,
		author,
		repo,
	)

	res, err := httpClient.Get(fmt.Sprintf(refsFetchURLTemplate, githubRoot))
	if err != nil {
		return Refs{}, errors.New(errorRefsFetchNetworkFailure)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return Refs{}, fmt.Errorf(errorRefsFetchNoSuchRepo, githubRoot)
	} else if res.StatusCode >= 500 {
		// FYI no reliable way to get test coverage here; this never happens
		return Refs{}, fmt.Errorf(errorRefsFetchGithubError, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// FYI no reliable way to get test coverage here; this never happens
		return Refs{}, fmt.Errorf(errorRefsFetchGithubParseError, err)
	}

	return NewRefs(data)
}

// Reserialize changes the refs data to incorporate the selected version as
// the HEAD instead of the default HEAD.
//
// This code was written by Gustavo Niemeyer, Nathan Youngman and
// Geert-Johan Riemer.
func (refsData Refs) Reserialize(versionRefName, versionRefHash string) []byte {
	var (
		buf bytes.Buffer

		data                 = refsData.Data
		dataLen              = refsData.DataLen
		dataStr              = refsData.DataStr
		indexHeadLineEnd     = refsData.IndexHeadLineEnd
		indexHeadLineStart   = refsData.IndexHeadLineStart
		indexMasterLineEnd   = refsData.IndexMasterLineEnd
		indexMasterLineStart = refsData.IndexMasterLineStart
	)

	// Size the buffer to be a little bigger than
	buf.Grow(dataLen + 256)

	// Copy the header as-is.
	buf.Write(data[:indexHeadLineStart])

	// Extract the original capabilities.
	caps := ""
	indexNullByte := strings.Index(
		dataStr[indexHeadLineStart:indexHeadLineEnd],
		"\x00",
	)

	// IF we found a zero byte, replace the symref with an oldref
	if indexNullByte > 0 {
		caps = strings.Replace(
			dataStr[indexHeadLineStart+indexNullByte+1:indexHeadLineEnd-1],
			refsSymRefAssignment,
			refsOldRefAssignment,
			-1,
		)
	}

	// Insert the HEAD reference line with the right hash and a proper symref
	// capability.
	var line string
	if strings.HasPrefix(versionRefName, refsHeadPrefix) {
		if caps == "" {
			line = fmt.Sprintf(
				refsAugmentedSymrefHeadLineFormat,
				versionRefHash,
				versionRefName,
			)
		} else {
			line = fmt.Sprintf(
				refsAugmentedSymrefHeadLineWithCapsFormat,
				versionRefHash,
				versionRefName,
				caps,
			)
		}
	} else {
		if caps == "" {
			line = fmt.Sprintf(refsAugmentedHeadLineFormat, versionRefHash)
		} else {
			line = fmt.Sprintf(refsAugmentedHeadLineWithCapsFormat, versionRefHash, caps)
		}
	}
	fmt.Fprintf(&buf, "%04x%s", 4+len(line), line)

	// Insert the master reference line.
	line = fmt.Sprintf(refsMasterLineFormat, versionRefHash)
	fmt.Fprintf(&buf, refsLineFormat, 4+len(line), line)

	// Append the rest, dropping the original master line if necessary.
	if indexMasterLineStart > 0 {
		buf.Write(data[indexHeadLineEnd:indexMasterLineStart])
		buf.Write(data[indexMasterLineEnd:])
	} else {
		buf.Write(data[indexHeadLineEnd:])
	}

	return buf.Bytes()
}

// CheckIfRefExists downloads and processes refs data from Github and checks
// whether a given ref exists in the remote refs list.
func CheckIfRefExists(author, repo string, ref string) (bool, error) {
	ref = BuildGitHubBranch(ref)
	repo = BuildNewGitHubRepoName(author, repo)
	author = GitHubGophrPackageOrgName
	githubRoot := fmt.Sprintf(
		githubRootTemplate,
		author,
		repo,
	)

	res, err := httpClient.Get(fmt.Sprintf(refsFetchURLTemplate, githubRoot))
	if err != nil {
		return false, errors.New(errorRefsFetchNetworkFailure)
	}

	defer res.Body.Close()

	if res.StatusCode >= 400 && res.StatusCode < 500 {
		return false, fmt.Errorf(errorRefsFetchNoSuchRepo, githubRoot)
	} else if res.StatusCode >= 500 {
		// FYI no reliable way to get test coverage here; this never happens
		return false, fmt.Errorf(errorRefsFetchGithubError, res.Status)
	}

	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		// FYI no reliable way to get test coverage here; this never happens
		return false, fmt.Errorf(errorRefsFetchGithubParseError, err)
	}

	refsString := string(data)
	refExists := strings.Contains(refsString, ref)

	return refExists, nil
}
