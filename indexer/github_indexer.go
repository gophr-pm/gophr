package main

import (
	"log"
	"sync"
	"time"

	"github.com/gophr-pm/gophr/lib/config"
	"github.com/gophr-pm/gophr/lib/db"
	"github.com/gophr-pm/gophr/lib/dtos"
	"github.com/gophr-pm/gophr/lib/github"
)

type packageRepoTuple struct {
	pkg      *models.PackageModel
	repoData dtos.GithubRepo
}

var requestTimeBuffer = 50 * time.Millisecond

// ReIndexPackageGitHubStats is a service dedicated to fetching Github repo metadata
// for each package in our DB and updating metadata
func ReIndexPackageGitHubStats(conf *config.Config, q db.Queryable) {
	log.Println("Reindexing packageModel github data")
	packageModels, err := models.ScanAllPackageModels(session)
	numPackageModels := len(packageModels)
	log.Printf("%d packages found", numPackageModels)

	if err != nil || numPackageModels == 0 {
		log.Println("Error retrieving querying package data")
		log.Fatalln(err)
	}

	log.Println("Initializing gitHub component")
	gitHubRequestService := github.NewRequestService(
		github.RequestServiceArgs{
			ForIndexer: true,
			Conf:       conf,
			Queryable:  q,
		},
	)

	var wg sync.WaitGroup
	nbConcurrentInserts := 20
	packageChan := make(chan packageRepoTuple, 20)

	log.Printf("Spinning up %d consumers", nbConcurrentInserts)
	for i := 0; i < nbConcurrentInserts; i++ {
		wg.Add(1)
		go func() {
			for tuple := range packageChan {
				indexTime := time.Now()
				packageModel := tuple.pkg
				packageModel.Description = &tuple.repoData.Description
				packageModel.IndexTime = &indexTime
				packageModel.Stars = &tuple.repoData.Stars
				err := models.InsertPackage(session, packageModel)
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
		packageModelGitHubData, err := gitHubRequestService.FetchGitHubDataForPackageModel(*packageModel.Author, *packageModel.Repo)

		if err == nil {
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

		packageChan <- packageRepoTuple{
			pkg:      packageModel,
			repoData: packageModelGitHubData,
		}
		time.Sleep(requestTimeBuffer)
	}

	close(packageChan)
	wg.Wait()
	log.Println("Finished testing star fetching")
}
