package semver

import (
	"bytes"
	"errors"
	"strconv"
	"strings"
)

// SemverCandidate is a semver version that has been confirmed to exist for a
// given package. It carries versioning metadata, but it also has git ref info
// so that the commit of the version can be isolated.
type SemverCandidate struct {
	GitRefHash              string
	GitRefName              string
	GitRefLabel             string
	MajorVersion            int
	MinorVersion            int
	PatchVersion            int
	PrereleaseLabel         string
	PrereleaseVersion       int
	PrereleaseVersionExists bool
}

// NewSemverCandidate creates a new instance of a SemverCandidate from a variety
// of data points about a specific version-related git ref.
func NewSemverCandidate(
	gitRefHash string,
	gitRefName string,
	gitRefLabel string,
	majorVersion string,
	minorVersion string,
	patchVersion string,
	prereleaseLabel string,
	prereleaseVersion string,
) (SemverCandidate, error) {
	var (
		err                     error
		majorVersionNumber      int
		minorVersionNumber      int
		patchVersionNumber      int
		prereleaseVersionNumber int
	)

	if len(gitRefHash) == 0 {
		return SemverCandidate{}, errors.New("Git ref hash is required")
	} else if len(gitRefName) == 0 {
		return SemverCandidate{}, errors.New("Git ref name is required")
	}

	if len(majorVersion) > 0 {
		majorVersionNumber, err = strconv.Atoi(majorVersion)
		if err != nil {
			return SemverCandidate{}, err
		}
	} else {
		return SemverCandidate{}, errors.New("Major version is required")
	}

	if len(minorVersion) > 0 {
		minorVersionNumber, err = strconv.Atoi(minorVersion)
		if err != nil {
			return SemverCandidate{}, err
		}
	} else {
		minorVersionNumber = 0
	}

	if len(patchVersion) > 0 {
		patchVersionNumber, err = strconv.Atoi(patchVersion)
		if err != nil {
			return SemverCandidate{}, err
		}
	} else {
		patchVersionNumber = 0
	}

	if len(prereleaseVersion) > 0 {
		prereleaseVersionNumber, err = strconv.Atoi(prereleaseVersion)
		if err != nil {
			return SemverCandidate{}, err
		}
	} else {
		prereleaseVersionNumber = 0
	}

	return SemverCandidate{
		GitRefHash:              gitRefHash,
		GitRefName:              gitRefName,
		GitRefLabel:             gitRefLabel,
		MajorVersion:            majorVersionNumber,
		MinorVersion:            minorVersionNumber,
		PatchVersion:            patchVersionNumber,
		PrereleaseLabel:         prereleaseLabel,
		PrereleaseVersion:       prereleaseVersionNumber,
		PrereleaseVersionExists: (len(prereleaseLabel) > 0),
	}, nil
}

// CompareTo compares the current candidate to another candidate and returns a
// number indicating the relationship between the two. -1 means this candidate
// is lower than the other. 1 implies the opposite. 0 means that the candidates
// are functionally equivalent.
func (candidate SemverCandidate) CompareTo(other SemverCandidate) int {
	if candidate.MajorVersion > other.MajorVersion {
		return 1
	} else if candidate.MajorVersion < other.MajorVersion {
		return -1
	} else if candidate.MinorVersion > other.MinorVersion {
		return 1
	} else if candidate.MinorVersion < other.MinorVersion {
		return -1
	} else if candidate.PatchVersion > other.PatchVersion {
		return 1
	} else if candidate.PatchVersion < other.PatchVersion {
		return -1
	} else if len(candidate.PrereleaseLabel) > 0 && len(other.PrereleaseLabel) == 0 {
		// Prerelease immediately means that the version is lesser
		return -1
	} else if len(candidate.PrereleaseLabel) == 0 && len(other.PrereleaseLabel) > 0 {
		// Prerelease immediately means that the version is lesser
		return 1
	} else if candidate.PrereleaseLabel != other.PrereleaseLabel {
		// Prerelease labels don't match, so return a comparison
		return strings.Compare(candidate.PrereleaseLabel, other.PrereleaseLabel)
	} else if candidate.PrereleaseVersion > other.PrereleaseVersion {
		// Prerelease labels are identical, so continue comparison to versions
		return 1
	} else if candidate.PrereleaseVersion < other.PrereleaseVersion {
		return -1
	}
	// If we got this far, then the versions are clearly identical
	// (which is super weird)
	return 0
}

