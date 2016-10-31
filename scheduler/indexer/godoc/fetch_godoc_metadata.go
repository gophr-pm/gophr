package godoc

/*
import (
	"log"
	"strings"
	"unicode/utf8"

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
	for _, node := range doc.Nodes {
		log.Println(node)
	}

	log.Println(doc.Html())
	sec := doc.Find("tr")
	log.Println("SEC", sec)

	// Traverse the godoc.org/index html and find every instance of <tr>.
	// This is because Godoc organizes their packages in tables.
	doc.Find("tr").Each(func(i int, s *goquery.Selection) {
		log.Println(s)
		children := s.Children()
		godocMetadata = PackageMetadata{}

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
*/
