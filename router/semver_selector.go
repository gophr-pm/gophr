package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	semverSelectorTildeChar                 = '~'
	semverSelectorCaratChar                 = '^'
	semverSelectorWildcardChar              = 'x'
	semverSelectorLessThanChar              = '-'
	semverSelectorSeparatorChar             = '.'
	semverSelectorGreaterThanChar           = '+'
	semverSelectorPrereleaseLabelPrefixChar = '-'
)

const (
	errorSemverParseFailureInvalidSegment                  = "Invalid semver %s specified: %s"
	errorSemverParseFailureVersionTerminated               = "Could not parse the %s segment because the version was already complete"
	errorSemverParseFailureMissingMajorVersion             = "SemverSelector major segment was unspecified"
	errorSemverParseFailurePrefixMixedWithWildcard         = "Version prefixes cannot be mixed with version wildcards"
	errorSemverParseFailureSuffixMixedWithPrefixOrWildcard = "Version suffixes cannot be mixed with version wildcards or prefixes"
)

const (
	semverSelectorPrefixNone  = iota
	semverSelectorPrefixTilde = iota
	semverSelectorPrefixCarat = iota
)

const (
	semverSelectorSuffixNone        = iota
	semverSelectorSuffixLessThan    = iota
	semverSelectorSuffixGreaterThan = iota
)

const (
	semverSegmentTypeNumber      = iota
	semverSegmentTypeWildcard    = iota
	semverSegmentTypeUnspecified = iota
)

const (
	semverSegmentNamePrefix     = "prefix"
	semverSegmentNameMajor      = "major"
	semverSegmentNameMinor      = "minor"
	semverSegmentNamePatch      = "patch"
	semverSegmentNamePrerelease = "pre-release"
	semverSegmentNameSuffix     = "suffix"
)

// SemverSelectorSegment is the atomic unit of a semver version.
type SemverSelectorSegment struct {
	Type   int
	Number int
}

// SemverSelector is a semver version selector. It can either specify a range of
// versions that it matches or refer to one specific version.
type SemverSelector struct {
	Prefix            int
	Suffix            int
	IsFlexible        bool
	MajorVersion      SemverSelectorSegment
	MinorVersion      SemverSelectorSegment
	PatchVersion      SemverSelectorSegment
	PrereleaseLabel   string
	PrereleaseVersion SemverSelectorSegment
}

