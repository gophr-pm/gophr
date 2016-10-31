package pkg

import (
	"time"

	"github.com/gophr-pm/gophr/lib/dtos"
)

// Details holds package details. Usually used in singular situations.
type Details struct {
	Summary

	TrendScore              float32
	DateDiscovered          time.Time
	DateLastIndexed         time.Time
	AllTimeVersionDownloads map[string]int64
}

// toDTO turns a summary into its most appropriate DTO.
func (d Details) toDTO() dtos.PackageDetails {
	versions := make([]dtos.PackageVersion, len(d.AllTimeVersionDownloads))

	for name, downloads := range d.AllTimeVersionDownloads {
		versions = append(versions, dtos.PackageVersion{
			Name:             name,
			AllTimeDownloads: downloads,
		})
	}

	return dtos.PackageDetails{
		Repo:            d.Repo,
		Stars:           d.Stars,
		Author:          d.Author,
		Awesome:         d.Awesome,
		Versions:        versions,
		TrendScore:      d.TrendScore,
		Description:     d.Description,
		DateDiscovered:  d.DateDiscovered,
		DateLastIndexed: d.DateLastIndexed,

		Downloads: dtos.PackageDownloads{
			Daily:   d.DailyDownloads,
			Weekly:  d.WeeklyDownloads,
			Monthly: d.MonthlyDownloads,
			AllTime: d.AllTimeDownloads,
		},
	}
}

// ToJSON turns a details struct into JSON.
func (d Details) ToJSON() ([]byte, error) {
	dto := d.toDTO()
	return dto.MarshalJSON()
}
