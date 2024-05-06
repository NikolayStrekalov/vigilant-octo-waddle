package agent

import (
	"flag"
	"fmt"
	"os"
	"strconv"
)

type Conf struct {
	ServerAddress  string
	PollInterval   uint
	ReportInterval uint
}

var defaultPollInterval uint = 2
var defalutReportInterval uint = 10
var Config = Conf{
	PollInterval:   defaultPollInterval,
	ReportInterval: defalutReportInterval,
	ServerAddress:  "localhost:8080",
}

func setupConfig() {
	flag.StringVar(&Config.ServerAddress, "a", "localhost:8080", "Эндпоинт сервера HOST:PORT")
	flag.UintVar(&Config.PollInterval, "p", Config.PollInterval, "Частота опроса метрик в секундах, больше нуля")
	flag.UintVar(&Config.ReportInterval, "r", Config.ReportInterval, "Частота отправки метрик в секундах, больше нуля")
	flag.Parse()
	if len(flag.Args()) > 0 || Config.PollInterval == 0 || Config.ReportInterval == 0 {
		flag.PrintDefaults()
		os.Exit(exitCodeWrongArgs)
	}

	if envAddress := os.Getenv("ADDRESS"); envAddress != "" {
		Config.ServerAddress = envAddress
	}
	if reportInterval := os.Getenv("REPORT_INTERVAL"); reportInterval != "" {
		val, err := strconv.Atoi(reportInterval)
		if err != nil || val <= 0 {
			fmt.Println("wrong REPORT_INTERVAL value")
			os.Exit(exitCodeMisconfigured)
		}
		Config.ReportInterval = uint(val)
	}
	if pollInterval := os.Getenv("POLL_INTERVAL"); pollInterval != "" {
		val, err := strconv.Atoi(pollInterval)
		if err != nil || val <= 0 {
			fmt.Println("wrong POLL_INTERVAL value")
			os.Exit(exitCodeMisconfigured)
		}
		Config.PollInterval = uint(val)
	}

	ReportBaseURL = fmt.Sprintf("http://%s/update/", Config.ServerAddress)
}
