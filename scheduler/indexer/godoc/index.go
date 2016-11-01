package godoc

import (
	"log"

	"github.com/PuerkitoBio/goquery"
	"github.com/gophr-pm/gophr/lib"
)

// Index is responsible for looking up every package in godoc.org/Index
// and persisting the packages.
func Index() {
	_, session := lib.Init()
	defer session.Close()

	log.Println("Fetching godoc package metadata.")
	packageMetadata, err := FetchPackageMetadata(
		FetchPackageMetadataArgs{
			ParseHTML: goquery.NewDocument,
		})
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Batch uploading godoc packages.")
	persistGodocPackages(session, packageMetadata)
}
