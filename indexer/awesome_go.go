package main

import (
	"strings"

	"github.com/PuerkitoBio/goquery"
)

func fetchAwesomeGoList() (map[string]bool, error) {
	doc, err := goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		return nil, err
	}

	var awesomeGoMap = make(map[string]bool)

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		childURL, _ := s.Attr("href")
		childURL = strings.Trim(childURL, "/")

		awesomeGoMap[childURL] = true
	})

	return awesomeGoMap, nil
}
