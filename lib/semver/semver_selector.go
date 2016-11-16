package semver

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	// SemverSelectorTildeChar is the character that represents flexible patch
	// version selection.
	SemverSelectorTildeChar = '~'
	// SemverSelectorCaratChar is the character that represents flexible minor &
	// patch version selection.
	SemverSelectorCaratChar = '^'
	// SemverSelectorWildcardChar is the character that represents variable major,
	// minor, patch or pre-release version selection.
	SemverSelectorWildcardChar = 'x'
	// SemverSelectorLessThanChar is the character that represents the less-than
	// version inequality.
	SemverSelectorLessThanChar = '-'
	// SemverSelectorSeparatorChar is the character that separates segements of a
	// semver version.
	SemverSelectorSeparatorChar = '.'
	// SemverSelectorGreaterThanChar is the character that represents the
	// greater-than version inequality.
	SemverSelectorGreaterThanChar = '+'
	// SemverSelectorPrereleaseLabelPrefixChar is the character that separates the
	// patch version segment and the pre-release label.
	SemverSelectorPrereleaseLabelPrefixChar = '-'
)

const (
	errorSemverParseFailureInvalidSegment                  = "Invalid semver %s specified: %s"
	errorSemverParseFailureVersionTerminated               = "Could not parse the %s segment because the version was already complete"
	errorSemverParseFailureMissingMajorVersion             = "SemverSelector major segment was unspecified"
	errorSemverParseFailurePrefixMixedWithWildcard         = "Version prefixes cannot be mixed with version wildcards"
	errorSemverParseFailureSuffixMixedWithPrefixOrWildcard = "Version suffixes cannot be mixed with version wildcards or prefixes"
)

const (
	// SemverSelectorPrefixNone is the prefix enum value for an unspecified
	// prefix.
	SemverSelectorPrefixNone = iota
	// SemverSelectorPrefixTilde is the prefix enum value for a tilde prefix.
	SemverSelectorPrefixTilde = iota
	// SemverSelectorPrefixCarat is the prefix enum value for a carat prefix.
	SemverSelectorPrefixCarat = iota
)

const (
	// SemverSelectorSuffixNone is the suffix enum value for an unspecified
	// suffix.
	SemverSelectorSuffixNone = iota
	// SemverSelectorSuffixLessThan is the suffix enum value for a less-than
	// suffix.
	SemverSelectorSuffixLessThan = iota
	// SemverSelectorSuffixGreaterThan is the suffix enum value for a greater-than
	// suffix.
	SemverSelectorSuffixGreaterThan = iota
)

const (
	// SemverSegmentTypeNumber is the segment type enum value for segment of type
	// number.
	SemverSegmentTypeNumber = iota
	// SemverSegmentTypeWildcard is the segment type enum value for segment of
	// type wildcard.
	SemverSegmentTypeWildcard = iota
	// SemverSegmentTypeUnspecified is the segment type enum value for an
	// unspecified segment.
	SemverSegmentTypeUnspecified = iota
)

