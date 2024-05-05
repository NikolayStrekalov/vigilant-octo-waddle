package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
)

func prepareRoutes(r *chi.Mux) {
	r.Get(getMetricPath, metricHandler)
	r.Post(updatePath, updateHandler)
}

var getMetricPath = "/value/{kind}/{name}"

func metricHandler(res http.ResponseWriter, req *http.Request) {
	kind := chi.URLParam(req, "kind")
	name := chi.URLParam(req, "name")
	switch kind {
	case "gauge":
		v, err := storage.getGauge(name)
		if err != nil {
			http.Error(res, "Metric not found!", http.StatusNotFound)
			return
		}
		if _, err := io.WriteString(res, strconv.FormatFloat(v, 'f', -1, 64)); err != nil {
			http.Error(res, "InternalServerError", http.StatusInternalServerError)
		}
	case "counter":
		v, err := storage.getCounter(name)
		if err != nil {
			http.Error(res, "Metric not found!", http.StatusNotFound)
			return
		}
		if _, err := io.WriteString(res, strconv.FormatInt(v, 10)); err != nil {
			http.Error(res, "InternalServerError", http.StatusInternalServerError)
		}
	default:
		http.Error(res, "Wrong metric type!", http.StatusNotFound)
		return
	}
}

var updatePath = "/update/{kind}/{name}/{value}"

func updateHandler(res http.ResponseWriter, req *http.Request) {
	kind := chi.URLParam(req, "kind")
	switch kind {
	case "gauge":
		val, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			http.Error(res, "Wrong float value!", http.StatusBadRequest)
			return
		}
		storage.UpdateGauge(chi.URLParam(req, "name"), val)
	case "counter":
		val, err := strconv.ParseInt(chi.URLParam(req, "value"), 10, 64)
		if err != nil {
			http.Error(res, "Wrong integer value!", http.StatusBadRequest)
			return
		}
		storage.IncrementCounter(chi.URLParam(req, "name"), val)
	default:
		http.Error(res, "Wrong metric type!", http.StatusBadRequest)
		return
	}
	storage.Log()
}
