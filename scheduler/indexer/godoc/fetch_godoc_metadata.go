package godoc

import (
	"log"
	"strings"

	"github.com/PuerkitoBio/goquery"
)

const (
	godocURL      = "https://godoc.org/-/index"
	gitHubBaseURL = "github.com"
)

// FetchPackageMetadata converts entries listed in godoc.org/index
// into a package metadata struct.
func FetchPackageMetadata(args FetchPackageMetadataArgs) ([]PackageMetadata, error) {
	var (
		metadataList  []PackageMetadata
		godocMetadata PackageMetadata
	)

	doc, err := args.ParseHTML(godocURL)
	if err != nil {
		return nil, err
	}

	// Traverse the godoc.org/index html and find every instance of <tr>.
	// This is because Godoc organizes their packages in tables.
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		children := s.Children()
		log.Println(children)
		godocMetadata = PackageMetadata{}

		// For each child in the tr element
		children.Each(func(i int, s2 *goquery.Selection) {
			childURL, childURLexists := s.Find("a").Attr("href")
			log.Println(childURL)
			childDescription := s.Text()

			if childURLexists == true {
				childURL = strings.Trim(childURL, "/")
				godocMetadata.githubURL = childURL
			}

			if len(childDescription) > 0 {
				// TODO check if description isn't just the url, if so dont set it
				godocMetadata.description = sanitizeUTF8String(strings.TrimPrefix(childDescription, childURL))
			}
		})

		githubURLTokens := strings.Split(godocMetadata.githubURL, "/")

		if len(githubURLTokens) == 3 && githubURLTokens[0] == gitHubBaseURL {
			githubURLTokens := strings.Split(godocMetadata.githubURL, "/")
			godocMetadata.author = githubURLTokens[1]
			godocMetadata.repo = githubURLTokens[2]

			metadataList = append(metadataList, godocMetadata)
		}
	})

	return metadataList, nil
}