const (
	// SemverSegmentNamePrefix is the name of the prefix semver segment.
	SemverSegmentNamePrefix = "prefix"
	// SemverSegmentNameMajor is the name of the major semver segment.
	SemverSegmentNameMajor = "major"
	// SemverSegmentNameMinor is the name of the minor semver segment.
	SemverSegmentNameMinor = "minor"
	// SemverSegmentNamePatch is the name of the patch semver segment.
	SemverSegmentNamePatch = "patch"
	// SemverSegmentNamePrerelease is the name of the pre-release semver segment.
	SemverSegmentNamePrerelease = "pre-release"
	// SemverSegmentNameSuffix is the name of the suffix semver segment.
	SemverSegmentNameSuffix = "suffix"
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
		if prefix[0] == SemverSelectorTildeChar {
			semver.Prefix = SemverSelectorPrefixTilde
			semver.IsFlexible = true
		} else if prefix[0] == SemverSelectorCaratChar {
			semver.Prefix = SemverSelectorPrefixCarat
			semver.IsFlexible = true
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				SemverSegmentNamePrefix,
				prefix)
		}
	} else {
		semver.Prefix = SemverSelectorPrefixNone
	}

	if len(majorVersion) > 0 {
		if number, err := strconv.Atoi(majorVersion); err == nil {
			semver.MajorVersion.Type = SemverSegmentTypeNumber
			semver.MajorVersion.Number = number
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				SemverSegmentNameMajor,
				majorVersion)
		}
	} else {
		return semver, errors.New(errorSemverParseFailureMissingMajorVersion)
	}

	if len(minorVersion) > 0 {
		if strings.ToLower(minorVersion)[0] == SemverSelectorWildcardChar {
			if semver.Prefix == SemverSelectorPrefixNone {
				if !semver.IsFlexible {
					semver.IsFlexible = true
				}
				semver.MinorVersion.Type = SemverSegmentTypeWildcard
				versionCompleted = true
			} else {
				return semver, errors.New(
					errorSemverParseFailurePrefixMixedWithWildcard)
			}
		} else if number, err := strconv.Atoi(minorVersion); err == nil {
			semver.MinorVersion.Type = SemverSegmentTypeNumber
			semver.MinorVersion.Number = number
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureInvalidSegment,
				SemverSegmentNameMinor,
				minorVersion)
		}
	} else {
		semver.MinorVersion.Type = SemverSegmentTypeUnspecified
		versionCompleted = true
	}

	if len(patchVersion) > 0 {
		if !versionCompleted {
			if strings.ToLower(patchVersion)[0] == SemverSelectorWildcardChar {
				if semver.Prefix == SemverSelectorPrefixNone {
					if !semver.IsFlexible {
						semver.IsFlexible = true
					}
					semver.PatchVersion.Type = SemverSegmentTypeWildcard
					versionCompleted = true
				} else {
					return semver, errors.New(
						errorSemverParseFailurePrefixMixedWithWildcard)
				}
			} else if number, err := strconv.Atoi(patchVersion); err == nil {
				semver.PatchVersion.Type = SemverSegmentTypeNumber
				semver.PatchVersion.Number = number
			} else {
				return semver, fmt.Errorf(
					errorSemverParseFailureInvalidSegment,
					SemverSegmentNamePatch,
					patchVersion)
			}
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureVersionTerminated,
				SemverSegmentNamePatch)
		}
	} else {
		semver.PatchVersion.Type = SemverSegmentTypeUnspecified
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
				SemverSegmentNamePrerelease,
				prereleaseVersion)
		}
	} else {
		if !versionCompleted {
			versionCompleted = true
		}
	}

	if len(prereleaseVersion) > 0 {
		if !versionCompleted {
			if strings.ToLower(prereleaseVersion)[0] == SemverSelectorWildcardChar {
				if semver.Prefix == SemverSelectorPrefixNone {
					if !semver.IsFlexible {
						semver.IsFlexible = true
					}
					semver.PrereleaseVersion.Type = SemverSegmentTypeWildcard
				} else {
					return semver, errors.New(
						errorSemverParseFailurePrefixMixedWithWildcard)
				}
			} else if number, err := strconv.Atoi(prereleaseVersion); err == nil {
				semver.PrereleaseVersion.Type = SemverSegmentTypeNumber
				semver.PrereleaseVersion.Number = number
			} else {
				return semver, fmt.Errorf(
					errorSemverParseFailureInvalidSegment,
					SemverSegmentNamePrerelease,
					prereleaseVersion)
			}
		} else {
			return semver, fmt.Errorf(
				errorSemverParseFailureVersionTerminated,
				SemverSegmentNamePrerelease)
		}
	} else {
		semver.PrereleaseVersion.Type = SemverSegmentTypeUnspecified
	}

	if len(suffix) > 0 {
		if !semver.IsFlexible {
			if suffix[0] == SemverSelectorGreaterThanChar {
				semver.Suffix = SemverSelectorSuffixGreaterThan
				semver.IsFlexible = true
			} else if suffix[0] == SemverSelectorLessThanChar {
				semver.Suffix = SemverSelectorSuffixLessThan
				semver.IsFlexible = true
			} else {
				return semver, fmt.Errorf(
					errorSemverParseFailureInvalidSegment,
					SemverSegmentNameSuffix,
					suffix)
			}
		} else {
			return semver, errors.New(
				errorSemverParseFailureSuffixMixedWithPrefixOrWildcard)
		}
	} else {
		semver.Suffix = SemverSelectorSuffixNone
	}

	return semver, nil
}

