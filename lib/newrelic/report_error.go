package nr

import (
	"log"

	newrelic "github.com/newrelic/go-agent"
)

// ReportNewRelicError reports errors to newrelic.
func ReportNewRelicError(txn newrelic.Transaction, err error) {
	log.Println("Reporting error to newrelic: ", err)
	newRelicErr := txn.NoticeError(err)
	if newRelicErr != nil {
		log.Println("New relic error occured: ", err)
	}
}
