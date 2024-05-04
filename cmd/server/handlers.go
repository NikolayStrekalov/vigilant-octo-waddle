package main

import (
	"net/http"
	"strconv"
	"strings"
)

var requiredURLParts = 5
var kindIndex = 2
var valueIndex = 4
var metricNameIndex = 3

func updateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
		return
	}
	pathParts := strings.Split(req.URL.Path, "/")
	switch len(pathParts) {
	case requiredURLParts:
		switch pathParts[kindIndex] {
		case "gauge":
			val, err := strconv.ParseFloat(pathParts[valueIndex], 64)
			if err != nil {
				http.Error(res, "Wrong float value!", http.StatusBadRequest)
				return
			}
			storage.UpdateGauge(pathParts[metricNameIndex], val)
		case "counter":
			val, err := strconv.ParseInt(pathParts[valueIndex], 10, 64)
			if err != nil {
				http.Error(res, "Wrong integer value!", http.StatusBadRequest)
				return
			}
			storage.IncrementCounter(pathParts[metricNameIndex], val)
		default:
			http.Error(res, "Wrong metric type!", http.StatusBadRequest)
			return
		}
	default:
		http.Error(
			res,
			"Use format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>",
			http.StatusNotFound,
		)
		return
	}
	// storage.Log()
}
