package server

import (
	"bytes"
	"compress/gzip"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/memstorage"
	"github.com/stretchr/testify/assert"
)

// Decompress распаковывает слайс байт.
func Decompress(data []byte) []byte {
	r, err := gzip.NewReader(bytes.NewReader(data))
	if err != nil {
		panic("Can't create reader")
	}
	defer func() { _ = r.Close() }()

	var b bytes.Buffer
	_, err = b.ReadFrom(r)
	if err != nil {
		panic("Can't decompress")
	}

	return b.Bytes()
}

func Compress(data []byte) []byte {
	var b bytes.Buffer
	w := gzip.NewWriter(&b)
	_, _ = w.Write(data)
	_ = w.Close()
	return b.Bytes()
}

func Test_appRouter(t *testing.T) {
	_ = logger.InitLog()
	r := appRouter()
	Storage, _, _ = memstorage.NewMemStorage("", false, 300)
	ts := httptest.NewServer(r)
	defer ts.Close()
	baseURL := ts.URL

	// Simple update and value requests
	res, err := http.Post(baseURL+"/update/counter/Name/43", "text/html", http.NoBody)
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	data, err := io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "", string(Decompress(data)))
	_ = res.Body.Close()

	res, err = http.Get(baseURL + "/value/counter/Name")
	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	data, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.Equal(t, "43", string(Decompress(data)))
	_ = res.Body.Close()

	// JSON update and value requests
	data = Compress([]byte(`{"id":"34","type":"gauge","value":74.092}`))
	req, _ := http.NewRequest(http.MethodPost, baseURL+"/update/", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	res, err = http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	data, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"id":"34","type":"gauge","value":74.092}`, string(Decompress(data)))
	_ = res.Body.Close()

	data = Compress([]byte(`{"id":"34","type":"gauge"}`))
	req, _ = http.NewRequest(http.MethodPost, baseURL+"/value/", bytes.NewReader(data))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")
	req.Header.Set("Accept-Encoding", "gzip")
	res, err = http.DefaultClient.Do(req)

	assert.Nil(t, err)
	assert.Equal(t, 200, res.StatusCode)
	assert.Equal(t, "gzip", res.Header.Get("Content-Encoding"))
	assert.Equal(t, "application/json", res.Header.Get("Content-Type"))
	data, err = io.ReadAll(res.Body)
	assert.Nil(t, err)
	assert.JSONEq(t, `{"id":"34","type":"gauge","value":74.092}`, string(Decompress(data)))
	_ = res.Body.Close()
}
