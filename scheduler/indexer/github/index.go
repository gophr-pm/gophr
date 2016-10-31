package githubIndexer

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/gophr-pm/gophr/lib/github"
)

// Index is a service dedicated to fetching Github repo metadata
// for each package in our DB and updating metadata.
func Index(args IndexArgs) error {
	conf, client := args.Init()
	defer client.Close()

	log.Println("Reindexing github data for packages.")
	packageModels, err := args.PackageRetriever(client)

	numPackageModels := len(packageModels)
	log.Printf("%d packages found", numPackageModels)

	if err != nil || numPackageModels == 0 {
		return fmt.Errorf("Failed to retrieve any packages from the db: %v", err)
	}

	gitHubRequestService := args.NewGithubRequestService(
		github.RequestServiceArgs{
			ForIndexer: true,
			Conf:       conf,
			Queryable:  client,
		},
	)

	log.Println("LOL")
	log.Println(gitHubRequestService)

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
				err := args.PackageInserter(client, packageModel)
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
				args.PackageDeleter(client, packageModel)
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
		time.Sleep(args.RequestTimeBuffer)
	}

	close(packageChan)
	wg.Wait()
	log.Println("Finished testing star fetching")

	return nil
}
