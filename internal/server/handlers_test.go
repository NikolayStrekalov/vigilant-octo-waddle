package server

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/memstorage"
	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_updateMetricHandler(t *testing.T) {
	type args struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test gauge update",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/gauge/RandomValue/3.1415", http.NoBody),
			},
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "Test counter update",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/counter/PollCount/31415", http.NoBody),
			},
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "Test fail slash end",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/gauge/RandomValue/3.1415/", http.NoBody),
			},
			want: want{
				code:        404,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test fail unknown kind",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/kind/RandomValue/3.1415", http.NoBody),
			},
			want: want{
				code:        400,
				response:    "Wrong metric type!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test fail wrong int",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/counter/PollCount/3.1415", http.NoBody),
			},
			want: want{
				code:        400,
				response:    "Wrong integer value!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test fail wrong float",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/gauge/RandomValue/e3.1415", http.NoBody),
			},
			want: want{
				code:        400,
				response:    "Wrong float value!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test fail no kind",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodPost, "http://localhost:8080/update/RandomValue/e3.1415", http.NoBody),
			},
			want: want{
				code:        404,
				response:    "404 page not found\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			prepareRoutes(r)
			r.ServeHTTP(tt.args.res, tt.args.req)
			res := tt.args.res.Result()

			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer func() {
				_ = res.Body.Close()
			}()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_metricHandler(t *testing.T) {
	type args struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test get gauge",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/gauge/RandomValue", http.NoBody),
			},
			want: want{
				code:        200,
				response:    "0.31",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test get counter",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/counter/PollCount", http.NoBody),
			},
			want: want{
				code:        200,
				response:    "-62",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test unknown kind",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/bool/PollCount", http.NoBody),
			},
			want: want{
				code:        404,
				response:    "Wrong metric type!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test wrong kind",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/gauge/PollCount", http.NoBody),
			},
			want: want{
				code:        404,
				response:    "Metric not found!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test gauge not found",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/gauge/something", http.NoBody),
			},
			want: want{
				code:        404,
				response:    "Metric not found!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test counter not found",
			args: args{
				res: httptest.NewRecorder(),
				req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/value/counter/something", http.NoBody),
			},
			want: want{
				code:        404,
				response:    "Metric not found!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		storage, _, _ := memstorage.NewMemStorage("", false, 300)
		storage.Gauge = map[string]float64{
			"RandomValue": 0.31,
			"qwer":        3.1415,
		}
		storage.Counter = map[string]int64{
			"PollCount": -62,
			"ewq":       9321,
		}
		Storage = storage
		t.Run(tt.name, func(t *testing.T) {
			r := chi.NewRouter()
			prepareRoutes(r)
			r.ServeHTTP(tt.args.res, tt.args.req)
			res := tt.args.res.Result()

			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)

			// получаем и проверяем тело запроса
			defer func() {
				_ = res.Body.Close()
			}()
			resBody, err := io.ReadAll(res.Body)

			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}

func Test_indexHandler(t *testing.T) {
	args := struct {
		res *httptest.ResponseRecorder
		req *http.Request
	}{
		res: httptest.NewRecorder(),
		req: httptest.NewRequest(http.MethodGet, "http://localhost:8080/", http.NoBody),
	}
	want := struct {
		code        int
		contentType string
	}{
		code:        200,
		contentType: "text/html",
	}

	r := chi.NewRouter()
	prepareRoutes(r)
	r.ServeHTTP(args.res, args.req)
	res := args.res.Result()

	// проверяем код ответа
	assert.Equal(t, want.code, res.StatusCode)

	defer func() {
		_ = res.Body.Close()
	}()
	_, err := io.ReadAll(res.Body)
	require.NoError(t, err)
	// проверяем Content-Type
	assert.Equal(t, want.contentType, res.Header.Get("Content-Type"))
}

func Test_bulkHandler(t *testing.T) {
	validBody := bytes.NewReader([]byte(`[{
		"id": "qwe",
		"type": "counter",
		"delta": 3
	}]`))
	invalidBody := bytes.NewReader([]byte(`[{
		"id": "qwe",
		"type": "counter",
		"delta":
	}]`))
	type args struct {
		res  *httptest.ResponseRecorder
		req  *http.Request
		json bool
	}
	type want struct {
		code        int
		response    string
		contentType string
	}
	tests := []struct {
		name string
		args args
		want want
	}{
		{
			name: "Test valid data",
			args: args{
				res:  httptest.NewRecorder(),
				req:  httptest.NewRequest(http.MethodPost, "http://localhost:8080/updates/", validBody),
				json: true,
			},
			want: want{
				code:        200,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "Test invalid data",
			args: args{
				res:  httptest.NewRecorder(),
				req:  httptest.NewRequest(http.MethodPost, "http://localhost:8080/updates/", invalidBody),
				json: true,
			},
			want: want{
				code:        400,
				response:    "Wrong json provided.\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
		{
			name: "Test GET method",
			args: args{
				res:  httptest.NewRecorder(),
				req:  httptest.NewRequest(http.MethodGet, "http://localhost:8080/updates/", validBody),
				json: true,
			},
			want: want{
				code:        405,
				response:    "",
				contentType: "",
			},
		},
		{
			name: "Test not json",
			args: args{
				res:  httptest.NewRecorder(),
				req:  httptest.NewRequest(http.MethodPost, "http://localhost:8080/updates/", validBody),
				json: false,
			},
			want: want{
				code:        400,
				response:    "Wrong Content-Type, use application/json!\n",
				contentType: "text/plain; charset=utf-8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.args.json {
				tt.args.req.Header.Add("Content-Type", "application/json")
			}
			r := chi.NewRouter()
			prepareRoutes(r)
			r.ServeHTTP(tt.args.res, tt.args.req)
			res := tt.args.res.Result()
			// проверяем код ответа
			assert.Equal(t, tt.want.code, res.StatusCode)

			defer func() {
				_ = res.Body.Close()
			}()
			resBody, err := io.ReadAll(res.Body)
			require.NoError(t, err)
			assert.Equal(t, tt.want.response, string(resBody))
			// проверяем Content-Type
			assert.Equal(t, tt.want.contentType, res.Header.Get("Content-Type"))
		})
	}
}
