package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gophr-pm/gophr/lib"
	"github.com/gophr-pm/gophr/lib/db/model/package/download"
	"github.com/gorilla/mux"
)

func main() {
	// Initialize the API.
	config, client := lib.Init()

	// Ensure that the client is closed eventually.
	defer client.Close()

	// Register all of the routes.
	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/blob/{%s}/{%s}/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo,
		urlVarSHA,
		urlVarPath),
		BlobHandler()).Methods("GET")
	r.HandleFunc(
		"/packages/new",
		GetNewPackagesHandler(client)).Methods("GET")
	r.HandleFunc(
		"/packages/search",
		SearchPackagesHandler(client)).Methods("GET")
	r.HandleFunc(
		"/packages/trending",
		GetTrendingPackagesHandler(client)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/packages/top/{%s}/{%s}",
		urlVarLimit,
		urlVarTimeSplit),
		GetTopPackagesHandler(client)).Methods("GET")
	r.HandleFunc(fmt.Sprintf(
		"/packages/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo),
		GetPackageHandler(client)).Methods("GET")

	// TODO(skeswa): delete this please.
	r.HandleFunc(fmt.Sprintf(
		"/downloads/{%s}/{%s}",
		urlVarAuthor,
		urlVarRepo),
		func(w http.ResponseWriter, r *http.Request) {
			vars := mux.Vars(r)
			author := vars[urlVarAuthor]
			repo := vars[urlVarRepo]

			fmt.Printf("got /downloads/%s/%s\n", author, repo)
			splits, err := download.GetSplits(client, author, repo)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("could not get downloads: " + err.Error()))
				return
			}

			splitsMap := map[string]int{
				"daily":   splits.Daily,
				"weekly":  splits.Weekly,
				"monthly": splits.Monthly,
				"allTime": splits.AllTime,
			}
			splitsJSON, err := json.Marshal(splitsMap)
			if err != nil {
				w.WriteHeader(500)
				w.Write([]byte("could not marshal json: " + err.Error()))
				return
			}

			respondWithJSON(w, splitsJSON)
		}).Methods("GET")

	// Start serving.
	log.Printf("Servicing HTTP requests on port %d.\n", config.Port)
	http.ListenAndServe(fmt.Sprintf(":%d", config.Port), r)
}
