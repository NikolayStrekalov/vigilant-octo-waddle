package server

type StorageOperations interface {
	GetGaugeList() []GaugeListItem
	GetCounterList() []CounterListItem
	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)
	UpdateGauge(string, float64)
	IncrementCounter(string, int64)
	Dump()
	Restore()
}

type GaugeListItem = struct {
	Name  string
	Value float64
}

type CounterListItem = struct {
	Name  string
	Value int64
}

var Storage StorageOperations
