package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"io/ioutil"
	"log"
	//"os"
	"net/http"
	"strings"
	"time"
)

type GoPackage struct {
	Description    string
	GitHubURL      string
	GoDocURL       string
	IndexTime      string
	HttpStatusCode int
}

func main() {
	log.Println("Begining download")
	doc, err := goquery.NewDocument("https://godoc.org/-/index")

	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished download")
	var goPackageList = make([]*GoPackage, 0)

	// For each tr element on the page
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		goPackage := new(GoPackage)

		// For each child in the tr element
		children.Each(func(i int, s2 *goquery.Selection) {
			childURL, childURLexists := s.Find("a").Attr("href")
			childDescription := s.Text()

			if childURLexists == true {
				goPackage.GitHubURL = childURL
			}

			if len(childDescription) > 0 {
				// TODO check if description isn't just the url, if so dont set it
				goPackage.Description = childDescription
			}
		})

		// Only continue if this goPackage contains a GitHub URL
		// TODO check to make sure github.com is the prefix
		if strings.Contains(goPackage.GitHubURL, "github.com") {
			// Build go doc url
			goPackage.GoDocURL = ("https://godoc.org/" + goPackage.GitHubURL)

			// Create Index Time
			t := time.Now()
			time := t.String()
			goPackage.IndexTime = time

			// Add goPackage to goPackageList
			goPackageList = append(goPackageList, goPackage)
		}
	})

	for _, goPackage := range goPackageList {
		resp, err := http.Get("http://www." + goPackage.GitHubURL)
		if err != nil {
			fmt.Println("ERROR")
		}
		defer resp.Body.Close()

		// TODO Check resp status code
		goPackage.HttpStatusCode = resp.StatusCode

		_, err = ioutil.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		// Call Router Functionality to Parse for array of versions

	}

	log.Println("Finished Building GoPackages")
	log.Println("Creating JSON dump")
	createJSONDump(goPackageList)
}

// Create a tmp JSON dump of all serialized goPackageData
func createJSONDump(goPackageList []*GoPackage) {
	var buffer bytes.Buffer

	for _, goPackage := range goPackageList {
		buffer.WriteString("{\"url\": \"" + goPackage.GitHubURL + "\", \"description\": \"" + goPackage.Description + "\", \"index_time\": \"" + goPackage.IndexTime + "\" },\n")
	}

	t := time.Now()
	time := t.String()
	jsonData := []byte(buffer.String())
	_ = ioutil.WriteFile("./"+time+"tmp.json", jsonData, 0644)
}

func printGoPackage(gp *GoPackage) {
	fmt.Println("GitHub URL = ", gp.GitHubURL)
	fmt.Println("Description = ", gp.Description)
	fmt.Println("Index time = ", gp.IndexTime)
	fmt.Println("GoDocURL = ", gp.GoDocURL, "\n")
}
