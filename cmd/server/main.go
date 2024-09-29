package main

import (
	"github.com/Xurliman/metrics-alert-system/cmd/server/handlers"
	"net/http"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /update/{type}/{name}/{value}", handlers.UpdateMetrics)
	return http.ListenAndServe(`:8080`, mux)
}
