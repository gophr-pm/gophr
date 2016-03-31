package main

import (
	"bytes"
	"fmt"
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

type SemverInequality int

const (
	SemverInequalityNone        = iota
	SemverInequalityLessThan    = iota
	SemverInequalityGreaterThan = iota
)

type SemverPartType int

const (
	SemverPartTypeNumber      = iota
	SemverPartTypeWildcard    = iota
	SemverPartTypeUnspecified = iota
)

type SemverPart struct {
	Type   SemverPartType
	Number int
}

type Semver struct {
	HasTilde          bool
	HasCarat          bool
	MajorVersion      SemverPart
	MinorVersion      SemverPart
	PatchVersion      SemverPart
	PrereleaseLabel   string
	PrereleaseVersion SemverPart
	VersionInequality SemverInequality
}

func (s Semver) String() {
	var (
		buffer                 bytes.Buffer
		versionStringCompleted = false
	)

	if s.HasTilde {
		buffer.WriteByte(semverTildeChar)
	} else if s.HasCarat {
		buffer.WriteByte(semverCaratChar)
	}

	if s.MajorVersion.Type == SemverPartTypeNumber {
		buffer.WriteString(strcons.Itoa(s.MajorVersion.Number))
	} else {
		var majorValue string

		if s.MajorVersion.Type == SemverPartTypeWildcard {
			majorValue = "a wildcard"
		} else {
			majorValue = "unspecified"
		}

		panic(fmt.Sprintf("Cannot stringify invalid semver (major is %s)", majorValue))
	}

	if s.Minor.Type == SemverPartTypeNumber {
		buffer.WriteByte(semverSeparatorChar)
		buffer.WriteString(strcons.Itoa(s.Minor.Number))
	} else if s.Minor.Type == SemverPartTypeWildcard {
		buffer.WriteByte(semverSeparatorChar)
		buffer.WriteByte(semverWildcardChar)
	} else {
		versionStringCompleted = true
	}

	if !versionStringCompleted {
		if s.Patch.Type == SemverPartTypeNumber {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteString(strcons.Itoa(s.Patch.Number))
		} else if s.Patch.Type == SemverPartTypeWildcard {
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
		if s.Prerelease.Type == SemverPartTypeNumber {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteString(strcons.Itoa(s.Prerelease.Number))
		} else if s.Prerelease.Type == SemverPartTypeWildcard {
			buffer.WriteByte(semverSeparatorChar)
			buffer.WriteByte(semverWildcardChar)
		}
	}

	if s.VersionInequality == SemverInequalityLessThan {
		buffer.WriteByte(semverLessThanChar)
	} else if s.VersionInequality == SemverInequalityGreaterThan {
		buffer.WriteByte(semverGreaterThanChar)
	}

	return buffer.String()
}
