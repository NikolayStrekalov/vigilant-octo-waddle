package server

import (
	"io"
	"net/http"
	"strconv"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/pgstorage"
	"github.com/mailru/easyjson"

	"github.com/go-chi/chi/v5"
)

const (
	indexPath                  = "/"
	getMetricPath              = "/value/{kind}/{name}"
	updateMetricPath           = "/update/{kind}/{name}/{value}"
	getMetricPathJSON          = "/value/"
	updateMetricPathJSON       = "/update/"
	pingPath                   = "/ping"
	bulkUpdatePath             = "/updates/"
	messageInternalServerError = "InternalServerError"
	gaugeKind                  = "gauge"
	counterKind                = "counter"
	metricNotFound             = "Metric not found!"
	wrongMetricType            = "Wrong metric type!"
	applicationJSONType        = "application/json"
)

func prepareRoutes(r *chi.Mux) {
	r.Get(indexPath, indexHandler)
	r.Get(getMetricPath, metricHandler)
	r.Post(updateMetricPath, updateMetricHandler)
	r.Post(getMetricPathJSON, metricJSONHandler)
	r.Post(updateMetricPathJSON, updateMetricJSONHandler)
	r.Get(pingPath, pingHandler)
	r.Post(bulkUpdatePath, bulkHandler)
}

func indexHandler(res http.ResponseWriter, req *http.Request) {
	counters := Storage.GetCounterList()
	gauges := Storage.GetGaugeList()
	html, err := renderIndexPage(counters, gauges)
	if err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", "text/html")
	res.WriteHeader(http.StatusOK)

	if _, err := html.WriteTo(res); err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
	}
}

func metricHandler(res http.ResponseWriter, req *http.Request) {
	kind := chi.URLParam(req, "kind")
	name := chi.URLParam(req, "name")
	switch kind {
	case gaugeKind:
		v, err := Storage.GetGauge(name)
		if err != nil {
			http.Error(res, metricNotFound, http.StatusNotFound)
			return
		}
		if _, err := io.WriteString(res, strconv.FormatFloat(v, 'f', -1, 64)); err != nil {
			http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		}
	case counterKind:
		v, err := Storage.GetCounter(name)
		if err != nil {
			http.Error(res, metricNotFound, http.StatusNotFound)
			return
		}
		if _, err := io.WriteString(res, strconv.FormatInt(v, 10)); err != nil {
			http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		}
	default:
		http.Error(res, wrongMetricType, http.StatusNotFound)
		return
	}
}

func updateMetricHandler(res http.ResponseWriter, req *http.Request) {
	kind := chi.URLParam(req, "kind")
	switch kind {
	case gaugeKind:
		val, err := strconv.ParseFloat(chi.URLParam(req, "value"), 64)
		if err != nil {
			http.Error(res, "Wrong float value!", http.StatusBadRequest)
			return
		}
		Storage.UpdateGauge(chi.URLParam(req, "name"), val)
	case counterKind:
		val, err := strconv.ParseInt(chi.URLParam(req, "value"), 10, 64)
		if err != nil {
			http.Error(res, "Wrong integer value!", http.StatusBadRequest)
			return
		}
		Storage.IncrementCounter(chi.URLParam(req, "name"), val)
	default:
		http.Error(res, wrongMetricType, http.StatusBadRequest)
		return
	}
}

func metricJSONHandler(res http.ResponseWriter, req *http.Request) {
	if val, ok := req.Header["Content-Type"]; !ok || val[0] != applicationJSONType {
		http.Error(res, "Wrong Content-Type, use application/json!", http.StatusBadRequest)
		return
	}
	m := models.Metrics{}
	data, err := io.ReadAll(req.Body)
	if err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		return
	}
	if err := easyjson.Unmarshal(data, &m); err != nil {
		http.Error(res, "Wrong json provided.", http.StatusBadRequest)
		return
	}
	switch m.MType {
	case counterKind:
		v, err := Storage.GetCounter(m.ID)
		if err != nil {
			http.Error(res, metricNotFound, http.StatusNotFound)
			return
		}
		m.Delta = &v
	case gaugeKind:
		v, err := Storage.GetGauge(m.ID)
		if err != nil {
			http.Error(res, metricNotFound, http.StatusNotFound)
			return
		}
		m.Value = &v
	default:
		http.Error(res, wrongMetricType, http.StatusBadRequest)
		return
	}
	rawBytes, err := easyjson.Marshal(&m)
	if err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", applicationJSONType)
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(rawBytes); err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
	}
}

func updateMetricJSONHandler(res http.ResponseWriter, req *http.Request) {
	if val, ok := req.Header["Content-Type"]; !ok || val[0] != applicationJSONType {
		http.Error(res, "Wrong Content-Type, use application/json!", http.StatusBadRequest)
		return
	}
	m := models.Metrics{}
	data, err := io.ReadAll(req.Body)
	defer func() { _ = req.Body.Close() }()
	if err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		return
	}
	if err := easyjson.Unmarshal(data, &m); err != nil {
		http.Error(res, "Wrong json provided.", http.StatusBadRequest)
		return
	}
	switch m.MType {
	case counterKind:
		if m.Delta == nil {
			http.Error(res, "Provide delta field for increment!", http.StatusBadRequest)
			return
		}
		Storage.IncrementCounter(m.ID, *m.Delta)
		v, err := Storage.GetCounter(m.ID)
		if err != nil {
			http.Error(res, metricNotFound, http.StatusNotFound)
			return
		}
		*m.Delta = v
	case gaugeKind:
		if m.Value == nil {
			http.Error(res, "Provide value field for update!", http.StatusBadRequest)
			return
		}
		Storage.UpdateGauge(m.ID, *m.Value)
		v, err := Storage.GetGauge(m.ID)
		if err != nil {
			http.Error(res, metricNotFound, http.StatusNotFound)
			return
		}
		*m.Value = v
	default:
		http.Error(res, wrongMetricType, http.StatusBadRequest)
		return
	}
	rawBytes, err := easyjson.Marshal(&m)
	if err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		return
	}
	res.Header().Set("Content-Type", applicationJSONType)
	res.WriteHeader(http.StatusOK)
	if _, err := res.Write(rawBytes); err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
	}
}

func pingHandler(res http.ResponseWriter, req *http.Request) {
	if pgstorage.Ping(ServerConfig.DatabaseDSN) {
		res.WriteHeader(http.StatusOK)
		return
	}
	res.WriteHeader(http.StatusInternalServerError)
}

func bulkHandler(res http.ResponseWriter, req *http.Request) {
	if val, ok := req.Header["Content-Type"]; !ok || val[0] != applicationJSONType {
		http.Error(res, "Wrong Content-Type, use application/json!", http.StatusBadRequest)
		return
	}
	metrics := models.MetricsSlice{}
	data, err := io.ReadAll(req.Body)
	defer func() { _ = req.Body.Close() }()
	if err != nil {
		http.Error(res, messageInternalServerError, http.StatusInternalServerError)
		return
	}

	if err := easyjson.Unmarshal(data, &metrics); err != nil {
		http.Error(res, "Wrong json provided.", http.StatusBadRequest)
		return
	}
	Storage.BulkUpdate(metrics)
	res.WriteHeader(http.StatusOK)
}
