package main

import (
	"log"
	"net/http"
)

type Mux struct {
	http.Handler
}

func handler(w http.ResponseWriter, r *http.Request) {
	log.Print(r.Host)
}

func main() {
	mux := http.NewServeMux()
	mux.HandleFunc("/", handler)

	log.Fatal(http.ListenAndServe(":8080", mux))
}