// NewSemverSelector creates a new semver version from the version selector
// regular expression capture groups.
func NewSemverSelector(
	prefix string,
	majorVersion string,
	minorVersion string,
	patchVersion string,
	prereleaseLabel string,
	prereleaseVersion string,
	suffix string) (SemverSelector, error) {
	// TODO(skeswa): implement this with full validation
	var (
		semver           SemverSelector
		versionCompleted = false
	)

	if len(prefix) > 0 {
		if prefix[0] == semverSelectorTildeChar {
			semver.Prefix = semverSelectorPrefixTilde
			semver.IsFlexible = true
		} else if prefix[0] == semverSelectorCaratChar {
			semver.Prefix = semverSelectorPrefixCarat
			semver.IsFlexible = true
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				semverSegmentNamePrefix,
				prefix)
		}
	} else {
		semver.Prefix = semverSelectorPrefixNone
	}

	if len(majorVersion) > 0 {
		if number, err := strconv.Atoi(majorVersion); err == nil {
			semver.MajorVersion.Type = semverSegmentTypeNumber
			semver.MajorVersion.Number = number
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				semverSegmentNameMajor,
				majorVersion)
		}
	} else {
		return semver, errors.New(errorSemverParseFailureMissingMajorVersion)
	}

	if len(minorVersion) > 0 {
		if strings.ToLower(minorVersion)[0] == semverSelectorWildcardChar {
			if semver.Prefix == semverSelectorPrefixNone {
				if !semver.IsFlexible {
					semver.IsFlexible = true
				}
				semver.MinorVersion.Type = semverSegmentTypeWildcard
				versionCompleted = true
			} else {
				return semver, errors.New(
					errorSemverParseFailurePrefixMixedWithWildcard)
			}
		} else if number, err := strconv.Atoi(minorVersion); err == nil {
			semver.MinorVersion.Type = semverSegmentTypeNumber
			semver.MinorVersion.Number = number
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				semverSegmentNameMinor,
				minorVersion)
		}
	} else {
		semver.MinorVersion.Type = semverSegmentTypeUnspecified
		versionCompleted = true
	}

	if len(patchVersion) > 0 {
		if !versionCompleted {
			if strings.ToLower(patchVersion)[0] == semverSelectorWildcardChar {
				if semver.Prefix == semverSelectorPrefixNone {
					if !semver.IsFlexible {
						semver.IsFlexible = true
					}
					semver.PatchVersion.Type = semverSegmentTypeWildcard
					versionCompleted = true
				} else {
					return semver, errors.New(
						errorSemverParseFailurePrefixMixedWithWildcard)
				}
			} else if number, err := strconv.Atoi(patchVersion); err == nil {
				semver.PatchVersion.Type = semverSegmentTypeNumber
				semver.PatchVersion.Number = number
			} else {
				return semver, fmt.Errorf(
					errorSemverParseFailureInvalidSegment,
					semverSegmentNamePatch,
					patchVersion)
			}
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureVersionTerminated,
				semverSegmentNamePatch)
		}
	} else {
		semver.PatchVersion.Type = semverSegmentTypeUnspecified
		if !versionCompleted {
			versionCompleted = true
		}
	}

	if len(prereleaseLabel) > 0 {
		if !versionCompleted {
			semver.PrereleaseLabel = prereleaseLabel
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				semverSegmentNamePrerelease,
				prereleaseVersion)
		}
	} else {
		if !versionCompleted {
			versionCompleted = true
		}
	}

	if len(prereleaseVersion) > 0 {
		if !versionCompleted {
			if strings.ToLower(prereleaseVersion)[0] == semverSelectorWildcardChar {
				if semver.Prefix == semverSelectorPrefixNone {
					if !semver.IsFlexible {
						semver.IsFlexible = true
					}
					semver.PrereleaseVersion.Type = semverSegmentTypeWildcard
				} else {
					return semver, errors.New(
						errorSemverParseFailurePrefixMixedWithWildcard)
				}
			} else if number, err := strconv.Atoi(prereleaseVersion); err == nil {
				semver.PrereleaseVersion.Type = semverSegmentTypeNumber
				semver.PrereleaseVersion.Number = number
			} else {
				return semver, fmt.Errorf(
					errorSemverParseFailureInvalidSegment,
					semverSegmentNamePrerelease,
					prereleaseVersion)
			}
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureVersionTerminated,
				semverSegmentNamePrerelease)
		}
	} else {
		semver.PrereleaseVersion.Type = semverSegmentTypeUnspecified
	}

	if len(suffix) > 0 {
		if !semver.IsFlexible {
			if suffix[0] == semverSelectorGreaterThanChar {
				semver.Suffix = semverSelectorSuffixGreaterThan
				semver.IsFlexible = true
			} else if suffix[0] == semverSelectorLessThanChar {
				semver.Suffix = semverSelectorSuffixLessThan
				semver.IsFlexible = true
			} else {
				return semver, fmt.Errorf(
					errorSemverParseFailureInvalidSegment,
					semverSegmentNameSuffix,
					suffix)
			}
		} else {
			return semver, errors.New(
				errorSemverParseFailureSuffixMixedWithPrefixOrWildcard)
		}
	} else {
		semver.Suffix = semverSelectorSuffixNone
	}

	return semver, nil
}

