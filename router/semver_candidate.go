package main

import (
	"errors"
	"strconv"
	"strings"
)

type SemverCandidate struct {
	GitRefHash              string
	GitRefName              string
	MajorVersion            int
	MinorVersion            int
	PatchVersion            int
	PrereleaseLabel         string
	PrereleaseVersion       int
	PrereleaseVersionExists bool
}

type SemverCandidateList []SemverCandidate

func NewSemverCandidate(
	gitRefHash string,
	gitRefName string,
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
		MajorVersion:            majorVersionNumber,
		MinorVersion:            minorVersionNumber,
		PatchVersion:            patchVersionNumber,
		PrereleaseLabel:         prereleaseLabel,
		PrereleaseVersion:       prereleaseVersionNumber,
		PrereleaseVersionExists: (len(prereleaseLabel) > 0),
	}, nil
}

func (this SemverCandidate) CompareTo(that SemverCandidate) int {
	if this.MajorVersion > that.MajorVersion {
		return 1
	} else if this.MajorVersion < that.MajorVersion {
		return -1
	} else if this.MinorVersion > that.MinorVersion {
		return 1
	} else if this.MinorVersion < that.MinorVersion {
		return -1
	} else if this.PatchVersion > that.PatchVersion {
		return 1
	} else if this.PatchVersion < that.PatchVersion {
		return -1
	} else if len(this.PrereleaseLabel) > 0 && len(that.PrereleaseLabel) == 0 {
		// Prerelease immediately means that the version is lesser
		return -1
	} else if len(this.PrereleaseLabel) == 0 && len(that.PrereleaseLabel) > 0 {
		// Prerelease immediately means that the version is lesser
		return 1
	} else if this.PrereleaseLabel != that.PrereleaseLabel {
		// Prerelease labels don't match, so return a comparison
		return strings.Compare(this.PrereleaseLabel, that.PrereleaseLabel)
	} else if this.PrereleaseVersion > that.PrereleaseVersion {
		// Prerelease labels are identical, so continue comparison to versions
		return 1
	} else if this.PrereleaseVersion < that.PrereleaseVersion {
		return -1
	}
	// If we got this far, then the versions are clearly identical
	// (which is super weird)
	return 0
}
