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

func main() {
	log.Println("Started Download of Godoc/index")

	doc, err := goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Download of Godoc/index\n")
	log.Println("Started scraping all github packages from Godoc/index")

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
		gitHubUrlTokens := strings.Split(goPackage.GitHubURL, "/")

		if strings.Contains(goPackage.GitHubURL, "github.com") && len(gitHubUrlTokens) == 3 {
			// Build go doc url
			goPackage.GoDocURL = ("https://godoc.org/" + goPackage.GitHubURL)

			// Create Index Time
			t := time.Now()
			time := t.String()
			goPackage.IndexTime = time

			goPackageMap[goPackage.GitHubURL] = goPackage
		}
	})

	log.Println("Finished scraping ", len(goPackageMap), " github packages from Godoc/index\n")
	log.Println("Started Download of awesome-go/README.md")

	doc, err = goquery.NewDocument("https://godoc.org/-/index")
	if err != nil {
		log.Fatal(err)
	}

	log.Println("Finished Download of awesome-go/README.md\n")
	log.Println("Started scraping awesome-go/README.md")

	awesomeGoCount := 0
	doc.Find("a").Each(func(i int, s *goquery.Selection) {
		childURL, _ := s.Attr("href")
		fmt.Println(childURL)
		goPackage, goPackageExists := goPackageMap[childURL]

		if goPackageExists {
			goPackage.AwesomeGo = true
			goPackageMap[childURL] = goPackage
			awesomeGoCount++
		}
	})

	log.Println("Finished scraping ", awesomeGoCount, " awesome-go/README.md\n")
	log.Println("Started retrieving versions for go packages")

	resc, errc := make(chan *GoPackage), make(chan *GoPackage)
	var goErrorPackageMap = make(map[string]*GoPackage)

	for _, goPackage := range goPackageMap {
		go func(goPackage *GoPackage) {
			refs, err := common.FetchRefs(goPackage.GitHubURL)
			if err != nil {
				log.Println("ERROR", goPackage.GitHubURL, " failed to return.", err)
				errc <- goPackage
				return
			}

			var versions []string
			for _, version := range refs.Candidates {
				versions = append(versions, version.String())
			}
			goPackage.Versions = versions
			resc <- goPackage
		}(goPackage)
	}

	// Create Errored out Map
	length := len(goPackageMap)
	for i := 0; i < length; i++ {
		select {
		case goPackage := <-resc:
			goPackageMap[goPackage.GitHubURL] = goPackage
		case errPackage := <-errc:
			delete(goPackageMap, errPackage.GitHubURL)
			goErrorPackageMap[errPackage.GitHubURL] = errPackage
		}
	}

	log.Println("SUCCESS: ", len(goPackageMap), " GoPackages were successfully built")
	log.Println("Creating JSON dump of ", len(goPackageMap), " valid go packages")
	createJSONDump(goPackageMap, "valid-go-packages")
	log.Println("Finished creating JSON dump\n")

	log.Println("WARNING: ", len(goErrorPackageMap), " GoPackages were not found on github")
	log.Println("Creating JSON dump of ", len(goPackageMap), " err packages")
	createJSONDump(goErrorPackageMap, "invalid-go-packages")
	log.Println("Finished creating JSON dump")
}

// Create a tmp JSON dump of all serialized goPackageData
func createJSONDump(goPackageMap map[string]*GoPackage, fileName string) {
	var buffer bytes.Buffer

	for _, goPackage := range goPackageMap {
		buffer.WriteString("{\"url\": \"" + goPackage.GitHubURL + "\", \"description\": \"" + goPackage.Description + "\", \"index_time\": \"" + goPackage.IndexTime + "\" }, \"versions\": \"" + fmt.Sprintf("%v", goPackage.Versions) + "\" \n")
	}

	t := time.Now()
	time := strings.Replace(t.String(), " ", "", -1)
	jsonData := []byte(buffer.String())
	_ = ioutil.WriteFile("./"+time+"-"+fileName+".json", jsonData, 0644)
}

func printGoPackage(goPackage *GoPackage) {
	fmt.Printf("%v", goPackage)
}
