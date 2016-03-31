package main

import (
	"bytes"
	"fmt"
	"strconv"
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

type SemverPartType int

const (
	semverPartTypeNumber      = iota
	semverPartTypeWildcard    = iota
	semverPartTypeUnspecified = iota
)

type SemverPrefix int

const (
	semverPrefixNone  = iota
	semverPrefixTilde = iota
	semverPrefixCarat = iota
)

type SemverSuffix int

const (
	semverSuffixNone        = iota
	semverSuffixLessThan    = iota
	semverSuffixGreaterThan = iota
)

type SemverPart struct {
	Type   SemverPartType
	Number int
}

type Semver struct {
	Prefix            SemverPrefix
	MajorVersion      SemverPart
	MinorVersion      SemverPart
	PatchVersion      SemverPart
	PrereleaseLabel   string
	PrereleaseVersion SemverPart
	Suffix            SemverSuffix
}

func NewSemver(prefix string, majorVersion string, minorVersion string, patchVersion string, prereleaseLabel string, prereleaseVersion string, suffix string) (Semver, error) {
	// TODO(skeswa): implement this with full validation
}

func (s Semver) String() {
	var (
		buffer                 bytes.Buffer
		versionStringCompleted = false
	)

	if s.prefix == semverPrefixTilde {
		buffer.WriteByte(semverTildeChar)
	} else if s.prefix == semverPrefixCarat {
		buffer.WriteByte(semverCaratChar)
	}

	if s.majorVersion.Type == semverPartTypeNumber {
		buffer.WriteString(strconv.Itoa(s.majorVersion.Number))
	} else {
		var majorValue string

		if s.majorVersion.Type == semverPartTypeWildcard {
			majorValue = "a wildcard"
		} else {
			majorValue = "unspecified"
		}

		panic(fmt.Sprintf("Cannot stringify invalid semver (major is %s)", majorValue))
	}

	if s.minorVersion.Type == semverPartTypeNumber {
		buffer.WriteByte(semverSeparatorChar)
		buffer.WriteString(strconv.Itoa(s.minorVersion.Number))
	} else if s.minorVersion.Type == semverPartTypeWildcard {
		buffer.WriteByte(semverSeparatorChar)
		buffer.WriteByte(semverWildcardChar)
	} else {
		versionStringCompleted = true
	}

	if !versionStringCompleted {
		if s.patchVersion.Type == semverPartTypeNumber {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.Patch.Number))
		} else if s.patchVersion.Type == semverPartTypeWildcard {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteByte(semverWildcardChar)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if len(s.prereleaseLabel) > 0 {
			buffer.WriteByte(semverPrereleaseLabelPrefixChar)
			buffer.WriteString(s.prereleaseLabel)
		} else {
			versionStringCompleted = true
		}
	}

	if !versionStringCompleted {
		if s.prerelease.Type == semverPartTypeNumber {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteString(strconv.Itoa(s.prerelease.Number))
		} else if s.prerelease.Type == semverPartTypeWildcard {
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
