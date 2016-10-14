package main

import (
	"fmt"
	"net/http"

	"github.com/robfig/cron"
)

func main() {
	c := cron.New()

	c.AddFunc("* * * * * *", func() { fmt.Println("Hey") })
	c.AddFunc("* * * * * *", Github_indexer)

	// Start the cron service
	c.Start()

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Requ	est) { fmt.Println("LOL") })
	http.ListenAndServe(":8080", nil)
}


func scanPackageModels(query *gocql.Query) ([]*PackageModel, error) {
	var (
		err          error
		scanError    error
		closeError   error
		packageModel *PackageModel

		repo        string
		author      string
		description string

		iter          = query.Iter()
		packageModels = make([]*PackageModel, 0)
	)

	for iter.Scan(&repo, &author, &description) {
		packageModel, err = NewPackageModelFromBulkSelect(author, repo, description)
		if err != nil {
			scanError = err
			break
		} else {
			packageModels = append(packageModels, packageModel)
		}
	}

	if err = iter.Close(); err != nil {
		closeError = err
	}

	if scanError != nil || closeError != nil {
		return nil, errors.NewQueryScanError(scanError, closeError)
	}

	return packageModels, nil
}
