package main

import (
	"log"
	"net/http"
)

func main() {
	results := make(chan result)

	for i := 0; i < 5000; i++ {
		go newRandomCheck(i, results).Perform()
	}

	http.HandleFunc("/results", wsHandler(newResultChan(results)))
	http.HandleFunc("/", newDashboard("/results"))

	log.Fatal(http.ListenAndServe(":9090", nil))
}
