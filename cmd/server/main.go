package main

import (
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

type MemStorage struct {
	gauge   map[string]float64
	counter map[string]int64
}

// In-memory storage for received data:
// 'gauge' field stores the received value,
// 'counter' field increments the received value.
var storage = MemStorage{
	gauge:   make(map[string]float64),
	counter: make(map[string]int64),
}

func updateHandler(res http.ResponseWriter, req *http.Request) {
	if req.Method != http.MethodPost {
		http.Error(res, "Only GET requests are allowed!", http.StatusBadRequest)
	}
	pathParts := strings.Split(req.URL.Path, "/")
	switch len(pathParts) {
	case 5:
		fmt.Println(pathParts[4])
		switch pathParts[2] {
		case "gauge":
			val, err := strconv.ParseFloat(pathParts[4], 64)
			if err != nil {
				http.Error(res, "Wrong float value!", http.StatusBadRequest)
			}
			storage.gauge[pathParts[3]] = val
		case "counter":
			val, err := strconv.ParseInt(pathParts[4], 10, 64)
			if err != nil {
				http.Error(res, "Wrong integer value!", http.StatusBadRequest)
			}
			storage.counter[pathParts[3]] += val
		default:
			http.Error(res, "Wrong metric type!", http.StatusBadRequest)
		}
	default:
		http.Error(res, "Use format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>", http.StatusNotFound)
	}
	fmt.Println(storage)
}

func main() {

	mux := http.NewServeMux()
	mux.HandleFunc(`/update/`, updateHandler)

	err := http.ListenAndServe(`:8080`, mux)
	if err != nil {
		panic(err)
	}
}
