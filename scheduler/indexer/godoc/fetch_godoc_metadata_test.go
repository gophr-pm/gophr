package godoc

import (
	"bytes"
	"errors"
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

		var godocHTMLString = "<tr><td><a href='/9fans.net/go/acme'>9fans.net/go/acme</a></td><td>Package acme is a simple interface for interacting with acme windows.</td></tr>"

		Convey("If fetching package metadata ", func() {
			actualOutput, err := FetchPackageMetadata(FetchPackageMetadataArgs{
				ParseHTML: func(url string) (*goquery.Document, error) {
					/*
						htmlNode, err := html.Parse(strings.NewReader(godocHTMLString))
						if err != nil {
							log.Println(err)
						}
						goquery.NewDocumentFromNode(htmlNode)
						return , nil
					*/
					reader := bytes.NewBufferString(godocHTMLString)
					log.Println("READER")
					log.Println(reader)
					doc, err := goquery.NewDocumentFromReader(reader)
					if err != nil {
						log.Println(err)
					}
					return doc, err
				},
			})

			So(err, ShouldBeNil)
			log.Println(actualOutput)
			//So(actualOutput, ShouldNotBeNil)
		})
	})
}
