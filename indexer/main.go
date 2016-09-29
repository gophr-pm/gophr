package main

import (
	"log"
	"sync"

	"github.com/gophr-pm/gophr/common"
	"github.com/gophr-pm/gophr/common/models"
)

func main() {
	log.Println("Preparing to initialize DB connection")
	conf, session := common.Init()
	defer session.Close()

	log.Println("Fetching godoc metadata")
	godocMetadataList, err := fetchGodocMetadata()
	if err != nil {
		log.Println("Failed to fetch godoc metadata")
		log.Fatalln(err)
	}

	log.Println("Fetching awesome go list")
	awesomeGoIndex, err := fetchAwesomeGoList()
	if err != nil {
		log.Println("Failed to fetch awesome go list")
		log.Fatalln(err)
	}

	log.Println("Preparing to build package models")
	packageModels, err := buildPackageModels(godocMetadataList, awesomeGoIndex)
	if err != nil {
		log.Println("Failed to build package models")
		log.Fatalln(err)
	}

	log.Println("Preparing to insert packages into database")

	var wg sync.WaitGroup
	var insertPackageErrors []error

	nbConcurrentInserts := 20
	packageChan := make(chan *models.PackageModel, nbConcurrentInserts)
	for i := 0; i < nbConcurrentInserts; i++ {
		wg.Add(1)
		go func() {
			for packageModel := range packageChan {
				err := models.InsertPackage(session, packageModel)
				if err != nil {
					insertPackageErrors = append(insertPackageErrors, err)
				}
			}
			wg.Done()
		}()
	}

	for _, packageModel := range packageModels {
		packageChan <- packageModel
	}
	close(packageChan)
	wg.Wait()

	for _, insertErr := range insertPackageErrors {
		log.Println(insertErr)
	}

	log.Println("Finished inserting packages into database")
	ReIndexPackageGitHubStats(conf, session)
}