// Matches simply determines whether the given candidate fits within the range
// defined by this version selector.
func (s SemverSelector) Matches(candidate SemverCandidate) bool {
	if s.IsFlexible {
		if s.Suffix == SemverSelectorSuffixGreaterThan {
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
			} else if s.PrereleaseVersion.Type == SemverSegmentTypeNumber {
				// The fact that we've gotten this far means that the pre-release labels
				// match - we just need to check that the version itself is greater
				return candidate.PrereleaseVersion >= s.PrereleaseVersion.Number
			} else {
				// If we got this far, the pre-release version of the selector was
				// unspecified, so, since anything is greater than nothing, we default
				// to true
				return true
			}
		} else if s.Suffix == SemverSelectorSuffixLessThan {
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
			} else if s.PrereleaseVersion.Type == SemverSegmentTypeNumber {
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
		} else if s.Prefix == SemverSelectorPrefixCarat {
			if s.MajorVersion.Number != candidate.MajorVersion {
				return false
			} else if s.MinorVersion.Number > candidate.MinorVersion {
				return false
			} else if s.PatchVersion.Number > candidate.PatchVersion {
				return false
			}

			return len(candidate.PrereleaseLabel) == 0
		} else if s.Prefix == SemverSelectorPrefixTilde {
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
			case SemverSegmentTypeWildcard, SemverSegmentTypeUnspecified:
				return true
			}
			switch s.PatchVersion.Type {
			case SemverSegmentTypeWildcard, SemverSegmentTypeUnspecified:
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
				if s.PrereleaseVersion.Type == SemverSegmentTypeUnspecified {
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

	if s.Prefix == SemverSelectorPrefixTilde {
		buffer.WriteByte(SemverSelectorTildeChar)
	} else if s.Prefix == SemverSelectorPrefixCarat {
		buffer.WriteByte(SemverSelectorCaratChar)
	}

	if s.MajorVersion.Type == SemverSegmentTypeNumber {
		buffer.WriteString(strconv.Itoa(s.MajorVersion.Number))
	} else {
		var majorValue string

		if s.MajorVersion.Type == SemverSegmentTypeWildcard {
			majorValue = "a wildcard"
		} else {
			majorValue = "unspecified"
		}

		panic(
			fmt.Sprintf("Cannot stringify invalid semver (major is %s)", majorValue))
	}

	if s.MinorVersion.Type == SemverSegmentTypeNumber {
		buffer.WriteByte(SemverSelectorSeparatorChar)
		buffer.WriteString(strconv.Itoa(s.MinorVersion.Number))
	} else if s.MinorVersion.Type == SemverSegmentTypeWildcard {
		buffer.WriteByte(SemverSelectorSeparatorChar)
		buffer.WriteByte(SemverSelectorWildcardChar)
	} else {
		versionStringCompleted = true
	}

	if !versionStringCompleted {
		if s.PatchVersion.Type == SemverSegmentTypeNumber {
			buffer.WriteByte(SemverSelectorSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.PatchVersion.Number))
		} else if s.PatchVersion.Type == SemverSegmentTypeWildcard {
			buffer.WriteByte(SemverSelectorSeparatorChar)
			buffer.WriteByte(SemverSelectorWildcardChar)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if len(s.PrereleaseLabel) > 0 {
			buffer.WriteByte(SemverSelectorPrereleaseLabelPrefixChar)
			buffer.WriteString(s.PrereleaseLabel)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if s.PrereleaseVersion.Type == SemverSegmentTypeNumber {
			buffer.WriteByte(SemverSelectorSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.PrereleaseVersion.Number))
		} else if s.PrereleaseVersion.Type == SemverSegmentTypeWildcard {
			buffer.WriteByte(SemverSelectorSeparatorChar)
			buffer.WriteByte(SemverSelectorWildcardChar)
		}
	}

	if s.Suffix == SemverSelectorSuffixLessThan {
		buffer.WriteByte(SemverSelectorLessThanChar)
	} else if s.Suffix == SemverSelectorSuffixGreaterThan {
		buffer.WriteByte(SemverSelectorGreaterThanChar)
	}

	return buffer.String()
}
