package main

import (
	"bytes"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/skeswa/gophr/common"
	"io/ioutil"
	"log"
	"strings"
	"time"
)

type GoPackage struct {
	Description    string
	GitHubURL      string
	GoDocURL       string
	IndexTime      string
	HttpStatusCode int
	AwesomeGo      bool
	Versions       []string
}

const (
	refsFetchURLTemplate = "https://%s.git/info/refs?service=git-upload-pack"
)

func main() {
	log.Println("Started Download of Godoc/index")

	doc, err := goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Download of Godoc/index")
	log.Println("Started scraping all github packages from Godoc/index")

	var goPackageList = make([]*GoPackage, 0)
	var goPackageMap = make(map[string]*GoPackage)

	// For each tr element on the page
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		goPackage := new(GoPackage)

		// For each child in the tr element
		children.Each(func(i int, s2 *goquery.Selection) {
			childURL, childURLexists := s.Find("a").Attr("href")
			childDescription := s.Text()

			if childURLexists == true {
				goPackage.GitHubURL = strings.Trim(childURL, "/")
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
			// Hash the goPackage.GitHubURL for lookup
			goPackageMap[goPackage.GitHubURL] = goPackage
		}
	})

	log.Println("Finished scraping all github packages from Godoc/index")
	log.Println("Started Download of awesome-go/README.md")

	doc, err = goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Download of awesome-go/README.md")
	log.Println("Started scraping awesome-go/README.md")

	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		childURL, _ := s.Attr("href")

		goPackage, goPackageExists := goPackageMap[childURL]

		if goPackageExists {
			goPackage.AwesomeGo = true
			goPackageMap[childURL] = goPackage
			printGoPackage(goPackage)
		}
	})

	log.Println("Finished scraping awesome-go/README.md")

	for _, goPackage := range goPackageMap {
		refs, err := common.FetchRefs(goPackage.GitHubURL)
		if err != nil {
			log.Println("ERROR", goPackage.GitHubURL, " failed to return.", err)
		}

		var versions []string
		for _, version := range refs.Candidates {
			versions = append(versions, version.String())
		}

		goPackage.Versions = versions
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
	fmt.Println("GoDocURL = ", gp.GoDocURL)
	fmt.Println("AwesomeGo = ", gp.AwesomeGo, "\n")
}
