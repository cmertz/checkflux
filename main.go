package main

import (
	"flag"
	"log"
	"net/http"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()

	results := make(chan result)

	for i := 0; i < 5000; i++ {
		go newRandomCheck(i, results).Perform()
	}

	http.HandleFunc("/results", wsHandler(newResultChan(results)))
	http.HandleFunc("/", newDashboard("/results"))

	log.Fatal(http.ListenAndServe(*addr, nil))
}
