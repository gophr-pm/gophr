package main

import (
	"bytes"
	"errors"
	"fmt"
	"strconv"
	"strings"
)

const (
	semverTildeChar                 = '~'
	semverCaratChar                 = '^'
	semverWildcardChar              = 'x'
	semverLessThanChar              = '-'
	semverSeparatorChar             = '.'
	semverGreaterThanChar           = '+'
	semverPrereleaseLabelPrefixChar = '-'
)

const (
	errorParseFailureInvalidSegment                  = "Invalid semver %s specified: %s"
	errorParseFailureVersionTerminated               = "Could not parse the %s segment because the version was already complete"
	errorParseFailureMissingMajorVersion             = "Semver major segment was unspecified"
	errorParseFailurePrefixMixedWithWildcard         = "Version prefixes cannot be mixed with version wildcards"
	errorParseFailureSuffixMixedWithPrefixOrWildcard = "Version suffixes cannot be mixed with version wildcards or prefixes"
)

const (
	semverPrefixNone  = iota
	semverPrefixTilde = iota
	semverPrefixCarat = iota
)

const (
	semverSuffixNone        = iota
	semverSuffixLessThan    = iota
	semverSuffixGreaterThan = iota
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

type SemverSegment struct {
	Type   int
	Number int
}

type Semver struct {
	Prefix            int
	Suffix            int
	IsFlexible        bool
	MajorVersion      SemverSegment
	MinorVersion      SemverSegment
	PatchVersion      SemverSegment
	PrereleaseLabel   string
	PrereleaseVersion SemverSegment
}

func NewSemver(
	prefix string,
	majorVersion string,
	minorVersion string,
	patchVersion string,
	prereleaseLabel string,
	prereleaseVersion string,
	suffix string,
) (Semver, error) {
	// TODO(skeswa): implement this with full validation
	var (
		semver           Semver
		versionCompleted = false
	)

	if len(prefix) > 0 {
		if prefix[0] == semverTildeChar {
			semver.Prefix = semverPrefixTilde
			semver.IsFlexible = true
		} else if prefix[0] == semverCaratChar {
			semver.Prefix = semverPrefixCarat
			semver.IsFlexible = true
		} else {
			return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNamePrefix, prefix)
		}
	} else {
		semver.Prefix = semverPrefixNone
	}

	if len(majorVersion) > 0 {
		if number, err := strconv.Atoi(majorVersion); err == nil {
			semver.MajorVersion.Type = semverSegmentTypeNumber
			semver.MajorVersion.Number = number
		} else {
			return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNameMajor, majorVersion)
		}
	} else {
		return semver, errors.New(errorParseFailureMissingMajorVersion)
	}

	if len(minorVersion) > 0 {
		if strings.ToLower(minorVersion)[0] == semverWildcardChar {
			if semver.Prefix == semverPrefixNone {
				if !semver.IsFlexible {
					semver.IsFlexible = true
				}
				semver.MinorVersion.Type = semverSegmentTypeWildcard
				versionCompleted = true
			} else {
				return semver, errors.New(errorParseFailurePrefixMixedWithWildcard)
			}
		} else if number, err := strconv.Atoi(minorVersion); err == nil {
			semver.MinorVersion.Type = semverSegmentTypeNumber
			semver.MinorVersion.Number = number
		} else {
			return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNameMinor, minorVersion)
		}
	} else {
		semver.MinorVersion.Type = semverSegmentTypeUnspecified
		versionCompleted = true
	}

	if len(patchVersion) > 0 {
		if !versionCompleted {
			if strings.ToLower(patchVersion)[0] == semverWildcardChar {
				if semver.Prefix == semverPrefixNone {
					if !semver.IsFlexible {
						semver.IsFlexible = true
					}
					semver.PatchVersion.Type = semverSegmentTypeWildcard
					versionCompleted = true
				} else {
					return semver, errors.New(errorParseFailurePrefixMixedWithWildcard)
				}
			} else if number, err := strconv.Atoi(patchVersion); err == nil {
				semver.PatchVersion.Type = semverSegmentTypeNumber
				semver.PatchVersion.Number = number
			} else {
				return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNamePatch, patchVersion)
			}
		} else {
			return semver, fmt.Errorf(errorParseFailureVersionTerminated, semverSegmentNamePatch)
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
			return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNamePrerelease, prereleaseVersion)
		}
	}

	if len(prereleaseVersion) > 0 {
		if !versionCompleted {
			if strings.ToLower(prereleaseVersion)[0] == semverWildcardChar {
				if semver.Prefix == semverPrefixNone {
					if !semver.IsFlexible {
						semver.IsFlexible = true
					}
					semver.PrereleaseVersion.Type = semverSegmentTypeWildcard
				} else {
					return semver, errors.New(errorParseFailurePrefixMixedWithWildcard)
				}
			} else if number, err := strconv.Atoi(prereleaseVersion); err == nil {
				semver.PrereleaseVersion.Type = semverSegmentTypeNumber
				semver.PrereleaseVersion.Number = number
			} else {
				return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNamePrerelease, prereleaseVersion)
			}
		} else {
			return semver, fmt.Errorf(errorParseFailureVersionTerminated, semverSegmentNamePrerelease)
		}
	} else {
		semver.PrereleaseVersion.Type = semverSegmentTypeUnspecified
	}

	if len(suffix) > 0 {
		if !semver.IsFlexible {
			if suffix[0] == semverGreaterThanChar {
				semver.Suffix = semverSuffixGreaterThan
			} else if suffix[0] == semverLessThanChar {
				semver.Suffix = semverSuffixLessThan
			} else {
				return semver, fmt.Errorf(errorParseFailureInvalidSegment, semverSegmentNameSuffix, suffix)
			}
		} else {
			return semver, errors.New(errorParseFailureSuffixMixedWithPrefixOrWildcard)
		}
	} else {
		semver.Suffix = semverSuffixNone
	}

	return semver, nil
}

