package main

import (
	"net/http"
	"strconv"
	"strings"
)

func updateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
	}
	pathParts := strings.Split(req.URL.Path, "/")
	switch len(pathParts) {
	case 5:
		switch pathParts[2] {
		case "gauge":
			val, err := strconv.ParseFloat(pathParts[4], 64)
			if err != nil {
				http.Error(res, "Wrong float value!", http.StatusBadRequest)
			}
			storage.UpdateGauge(pathParts[3], val)
		case "counter":
			val, err := strconv.ParseInt(pathParts[4], 10, 64)
			if err != nil {
				http.Error(res, "Wrong integer value!", http.StatusBadRequest)
			}
			storage.IncrementCounter(pathParts[3], val)
		default:
			http.Error(res, "Wrong metric type!", http.StatusBadRequest)
		}
	default:
		http.Error(res, "Use format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>", http.StatusNotFound)
	}
	// storage.Log()
}
