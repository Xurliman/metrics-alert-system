package handlers

import (
	"net/http"
	"strconv"
)

type MemStorage struct {
}

func UpdateMetrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metricName := r.PathValue("name")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	metricVal := r.PathValue("value")
	_, err := strconv.ParseInt(metricVal, 10, 64)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}
	w.WriteHeader(http.StatusOK)
}
