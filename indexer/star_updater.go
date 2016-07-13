package main

import (
	"log"
	"os"
	"sync"
	"time"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
)

func ReIndexPackageGitHubStats(session *gocql.Session) {
	log.Println("Starting ReIndexing PackageModel GitHub Data")
	packageModels, err := common.ScanAllPackageModels(session)
	log.Println(err)
	log.Printf("%d packages found", len(packageModels))

	totalNumPackages := len(packageModels)
	log.Printf("Total num packages found = %d ", totalNumPackages)

	log.Println("Initializing GitHub Component")
	gitHubRequestService := common.NewGitHubRequestService()

	var wg sync.WaitGroup
	nbConcurrentInserts := 20
	packageChan := make(chan common.GitHubPackageModelDTO, 20)

	log.Printf("Spinning up %d consumers", nbConcurrentInserts)
	for i := 0; i < nbConcurrentInserts; i++ {
		wg.Add(1)
		go func() {
			for gitHubPackageModelDTO := range packageChan {
				// TODO update packageModel with new github data and index time
				packageStarCount := int(gitHubPackageModelDTO.ResponseBody["stargazers_count"].(float64))
				indexTime := time.Now()
				log.Printf("StarCount %d \n", packageStarCount)
				log.Printf("New index time %s \n", indexTime)
				// TODO save data to DB
			}
			wg.Done()
		}()
	}

	log.Printf("Preparing to fetch stars for %d repos", totalNumPackages)
	for count, packageModel := range packageModels {
		log.Printf("PROCESSING PACKAGE %s/%s #%d \n", *packageModel.Author, *packageModel.Repo, count)
		packageModelGitHubData, err := gitHubRequestService.FetchGitHubDataForPackageModel(*packageModel)
		if err != nil {
			log.Printf("PANIC %v \n", err)
			os.Exit(3)
		}
		packageChan <- common.GitHubPackageModelDTO{Package: *packageModel, ResponseBody: packageModelGitHubData}
		time.Sleep(100 * time.Millisecond)
	}

	close(packageChan)
	wg.Wait()
	log.Println("Finished testing star fetching")
}
