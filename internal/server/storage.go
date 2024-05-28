package server

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/mailru/easyjson"
)

const dumpFilePermissions = 0o600

type StorageOperations interface {
	GetGaugeList() []GaugeListItem
	GetCounterList() []CounterListItem
	GetGauge(string, float64)
	GetCounter(string, int64)
	UpdateGauge(string, float64)
	IncrementCounter(string, int64)
	Dump(string)
	Restore(string)
}

//easyjson:json
type MemStorage struct {
	Gauge      map[string]float64
	Counter    map[string]int64
	muxGauge   *sync.RWMutex
	muxCounter *sync.RWMutex
}

var storage = MemStorage{
	Gauge:      make(map[string]float64),
	Counter:    make(map[string]int64),
	muxGauge:   &sync.RWMutex{},
	muxCounter: &sync.RWMutex{},
}

var errNotFound = errors.New("not found")

type GaugeListItem struct {
	Name  string
	Value float64
}

func (m MemStorage) GetGaugeList() []GaugeListItem {
	m.muxGauge.RLock()
	defer m.muxGauge.RUnlock()
	items := make([]GaugeListItem, 0, len(m.Gauge))
	for name, value := range m.Gauge {
		items = append(items, GaugeListItem{Name: name, Value: value})
	}
	return items
}

type CounterListItem struct {
	Name  string
	Value int64
}

func (m MemStorage) GetCounterList() []CounterListItem {
	m.muxCounter.RLock()
	defer m.muxCounter.RUnlock()
	items := make([]CounterListItem, 0, len(m.Counter))
	for name, value := range m.Counter {
		items = append(items, CounterListItem{Name: name, Value: value})
	}
	return items
}

func (m MemStorage) GetGauge(name string) (float64, error) {
	m.muxGauge.RLock()
	defer m.muxGauge.RUnlock()
	if v, ok := m.Gauge[name]; ok {
		return v, nil
	}
	return 0, errNotFound
}

func (m MemStorage) GetCounter(name string) (int64, error) {
	m.muxCounter.RLock()
	defer m.muxCounter.RUnlock()
	if v, ok := m.Counter[name]; ok {
		return v, nil
	}
	return 0, errNotFound
}

func (m MemStorage) UpdateGauge(name string, value float64) {
	m.muxGauge.Lock()
	m.Gauge[name] = value
	m.muxGauge.Unlock()
}

func (m MemStorage) IncrementCounter(name string, value int64) {
	m.muxCounter.Lock()
	m.Counter[name] += value
	m.muxCounter.Unlock()
}

func (m MemStorage) Dump(filePath string) {
	m.muxCounter.RLock()
	m.muxGauge.RLock()
	defer func() {
		m.muxCounter.RUnlock()
		m.muxGauge.RUnlock()
	}()
	data, err := easyjson.Marshal(storage)
	if err != nil {
		logger.Info("Error converting storage to json.", err)
		return
	}
	err = os.WriteFile(filePath, data, dumpFilePermissions)
	if err != nil {
		logger.Info("Error writing dump.", err)
		return
	}
}

func (m MemStorage) Restore(filePath string) {
	m.muxCounter.Lock()
	m.muxGauge.Lock()
	defer func() {
		m.muxCounter.Unlock()
		m.muxGauge.Unlock()
	}()
	data, err := os.ReadFile(filePath)
	if err != nil {
		logger.Info("Error reading dump.", err)
		return
	}
	err = easyjson.Unmarshal(data, &storage)
	if err != nil {
		logger.Info("Error unmarshalling dump.", err)
		return
	}
}

func (m MemStorage) Log() {
	m.muxGauge.RLock()
	fmt.Println(storage.Gauge)
	m.muxGauge.RUnlock()
	m.muxCounter.RLock()
	fmt.Println(storage.Counter)
	m.muxCounter.RUnlock()
}
