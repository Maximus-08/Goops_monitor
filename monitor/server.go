package main

import (
	"fmt"
	"net/http"
	"log"
)

func startServer(port string) {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	log.Printf("Starting test server on %s", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		log.Printf("Server failed: %v", err)
	}
}
