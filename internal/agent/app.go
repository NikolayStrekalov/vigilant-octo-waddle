package agent

import (
	"flag"
	"fmt"
	"os"
	"runtime"
)

var exitCodeWrongArgs = 2

func Start() {
	flag.StringVar(&Config.ServerAddress, "a", "localhost:8080", "Эндпоинт сервера HOST:PORT")
	flag.UintVar(&Config.PollInterval, "p", Config.PollInterval, "Частота опроса метрик в секундах, больше нуля")
	flag.UintVar(&Config.ReportInterval, "r", Config.ReportInterval, "Частота отправки метрик в секундах, больше нуля")
	flag.Parse()
	if len(flag.Args()) > 0 || Config.PollInterval == 0 || Config.ReportInterval == 0 {
		flag.PrintDefaults()
		os.Exit(exitCodeWrongArgs)
	}
	ReportBaseURL = fmt.Sprintf("http://%s/update/", Config.ServerAddress)

	go collectStats()
	go reportStats()
	runtime.Goexit()
}
