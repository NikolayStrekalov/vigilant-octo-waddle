package agent

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"time"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
	"github.com/mailru/easyjson"
)

var ReportBaseURL = "http://localhost:8080/update/"
var ReportBulkURL = "http://localhost:8080/updates/"

func sendStat(kind StatKind, name StatName, value string) {
	path, err := url.JoinPath(ReportBaseURL, string(kind), string(name), value)
	if err != nil {
		fmt.Println("Fail to construct server url.", err)
		return
	}
	resp, err := http.Post(path, "text/plain", http.NoBody)
	if err != nil {
		fmt.Println("Post error:", err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		fmt.Println("Wrong request code:", resp.StatusCode)
	}
}

func sendStatJSON(m easyjson.Marshaler, toURL string) {
	data, err := easyjson.Marshal(m)
	if err != nil {
		fmt.Println("Fail to serialize metric.", err)
		return
	}
	gzData, err := Compress(data)
	if err != nil {
		fmt.Println("Compress error:", err)
		return
	}
	req, err := http.NewRequest(http.MethodPost, toURL, bytes.NewReader(gzData))
	if err != nil {
		fmt.Println("Create request error:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Post error:", err)
		return
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Response read error:", err)
			return
		}
		fmt.Println(ReportBaseURL, "Wrong request code:", resp.StatusCode, string(data))
	}
}

func reportStats() {
	var err error
	for {
		statMutex.Lock()
		runtime.ReadMemStats(&RuntimeStats)
		r := reflect.ValueOf(RuntimeStats)
		metrics := models.MetricsSlice{}
		for _, statName := range runtimeStatList {
			f := reflect.Indirect(r).FieldByName(string(statName))

			runtimeMetrics := models.Metrics{
				ID:    string(statName),
				MType: "gauge",
				Value: new(float64),
			}
			if *runtimeMetrics.Value, err = getFloatStat(f); err != nil {
				fmt.Println(err)
			}
			metrics = append(metrics, runtimeMetrics)
		}
		randomMetrics := models.Metrics{
			ID:    string(statRandomValue),
			MType: "gauge",
			Value: new(float64),
		}
		*randomMetrics.Value = RandomValue
		metrics = append(metrics, randomMetrics)
		pollMetrics := models.Metrics{
			ID:    string(statPollCount),
			MType: "counter",
			Delta: new(int64),
		}
		*pollMetrics.Delta = int64(PollCount)
		metrics = append(metrics, pollMetrics)
		PollCount = 0
		statMutex.Unlock()

		go sendStatJSON(metrics, ReportBulkURL)
		time.Sleep(time.Duration(Config.ReportInterval) * time.Second)
	}
}
