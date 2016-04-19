package main

import (
	"log"
	"sync"

	"github.com/gocql/gocql"
	"github.com/skeswa/gophr/common"
)

func main() {
	log.Println("Establishing connection to DB")
	cluster := gocql.NewCluster("gophr.dev")
	cluster.ProtoVersion = 4
	cluster.Keyspace = "gophr"
	cluster.Consistency = gocql.One
	session, err := cluster.CreateSession()
	if err != nil {
		log.Println("Connection failed to establish successfully")
		log.Fatalln(err)
	}
	defer session.Close()
	log.Println("Connection established successfully")

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
	for _, packageModel := range packageModels {
		err = common.InsertPackage(session, packageModel)
		go func(packageModel *common.PackageModel) {
			wg.Add(1)
			if err != nil {
				json, errz := packageModel.MarshalJSON()
				if errz == nil {
					log.Fatalln(string(json[:]), err)
				}
			}
			wg.Done()
		}(packageModel)
	}
	wg.Wait()

	log.Println("Finished inserting packages into database")
}
