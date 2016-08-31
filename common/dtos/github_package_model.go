package dtos

import "github.com/skeswa/gophr/common/models"

// GitHubPackageModelDTO TODO Optimize this
type GitHubPackageModelDTO struct {
	Package      models.PackageModel
	ResponseBody map[string]interface{}
}