func (s Semver) String() string {
	var (
		buffer                 bytes.Buffer
		versionStringCompleted = false
	)

	if s.Prefix == semverPrefixTilde {
		buffer.WriteByte(semverTildeChar)
	} else if s.Prefix == semverPrefixCarat {
		buffer.WriteByte(semverCaratChar)
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

		panic(fmt.Sprintf("Cannot stringify invalid semver (major is %s)", majorValue))
	}

	if s.MinorVersion.Type == semverSegmentTypeNumber {
		buffer.WriteByte(semverSeparatorChar)
		buffer.WriteString(strconv.Itoa(s.MinorVersion.Number))
	} else if s.MinorVersion.Type == semverSegmentTypeWildcard {
		buffer.WriteByte(semverSeparatorChar)
		buffer.WriteByte(semverWildcardChar)
	} else {
		versionStringCompleted = true
	}

	if !versionStringCompleted {
		if s.PatchVersion.Type == semverSegmentTypeNumber {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.PatchVersion.Number))
		} else if s.PatchVersion.Type == semverSegmentTypeWildcard {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteByte(semverWildcardChar)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if len(s.PrereleaseLabel) > 0 {
			buffer.WriteByte(semverPrereleaseLabelPrefixChar)
			buffer.WriteString(s.PrereleaseLabel)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if s.PrereleaseVersion.Type == semverSegmentTypeNumber {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.PrereleaseVersion.Number))
		} else if s.PrereleaseVersion.Type == semverSegmentTypeWildcard {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteByte(semverWildcardChar)
		}
	}

	if s.Suffix == semverSuffixLessThan {
		buffer.WriteByte(semverLessThanChar)
	} else if s.Suffix == semverSuffixGreaterThan {
		buffer.WriteByte(semverGreaterThanChar)
	}

	return buffer.String()
}
