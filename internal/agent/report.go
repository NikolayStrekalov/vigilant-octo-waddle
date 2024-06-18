package agent

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"time"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
	"github.com/mailru/easyjson"

	"github.com/avast/retry-go/v4"
)

var ReportBaseURL = "http://localhost:8080/update/"
var ReportBulkURL = "http://localhost:8080/updates/"

const maxRequestAttempts = 4

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

func sendStatJSON(m easyjson.Marshaler, toURL string) error {
	data, err := easyjson.Marshal(m)
	if err != nil {
		return fmt.Errorf("fail to serialize metric: %w", err)
	}
	gzData, err := Compress(data)
	if err != nil {
		return fmt.Errorf("compress error: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, toURL, bytes.NewReader(gzData))
	if err != nil {
		return fmt.Errorf("create request error: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Content-Encoding", "gzip")

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*time.Duration(Config.ReportInterval))
	defer cancel()
	resp, err := retry.DoWithData(
		func() (*http.Response, error) {
			return http.DefaultClient.Do(req)
		},
		retry.Context(ctx),
		retry.Attempts(maxRequestAttempts),
		retry.DelayType(func(n uint, err error, config *retry.Config) time.Duration {
			return time.Duration(1+n*2) * time.Second
		}),
	)
	if err != nil {
		return fmt.Errorf("post error: %w", err)
	}
	defer func() {
		_ = resp.Body.Close()
	}()
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf("response read error: %w", err)
		}
		return fmt.Errorf("wrong response code: %d, data: %s", resp.StatusCode, string(data))
	}
	return nil
}

func reportStats() {
	var err error
	ticker := time.NewTicker(time.Duration(Config.ReportInterval) * time.Second)
	for {
		<-ticker.C
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
		*pollMetrics.Delta = PollCount
		metrics = append(metrics, pollMetrics)
		PollCount = 0
		statMutex.Unlock()

		err := sendStatJSON(metrics, ReportBulkURL)
		if err != nil {
			fmt.Println(err)
			statMutex.Lock()
			PollCount += *pollMetrics.Delta
			statMutex.Unlock()
		}
	}
}
