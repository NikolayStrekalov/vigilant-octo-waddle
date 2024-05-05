package agent

import (
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"runtime"
	"time"
)

var ReportBaseURL = "http://localhost:8080/update/"

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

func reportStats(stats *runtime.MemStats) {
	for {
		m.Lock()
		runtime.ReadMemStats(stats)
		r := reflect.ValueOf(*stats)
		for _, statName := range runtimeStatList {
			f := reflect.Indirect(r).FieldByName(string(statName))
			go sendStat(gaugeKind, statName, getFormatedStat(f))
		}
		go sendStat(gaugeKind, statRandomValue, getFormatedStat(reflect.ValueOf(RandomValue)))
		go sendStat(counterKind, statPollCount, getFormatedStat(reflect.ValueOf(PollCount)))
		PollCount = 0
		m.Unlock()
		time.Sleep(time.Duration(Config.ReportInterval) * time.Second)
	}
}
