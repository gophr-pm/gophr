package github

import (
	"strings"
	"unicode/utf8"

	"github.com/PuerkitoBio/goquery"
)

type godocMetadata struct {
	githubURL   string
	description string
	author      string
	repo        string
}

const (
	godocURL = "https://godoc.org/-/index"
)

// fetchGodocMetadata converts entries listed in godoc/Index
// into a godocMetadata struct.
func fetchGodocMetadata() ([]godocMetadata, error) {
	var (
		godocMetadataList []godocMetadata
		metadata          godocMetadata
	)

	doc, err := goquery.NewDocument(godocURL)
	if err != nil {
		return nil, err
	}

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		metadata = godocMetadata{}

		// For each child in the tr element
		children.Each(func(i int, s2 *goquery.Selection) {
			childURL, childURLexists := s.Find("a").Attr("href")
			childDescription := s.Text()

			if childURLexists == true {
				childURL = strings.Trim(childURL, "/")
				metadata.githubURL = childURL
			}

			if len(childDescription) > 0 {
				// TODO check if description isn't just the url, if so dont set it
				metadata.description = sanitizeUTF8String(strings.TrimPrefix(childDescription, childURL))
			}
		})

		githubURLTokens := strings.Split(metadata.githubURL, "/")

		if len(githubURLTokens) == 3 && githubURLTokens[0] == "github.com" {
			githubURLTokens := strings.Split(metadata.githubURL, "/")
			metadata.author = githubURLTokens[1]
			metadata.repo = githubURLTokens[2]

			godocMetadataList = append(godocMetadataList, metadata)
		}
	})

	return godocMetadataList, nil
}

func sanitizeUTF8String(str string) string {
	if !utf8.ValidString(str) {
		v := make([]rune, 0, len(str))
		for i, r := range str {
			if r == utf8.RuneError {
				_, size := utf8.DecodeRuneInString(str[i:])
				if size == 1 {
					continue
				}
			}
			v = append(v, r)
		}
		return string(v)
	}

	return str
}
