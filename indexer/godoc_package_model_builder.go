package main

import (
	"log"
	"sync"
	"time"

	"github.com/skeswa/gophr/common"
)

func fetchPackageVersions(metadata godocMetadata) ([]string, error) {
	refs, err := common.FetchRefs(metadata.author, metadata.repo)
	if err != nil {
		return nil, err
	}

	var versions []string
	for _, version := range refs.Candidates {
		versions = append(versions, version.String())
	}

	return versions, nil
}

func buildPackageModels(godocMetadataList []godocMetadata, awesomeGoIndex map[string]bool) ([]*common.PackageModel, error) {
	nbConcurrentGet := 20
	metadataChan := make(chan godocMetadata, nbConcurrentGet)
	packageModelChan := make(chan *common.PackageModel, nbConcurrentGet)

	var wg sync.WaitGroup
	for i := 0; i < nbConcurrentGet; i++ {
		wg.Add(1)
		go func() {
			for metadata := range metadataChan {
				log.Printf("Fetching package versions for \"%s\"", metadata.githubURL)
				packageVersions, err := fetchPackageVersions(metadata)
				if err == nil {
					log.Printf("\"%s\" versions retrieved successfully", metadata.githubURL)
					_, isAwesome := awesomeGoIndex[metadata.githubURL]
					packageModel, err := common.NewPackageModelForInsert(
						metadata.author,
						true,
						metadata.repo,
						packageVersions,
						"godoc.org/"+metadata.githubURL,
						time.Now(),
						isAwesome,
						metadata.description,
					)
					if err == nil {
						packageModelChan <- packageModel
					} else {
						packageModelChan <- nil
					}
				} else {
					log.Printf("\"%s\" versions failed to retrieve successfully", metadata.githubURL)
					packageModelChan <- nil
				}
			}
			log.Println("wait group done")
			wg.Done()
		}()
	}

	var packageModels []*common.PackageModel
	go func() {
		wg.Add(1)
		for i := 0; i < len(godocMetadataList); i++ {
			log.Printf("Waiting for %d out of %d of godocMetadataList length \n", i+1, len(godocMetadataList))
			packageModel := <-packageModelChan
			log.Println("Recieved one package from packageModelChan")
			if packageModel != nil {
				packageModels = append(packageModels, packageModel)
				log.Println("Appending package to packageModels slice")
			}
		}
		log.Println("wait group done")
		wg.Done()
	}()

	for _, metadata := range godocMetadataList {
		metadataChan <- metadata
		log.Println("Queuing package into metadataChan")
	}

	close(metadataChan)

	wg.Wait()
	log.Println("Done waiting")

	close(packageModelChan)

	return packageModels, nil
}
