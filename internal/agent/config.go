package agent

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
