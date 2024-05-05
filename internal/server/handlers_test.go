package server

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_updateHandler(t *testing.T) {
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
				response:    "Use format http://<АДРЕС_СЕРВЕРА>/update/<ТИП_МЕТРИКИ>/<ИМЯ_МЕТРИКИ>/<ЗНАЧЕНИЕ_МЕТРИКИ>\n",
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			updateHandler(tt.args.res, tt.args.req)
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
