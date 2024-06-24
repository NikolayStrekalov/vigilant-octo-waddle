package agent

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"golang.org/x/sync/semaphore"
)

type Conf struct {
	ServerAddress  string `json:"address"`
	SignKey        string `json:"key"`
	PollInterval   uint   `json:"poll"`
	ReportInterval uint   `json:"report"`
	RateLimit      uint   `json:"limit"`
}

func (s *Conf) log() {
	lg, err := json.Marshal(Config)
	if err != nil {
		fmt.Println("error serializing config:", err)
	}
	fmt.Println("config:", string(lg))
}

var defaultPollInterval uint = 2
var defalutReportInterval uint = 10
var defaultRateLimit uint = 1
var Config = Conf{
	PollInterval:   defaultPollInterval,
	ReportInterval: defalutReportInterval,
	RateLimit:      defaultRateLimit,
	ServerAddress:  "localhost:8080",
	SignKey:        "",
}
var RequestLimiter *semaphore.Weighted

func setupConfig() {
	flag.StringVar(&Config.ServerAddress, "a", "localhost:8080", "Эндпоинт сервера HOST:PORT")
	flag.StringVar(&Config.SignKey, "k", "", "Ключ для подписи")
	flag.UintVar(&Config.PollInterval, "p", Config.PollInterval, "Частота опроса метрик в секундах, больше нуля")
	flag.UintVar(&Config.ReportInterval, "r", Config.ReportInterval, "Частота отправки метрик в секундах, больше нуля")
	flag.UintVar(&Config.RateLimit, "l", Config.RateLimit, "Максимальное число одновременных исходящих запросов")
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
	if envSignKey := os.Getenv("KEY"); envSignKey != "" {
		Config.SignKey = envSignKey
	}

	ReportBaseURL = fmt.Sprintf("http://%s/update/", Config.ServerAddress)
	ReportBulkURL = fmt.Sprintf("http://%s/updates/", Config.ServerAddress)
	RequestLimiter = semaphore.NewWeighted(int64(Config.RateLimit))

	Config.log()
}
