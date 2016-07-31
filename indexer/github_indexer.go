package main

import (
	"log"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
	"github.com/skeswa/gophr/common/models"
)

var requestTimeBuffer time.Duration = 50 * time.Millisecond

func ReIndexPackageGitHubStats(session *gocql.Session) {
	log.Println("Reindexing packageModel github data")
	packageModels, err := models.ScanAllPackageModels(session)
	numPackageModels := len(packageModels)
	log.Printf("%d packages found", numPackageModels)

	if err != nil || numPackageModels == 0 {
		log.Println("Error retrieving querying package data")
		log.Fatalln(err)
	}

	log.Println("Initializing gitHub component")
	gitHubRequestService := common.NewGitHubRequestService()

	var wg sync.WaitGroup
	nbConcurrentInserts := 20
	packageChan := make(chan common.GitHubPackageModelDTO, 20)

	log.Printf("Spinning up %d consumers", nbConcurrentInserts)
	for i := 0; i < nbConcurrentInserts; i++ {
		wg.Add(1)
		go func() {
			for gitHubPackageModelDTO := range packageChan {
				packageStarCount := common.ParseStarCount(gitHubPackageModelDTO.ResponseBody)
				log.Printf("Star count %d \n", packageStarCount)
				indexTime := time.Now()
				log.Printf("New index time %s \n", indexTime)
				packageModel := gitHubPackageModelDTO.Package
				packageModel.IndexTime = &indexTime
				packageModel.Stars = &packageStarCount
				err := models.InsertPackage(session, &packageModel)
				if err != nil {
					log.Println("Could not insert packageModel, error occured")
					log.Println(err)
				}
			}
			wg.Done()
		}()
	}

	log.Printf("Preparing to fetch stars for %d repos", numPackageModels)
	for count, packageModel := range packageModels {
		log.Printf("Processing package %s/%s #%d \n", *packageModel.Author, *packageModel.Repo, count)
		packageModelGitHubData, err := gitHubRequestService.FetchGitHubDataForPackageModel(*packageModel)

		if packageModelGitHubData == nil && err == nil {
			log.Println(err)
			wg.Add(1)
			go func() {
				log.Println("Preparing to delete packageModel")
				models.DeletePackageModel(session, packageModel)
				wg.Done()
			}()
		} else if err != nil {
			log.Println("Package could not be successfully retrieved from Github. Error occured")
			log.Println(err)
		}

		packageChan <- common.GitHubPackageModelDTO{Package: *packageModel, ResponseBody: packageModelGitHubData}
		time.Sleep(requestTimeBuffer)
	}

	close(packageChan)
	wg.Wait()
	log.Println("Finished testing star fetching")
}
