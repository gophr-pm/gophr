package common

import "github.com/pquerna/ffjson/ffjson"

//go:generate ffjson $GOFILE

// PackageDTO is the data transfer object for an individual package.
type PackageDTO struct {
	Repo        string   `json:"repo"`
	Author      string   `json:"author"`
	Awesome     bool     `json:"awesome"`
	Versions    []string `json:"versions"`
	GodocURL    string   `json:"godocURL"`
	Description string   `json:"description"`
}

// PackageListItemDTO is the data transfer object for one package represented in
// a list of other package.
type PackageListItemDTO struct {
	Repo        string `json:"repo,omitempty"`
	Author      string `json:"author,omitempty"`
	Awesome     bool   `json:"awesome,omitempty"`
	Description string `json:"description,omitempty"`
}

// NewPackageListItemDTO creates a new PackageListItemDTO.
func NewPackageListItemDTO(packageModel *PackageModel) *PackageListItemDTO {
	var (
		awesome                   bool
		repo, author, description string
	)

	// TODO(skeswa): think on whether there is a cleaner way to handle unexpected
	// nils.

	if packageModel.Repo != nil {
		repo = *packageModel.Repo
	}

	if packageModel.Author != nil {
		author = *packageModel.Author
	}

	if packageModel.Description != nil {
		description = *packageModel.Description
	}

	if packageModel.AwesomeGo != nil {
		awesome = *packageModel.AwesomeGo
	}

	return &PackageListItemDTO{
		Repo:        repo,
		Author:      author,
		Awesome:     awesome,
		Description: description,
	}
}

// PackageListDTO is the data transfer object for a list of PackageListItemDTOs.
type PackageListDTO []*PackageListItemDTO

// NewPackageListDTOFromPackageModelList creates a new PackageListDTO from a
// list of package models.
func NewPackageListDTOFromPackageModelList(packageModels []*PackageModel) PackageListDTO {
	packageListItemDTOs := make(PackageListDTO, 0)

	for _, packageModel := range packageModels {
		packageListItemDTOs = append(packageListItemDTOs, NewPackageListItemDTO(packageModel))
	}

	return packageListItemDTOs
}

// MarshalJSON returns the JSON encoding of the PackageListDTO.
func (dto PackageListDTO) MarshalJSON() ([]byte, error) {
	return ffjson.Marshal(dto)
}