// Matches simply determines whether the given candidate fits within the range
// defined by this version selector.
func (s SemverSelector) Matches(candidate SemverCandidate) bool {
	if s.IsFlexible {
		if s.Suffix == semverSelectorSuffixGreaterThan {
			if s.MajorVersion.Number > candidate.MajorVersion {
				return false
			} else if s.MajorVersion.Number < candidate.MajorVersion {
				return true
			} else if s.MinorVersion.Number > candidate.MinorVersion {
				return false
			} else if s.MinorVersion.Number < candidate.MinorVersion {
				return true
			} else if s.PatchVersion.Number > candidate.PatchVersion {
				return false
			} else if s.PatchVersion.Number < candidate.PatchVersion {
				return true
			} else if len(s.PrereleaseLabel) == 0 && len(candidate.PrereleaseLabel) > 0 {
				// Don't match a pre-release candidate if possible
				return false
			} else if len(s.PrereleaseLabel) > 0 && len(candidate.PrereleaseLabel) == 0 {
				// If the selector has a pre-release, and the candidate doesn't, then
				// we can conclude that it is greater
				return true
			} else if s.PrereleaseLabel != candidate.PrereleaseLabel {
				// If the selector's pre-release doesn't match the candidate's
				// pre-release, then they should not match
				return false
			} else if s.PrereleaseVersion.Type == semverSegmentTypeNumber {
				// The fact that we've gotten this far means that the pre-release labels
				// match - we just need to check that the version itself is greater
				return candidate.PrereleaseVersion >= s.PrereleaseVersion.Number
			} else {
				// If we got this far, the pre-release version of the selector was
				// unspecified, so, since anything is greater than nothing, we default
				// to true
				return true
			}
		} else if s.Suffix == semverSelectorSuffixLessThan {
			if s.MajorVersion.Number > candidate.MajorVersion {
				return true
			} else if s.MajorVersion.Number < candidate.MajorVersion {
				return false
			} else if s.MinorVersion.Number > candidate.MinorVersion {
				return true
			} else if s.MinorVersion.Number < candidate.MinorVersion {
				return false
			} else if s.PatchVersion.Number > candidate.PatchVersion {
				return true
			} else if s.PatchVersion.Number < candidate.PatchVersion {
				return false
			} else if len(s.PrereleaseLabel) == 0 && len(candidate.PrereleaseLabel) > 0 {
				// Don't match a pre-release candidate if possible
				return false
			} else if len(s.PrereleaseLabel) > 0 && len(candidate.PrereleaseLabel) == 0 {
				// If the selector has a pre-release, and the candidate doesn't, then
				// the candidate probably doesn't match given that pre-release means the
				// version is "less than" one without a pre-release
				return false
			} else if s.PrereleaseLabel != candidate.PrereleaseLabel {
				// If the selector's pre-release doesn't match the candidate's
				// pre-release, then they should not match
				return false
			} else if s.PrereleaseVersion.Type == semverSegmentTypeNumber {
				// The fact that we've gotten this far means that the pre-release labels
				// match - we just need to check that the version itself is lesser
				return candidate.PrereleaseVersion <= s.PrereleaseVersion.Number
			} else {
				// If we got this far, the pre-release version of the selector was
				// unspecified. We treat that virtually as pre-release version 0. Since
				// only a candidate with pre-release version 0 could match, we'll
				// make it the return condition.
				return candidate.PrereleaseVersion == 0
			}
		} else if s.Prefix == semverSelectorPrefixCarat {
			if s.MajorVersion.Number != candidate.MajorVersion {
				return false
			} else if s.MinorVersion.Number > candidate.MinorVersion {
				return false
			} else if s.PatchVersion.Number > candidate.PatchVersion {
				return false
			}

			return len(candidate.PrereleaseLabel) == 0
		} else if s.Prefix == semverSelectorPrefixTilde {
			if s.MajorVersion.Number != candidate.MajorVersion {
				return false
			} else if s.MinorVersion.Number != candidate.MinorVersion {
				return false
			} else if s.PatchVersion.Number > candidate.PatchVersion {
				return false
			}

			return len(candidate.PrereleaseLabel) == 0
		} else {
			// This means that we have at least one wildcard
			if s.MajorVersion.Number != candidate.MajorVersion {
				return false
			}
			switch s.MinorVersion.Type {
			case semverSegmentTypeWildcard, semverSegmentTypeUnspecified:
				return true
			}
			switch s.PatchVersion.Type {
			case semverSegmentTypeWildcard, semverSegmentTypeUnspecified:
				return true
			}

			return s.PrereleaseLabel == candidate.PrereleaseLabel
		}
	} else {
		primaryVersionsMatch := s.MajorVersion.Number == candidate.MajorVersion &&
			s.MinorVersion.Number == candidate.MinorVersion &&
			s.PatchVersion.Number == candidate.PatchVersion

		if len(s.PrereleaseLabel) > 0 {
			matchesUpToLabel := primaryVersionsMatch &&
				s.PrereleaseLabel == candidate.PrereleaseLabel

			if matchesUpToLabel {
				if s.PrereleaseVersion.Type == semverSegmentTypeUnspecified {
					return true
				}

				return s.PrereleaseVersion.Number == candidate.PrereleaseVersion
			}

			return false
		}

		return primaryVersionsMatch
	}
}

