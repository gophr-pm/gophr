package pkg

import (
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/pquerna/ffjson/ffjson"
)

// Summary summarizes a package. Usually used in plural situations.
type Summary struct {
	Repo             string
	Stars            int
	Author           string
	Awesome          bool
	Description      string
	DailyDownloads   int64
	WeeklyDownloads  int64
	MonthlyDownloads int64
	AllTimeDownloads int64
}

// Summaries is a list of summary structs.
type Summaries []Summary

// toDTO turns a summary into its most appropriate DTO.
func (s Summary) toDTO() dtos.PackageSummary {
	return dtos.PackageSummary{
		Repo:        s.Repo,
		Stars:       s.Stars,
		Author:      s.Author,
		Awesome:     s.Awesome,
		Description: s.Description,

		Downloads: dtos.PackageDownloads{
			Daily:   s.DailyDownloads,
			Weekly:  s.WeeklyDownloads,
			Monthly: s.MonthlyDownloads,
			AllTime: s.AllTimeDownloads,
		},
	}
}

// ToJSON turns a summary into JSON.
func (s Summary) ToJSON() ([]byte, error) {
	dto := s.toDTO()
	return dto.MarshalJSON()
}

// ToJSON turns summaries into JSON.
func (s Summaries) ToJSON() ([]byte, error) {
	dtos := make([]dtos.PackageSummary, len(s))

	for _, summary := range s {
		dtos = append(dtos, summary.toDTO())
	}

	return ffjson.Marshal(dtos)
}
