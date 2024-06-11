package server

import (
	"errors"
)

type StorageOperations interface {
	GetGaugeList() []GaugeListItem
	GetCounterList() []CounterListItem
	GetGauge(string) (float64, error)
	GetCounter(string) (int64, error)
	UpdateGauge(string, float64)
	IncrementCounter(string, int64)
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
var errStorageNotInitialized = errors.New("storage require initialization")

type DeferredStorage struct {
	reinit        *func()
	checkResource *func() bool
}

func (d *DeferredStorage) tryInit() {
	if (*d.checkResource)() {
		(*d.reinit)()
	}
}

func (d *DeferredStorage) GetGaugeList() []GaugeListItem {
	d.tryInit()
	return []GaugeListItem{}
}

func (d *DeferredStorage) GetCounterList() []CounterListItem {
	d.tryInit()
	return []CounterListItem{}
}

func (d *DeferredStorage) GetGauge(name string) (float64, error) {
	d.tryInit()
	return 0, errStorageNotInitialized
}

func (d *DeferredStorage) GetCounter(name string) (int64, error) {
	d.tryInit()
	return 0, errStorageNotInitialized
}

func (d *DeferredStorage) UpdateGauge(name string, value float64) {
	d.tryInit()
}

func (d *DeferredStorage) IncrementCounter(name string, value int64) {
	d.tryInit()
}
