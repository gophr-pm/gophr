package godoc

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"testing"

	"github.com/PuerkitoBio/goquery"
	. "github.com/smartystreets/goconvey/convey"
)

func TestFetchGoDocList(t *testing.T) {
	Convey("Given a godoc url", t, func() {

		Convey("If it fails to fetch the html from godoc.org/index, it should return an error", func() {
			actualOutput, err := FetchPackageMetadata(FetchPackageMetadataArgs{
				ParseHTML: func(url string) (*goquery.Document, error) {
					return nil, errors.New("Failed to restrieve html from godoc")
				},
			})

			So(err, ShouldNotBeNil)
			So(actualOutput, ShouldBeNil)
		})

		var godocHTMLString = `<table><tr><td><a href='/github.com/gophr-pm/gophr'>github.com/gophr-pm/gophr</a></td><td>Best Package Manager in the world.</td></tr></table>`

		Convey("If it retrieves a package successfully, it should return that package as a list of PackageMetadata", func() {
			actualOutput, err := FetchPackageMetadata(FetchPackageMetadataArgs{
				ParseHTML: func(url string) (*goquery.Document, error) {
					buf := bytes.NewBufferString(godocHTMLString)
					doc, err := goquery.NewDocumentFromReader(buf)
					if err != nil {
						fmt.Println(err)
					}
					return doc, err
				},
			})

			So(err, ShouldBeNil)
			log.Println(actualOutput)
			expectedOutput := []PackageMetadata{
				PackageMetadata{
					githubURL:   "github.com/gophr-pm/gophr",
					description: "Best Package Manager in the world.",
					author:      "gophr-pm",
					repo:        "gophr",
				},
			}
			So(actualOutput, ShouldResemble, expectedOutput)
		})
	})
}
