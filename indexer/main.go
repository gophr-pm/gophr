package main

import (
	"bytes"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	"strings"
)

func main() {
	log.Println("Begining download")
	doc, err := goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished download")
	var buffer bytes.Buffer

	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()

		var (
			githubURL   string
			description string
			//goDocURL    string
		)

		children.Each(func(i int, s2 *goquery.Selection) {
			childURL, childURLexists := s.Find("a").Attr("href")
			childDescription := s.Text()

			if childURLexists == true {
				githubURL = childURL
			}

			if len(childDescription) > 0 {
				description = childDescription
			}
		})

		if strings.Contains(githubURL, "github.com") {
			// TODO build go doc url
			//goDocUrl =

			// TODO check for descriptions with just github url and remove
			if githubURL == "/"+githubURL {
				description = ""
			}

			buffer.WriteString("{\"url\": \"" + githubURL + "\", \"description\": \"" + description + "\"},\n")
		}
	})

	// TODO temp write JSON files
	jsonData := []byte(buffer.String())
	_ = ioutil.WriteFile("./tmp.json", jsonData, 0644)
	log.Println("Finished")
}
