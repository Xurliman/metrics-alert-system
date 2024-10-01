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
	metricType := r.PathValue("type")
	if metricType != "gauge" && metricType != "counter" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	metricName := r.PathValue("name")
	if metricName == "" {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	metricVal := r.PathValue("value")
	if metricType == "counter" {
		_, err := strconv.ParseInt(metricVal, 10, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	if metricType == "gauge" {
		_, err := strconv.ParseFloat(metricVal, 64)
		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
	}
	w.WriteHeader(http.StatusOK)
}
