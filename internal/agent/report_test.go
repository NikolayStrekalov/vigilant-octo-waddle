package agent

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/semaphore"
)

func Test_sendStat(t *testing.T) {
	type args struct {
		kind  StatKind
		name  StatName
		value string
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
	}{
		{
			name: "Test Post path 1",
			args: args{
				kind:  counterKind,
				name:  statPollCount,
				value: "314",
			},
			wantPath: "/counter/PollCount/314",
		},
		{
			name: "Test Post path 2",
			args: args{
				kind:  gaugeKind,
				name:  statBuckHashSys,
				value: "3.1415",
			},
			wantPath: "/gauge/BuckHashSys/3.1415",
		},
	}
	var requestPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
		requestPath = r.URL.Path
	}))
	defer ts.Close()
	ReportBaseURL = ts.URL
	for _, tt := range tests {
		sendStat(tt.args.kind, tt.args.name, tt.args.value)
		assert.Equal(t, requestPath, tt.wantPath)
	}
}

func Ptr[T any](v T) *T {
	return &v
}

func Test_sendStatJSON(t *testing.T) {
	type args struct {
		m *models.Metrics
	}
	tests := []struct {
		name    string
		args    args
		wantStr string
	}{
		{
			name:    "Empty model",
			args:    args{m: &models.Metrics{}},
			wantStr: `{"id":"","type":""}`,
		},
		{
			name:    "Counter model",
			args:    args{m: &models.Metrics{ID: "2", MType: "counter", Delta: Ptr(int64(5))}},
			wantStr: `{"id":"2","type":"counter","delta":5}`,
		},
		{
			name:    "Gauge model",
			args:    args{m: &models.Metrics{ID: "2", MType: "gauge", Value: Ptr(float64(3.1415))}},
			wantStr: `{"id":"2","type":"gauge","value":3.1415}`,
		},
	}

	var (
		requestBody           []byte
		contentEncodingHeader string
		contentTypeHeader     string
	)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
		contentTypeHeader = r.Header.Get("Content-Type")
		contentEncodingHeader = r.Header.Get("Content-Encoding")
		requestBody, _ = io.ReadAll(r.Body)
		_ = r.Body.Close()
	}))
	defer ts.Close()
	ReportBaseURL = ts.URL
	RequestLimiter = semaphore.NewWeighted(1)
	for _, tt := range tests {
		_ = sendStatJSON(tt.args.m, ReportBaseURL)
		assert.Equal(t, "application/json", contentTypeHeader)
		assert.Equal(t, "gzip", contentEncodingHeader)

		gz, err := gzip.NewReader(bytes.NewReader(requestBody))
		assert.Nil(t, err)
		_ = gz.Close()
		data, err := io.ReadAll(gz)
		assert.Nil(t, err)
		assert.JSONEq(t, tt.wantStr, string(data))
	}
}