// String returns a string-serialized version of the SemverCandidate. The git
// ref metadata is excluded such that the output of this function resembles a
// semver-compliant version string.
func (candidate SemverCandidate) String() string {
	var buffer bytes.Buffer

	buffer.WriteString(strconv.Itoa(candidate.MajorVersion))
	buffer.WriteByte(SemverSelectorSeparatorChar)
	buffer.WriteString(strconv.Itoa(candidate.MinorVersion))
	buffer.WriteByte(SemverSelectorSeparatorChar)
	buffer.WriteString(strconv.Itoa(candidate.PatchVersion))

	if len(candidate.PrereleaseLabel) > 0 {
		buffer.WriteByte(SemverSelectorPrereleaseLabelPrefixChar)
		buffer.WriteString(candidate.PrereleaseLabel)

		if candidate.PrereleaseVersion > 0 {
			buffer.WriteByte(SemverSelectorSeparatorChar)
			buffer.WriteString(strconv.Itoa(candidate.PrereleaseVersion))
		}
	}

	return buffer.String()
}

// SemverCandidateList is an abstraction for a slice of SemverCandidates with
// some useful properties; namely, a SemverCandidateList is sortable and knows
// how to reduce itself to only matches that are relevant.
type SemverCandidateList []SemverCandidate

func (list SemverCandidateList) Len() int {
	return len(list)
}

func (list SemverCandidateList) Swap(i, j int) {
	list[i], list[j] = list[j], list[i]
}

func (list SemverCandidateList) Less(i, j int) bool {
	return list[i].CompareTo(list[j]) < 0
}

// Match returns a new SemverCandidateList with only candidates that match the
// specified selector.
func (list SemverCandidateList) Match(selector SemverSelector) SemverCandidateList {
	var newList []SemverCandidate

	for _, candidate := range list {
		if selector.Matches(candidate) {
			newList = append(newList, candidate)
		}
	}

	return newList
}

// Lowest returns the candidate that appaeared first in the list (which is by
// default the lowest).
func (list SemverCandidateList) Lowest() *SemverCandidate {
	listLength := len(list)
	if listLength < 1 {
		return nil
	}

	return &list[0]
}

// Highest returns the candidate that appaeared last in the list (which is by
// default the highest).
func (list SemverCandidateList) Highest() *SemverCandidate {
	listLength := len(list)
	if listLength < 1 {
		return nil
	}

	return &list[listLength-1]
}

// Best returns the best version available in the candidate list according to
// the specified selector.
func (list SemverCandidateList) Best(selector SemverSelector) *SemverCandidate {
	var (
		matches    = list.Match(selector)
		matchesLen = len(matches)
	)

	if matchesLen == 0 {
		// If there are no options, return no options.
		return nil
	} else if matchesLen == 1 {
		// Hmm, I wonder which *one* is the best :P.
		return &matches[0]
	} else {
		var (
			selectorHasLessThan  = selector.Suffix == SemverSelectorSuffixLessThan
			selectorHasWildcards = selector.MinorVersion.Type == SemverSegmentTypeWildcard ||
				selector.PatchVersion.Type == SemverSegmentTypeWildcard ||
				selector.PrereleaseVersion.Type == SemverSegmentTypeWildcard
		)

		// Get the most recent version available (adjusting what variation is
		// possible).
		if selectorHasWildcards || selectorHasLessThan {
			return matches.Highest()
		}

		return matches.Lowest()
	}
}
