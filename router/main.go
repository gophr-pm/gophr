package main

import (
	"fmt"
	"net/http"
	"os"
	"strconv"
	"time"
)

func annoy() {
	for {
		fmt.Println("Still here fam")
		time.Sleep(1 * time.Second)
	}
}

func main() {
	http.HandleFunc("/status", healthCheckHandler)
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Sup, I'm routerr. I love %s!", r.URL.Path[1:])
	})
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

	go annoy()

	http.ListenAndServe(fmt.Sprintf(":%d", port), nil)
}
