package common

//go:generate ffjson $GOFILE

type PackageInstallDTO struct {
	Author  string      `json:"author"`
	Repo    string      `json:"repo"`
	Version *VersionDTO `json:"version"`
}

func NewPackageInstallDTO(author, repo, version string) *PackageInstallDTO {
	var versionDTO *VersionDTO
	if len(version) > 0 {
		versionDTO = NewVersionDTO(version)
	}

	return &PackageInstallDTO{
		Author:  author,
		Repo:    repo,
		Version: versionDTO,
	}
}
