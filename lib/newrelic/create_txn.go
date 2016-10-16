package nr

import (
	"log"
	"net/http"

	newrelic "github.com/newrelic/go-agent"
)

// CreateNewRelicTxn starts a transaction that will log activity
// to new relic.
func CreateNewRelicTxn(
	newRelicApp newrelic.Application,
	w *http.ResponseWriter,
	r *http.Request,
) newrelic.Transaction {
	var txn newrelic.Transaction
	log.Printf("Logging request for %s \n", r.URL.String())
	// Create a new relic transaction.
	txn = newRelicApp.StartTransaction(r.URL.String(), *w, r)

	return txn
}
