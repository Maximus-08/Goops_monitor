package main

import (
	"fmt"
	"net/http"
)

func startServer(port string) {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			fmt.Fprintln(w, "goops-monitor test server")
			fmt.Fprintln(w, "Endpoints: /health")
			return
		}
		http.NotFound(w, r)
	})
	
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		fmt.Fprintln(w, "OK")
	})

	LogInfo("Starting test server", "port", port)
	if err := http.ListenAndServe(port, nil); err != nil {
		LogError("Server failed", "error", err.Error())
	}
}
