package dtos

import "github.com/gophr-pm/gophr/lib/model"

// GitHubPackageModelDTO TODO Optimize this
type GitHubPackageModelDTO struct {
	Package      models.PackageModel
	ResponseBody map[string]interface{}
}
