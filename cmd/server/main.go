package main

import (
	"net/http"
	"strings"
)

func main() {
	if err := run(); err != nil {
		panic(err)
	}
}

func run() error {
	http.HandleFunc("/update/counter/", webhook)
	http.HandleFunc("/update/gauge/", webhook)
	return http.ListenAndServe(`:8080`, nil)
}

func webhook(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	path := r.URL.Path
	nameAndVal := strings.ReplaceAll(path, "/update/counter/", "")
	if nameAndVal == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	w.WriteHeader(http.StatusOK)
	return
}
