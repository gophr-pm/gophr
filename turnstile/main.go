package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"

	"github.com/gocql/gocql"
	"github.com/gorilla/mux"
)

func main() {
	cluster := gocql.NewCluster("gophr-db")
	cluster.ProtoVersion = 4
	cluster.Keyspace = "gophr"
	cluster.Consistency = gocql.One
	session, err := cluster.CreateSession()
	defer session.Close()

	if err != nil {
		log.Fatalln("Failed to connect to the database:", err)
	}

	r := mux.NewRouter()
	r.HandleFunc("/status", StatusHandler()).Methods("GET")

	portStr := os.Getenv("PORT")
	var port int
	if len(portStr) == 0 {
		fmt.Println("Port left unspecified; setting port to 3000.")
		port = 3000
	} else if portNum, err := strconv.Atoi(portStr); err == nil {
		fmt.Printf("Port was specified as %d.\n", portNum)
		port = portNum
	} else {
		fmt.Println("Port was invalid; setting port to 3000.")
		port = 3000
	}

	http.ListenAndServe(fmt.Sprintf(":%d", port), r)
}