func (s SemverSelector) String() string {
	var (
		buffer                 bytes.Buffer
		versionStringCompleted = false
	)

	if s.Prefix == semverSelectorPrefixTilde {
		buffer.WriteByte(semverSelectorTildeChar)
	} else if s.Prefix == semverSelectorPrefixCarat {
		buffer.WriteByte(semverSelectorCaratChar)
	}

	if s.MajorVersion.Type == semverSegmentTypeNumber {
		buffer.WriteString(strconv.Itoa(s.MajorVersion.Number))
	} else {
		var majorValue string

		if s.MajorVersion.Type == semverSegmentTypeWildcard {
			majorValue = "a wildcard"
		} else {
			majorValue = "unspecified"
		}

		panic(
			fmt.Sprintf("Cannot stringify invalid semver (major is %s)", majorValue))
	}

	if s.MinorVersion.Type == semverSegmentTypeNumber {
		buffer.WriteByte(semverSelectorSeparatorChar)
		buffer.WriteString(strconv.Itoa(s.MinorVersion.Number))
	} else if s.MinorVersion.Type == semverSegmentTypeWildcard {
		buffer.WriteByte(semverSelectorSeparatorChar)
		buffer.WriteByte(semverSelectorWildcardChar)
	} else {
		versionStringCompleted = true
	}

	if !versionStringCompleted {
		if s.PatchVersion.Type == semverSegmentTypeNumber {
			buffer.WriteByte(semverSelectorSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.PatchVersion.Number))
		} else if s.PatchVersion.Type == semverSegmentTypeWildcard {
			buffer.WriteByte(semverSelectorSeparatorChar)
			buffer.WriteByte(semverSelectorWildcardChar)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if len(s.PrereleaseLabel) > 0 {
			buffer.WriteByte(semverSelectorPrereleaseLabelPrefixChar)
			buffer.WriteString(s.PrereleaseLabel)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if s.PrereleaseVersion.Type == semverSegmentTypeNumber {
			buffer.WriteByte(semverSelectorSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.PrereleaseVersion.Number))
		} else if s.PrereleaseVersion.Type == semverSegmentTypeWildcard {
			buffer.WriteByte(semverSelectorSeparatorChar)
			buffer.WriteByte(semverSelectorWildcardChar)
		}
	}

	if s.Suffix == semverSelectorSuffixLessThan {
		buffer.WriteByte(semverSelectorLessThanChar)
	} else if s.Suffix == semverSelectorSuffixGreaterThan {
		buffer.WriteByte(semverSelectorGreaterThanChar)
	}

	return buffer.String()
}
