package godoc

import "log"

// Index is responsible for looking up every package in godoc.org/Index
// and persisting the packages.
func Index() {
	log.Println("Fetching godoc metadata")
	metadata, err := fetchMetadata()
	if err != nil {
		log.Fatalln(err)
	}

	// TODO(Shikkic): batch upload godoc packages
	log.Println(metadata)
}
