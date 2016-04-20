package main

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

func fetchGodocMetadata() ([]godocMetadata, error) {
	doc, err := goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		return nil, err
	}

	var godocMetadataList []godocMetadata

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		godocMetadata := godocMetadata{}

		// For each child in the tr element
		children.Each(func(i int, s2 *goquery.Selection) {
			childURL, childURLexists := s.Find("a").Attr("href")
			childDescription := s.Text()

			if childURLexists == true {
				childURL = strings.Trim(childURL, "/")
				godocMetadata.githubURL = childURL
			}

			if len(childDescription) > 0 {
				// TODO check if description isn't just the url, if so dont set it
				godocMetadata.description = sanitizeUTF8String(childDescription)
			}
		})

		githubURLTokens := strings.Split(godocMetadata.githubURL, "/")

		if len(githubURLTokens) == 3 && githubURLTokens[0] == "github.com" {
			githubURLTokens := strings.Split(godocMetadata.githubURL, "/")
			godocMetadata.author = githubURLTokens[1]
			godocMetadata.repo = githubURLTokens[2]

			godocMetadataList = append(godocMetadataList, godocMetadata)
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
