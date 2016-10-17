package awesome

import (
	"log"

	"github.com/gophr-pm/gophr/lib"
)

// Index is responsible for finding all go awesome packages
// and persisting them in `awsome_packages table` for later look up.
func Index() {
	_, session := common.Init()
	defer session.Close()

	log.Println("Fetching awesome go list.")
	awesomePackages, err := fetchAwesomeGoList()
	if err != nil {
		log.Fatalln(err)
	}

	log.Println("Persisting awesome go list.")
	if err = persistAwesomePackages(session, awesomePackages); err != nil {
		log.Fatalln(err)
	}
}
