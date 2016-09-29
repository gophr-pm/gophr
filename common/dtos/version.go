package dtos

import (
	"bytes"

	"github.com/gophr-pm/gophr/common/semver"
)

//go:generate ffjson $GOFILE

const (
	// VersionDTOTypeSemver is the type of version DTO representing a semver
	// version.
	VersionDTOTypeSemver = "semver"
	// VersionDTOTypeRefHash is the type of version DTO representing a git ref.
	VersionDTOTypeRefHash = "refhash"
)

const (
	refHashLength = 40
)

// VersionDTO is the data transfer object for a package version (either semver
// or git ref).
type VersionDTO struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

// NewVersionDTO creates a new VersionDTO.
func NewVersionDTO(versionStr string) *VersionDTO {
	if len(versionStr) == refHashLength {
		return &VersionDTO{
			Type:  VersionDTOTypeRefHash,
			Value: versionStr,
		}
	}

	return &VersionDTO{
		Type:  VersionDTOTypeSemver,
		Value: versionStr,
	}
}

// VersionListDTO is the data transfer object for a list of VersionDTOs.
type VersionListDTO []*VersionDTO

// NewVersionListDTO creates a new VersionListDTO.
func NewVersionListDTO() VersionListDTO {
	return make(VersionListDTO, 0)
}

// NewVersionListDTOFromVersionStrings builds a VersionListDTO from a list of
// version strings.
func NewVersionListDTOFromVersionStrings(versionStrings []string) VersionListDTO {
	versionDTOs := NewVersionListDTO()

	for _, versionStr := range versionStrings {
		versionDTOs = append(versionDTOs, NewVersionDTO(versionStr))
	}

	return versionDTOs
}

// NewVersionListDTOFromSemverCandidateList builds a VersionListDTO from a list
// of version candidates.
func NewVersionListDTOFromSemverCandidateList(candidates semver.SemverCandidateList) VersionListDTO {
	versionDTOs := NewVersionListDTO()

	for _, candidate := range candidates {
		versionDTOs = append(versionDTOs, NewVersionDTO(candidate.String()))
	}

	return versionDTOs
}

// MarshalJSON returns the JSON encoding of the VersionListDTO.
func (dto VersionListDTO) MarshalJSON() ([]byte, error) {
	var (
		err    error
		json   []byte
		buffer bytes.Buffer
	)

	buffer.WriteByte(JSONCharOpenArray)
	for i, item := range dto {
		json, err = item.MarshalJSON()
		if err != nil {
			return nil, err
		}

		if i > 0 {
			buffer.WriteByte(JSONCharArrayDelimiter)
		}

		buffer.Write(json)
	}
	buffer.WriteByte(JSONCharCloseArray)

	return buffer.Bytes(), nil
}
