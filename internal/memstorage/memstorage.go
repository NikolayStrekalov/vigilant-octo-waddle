package memstorage

import (
	"errors"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/models"
	"github.com/mailru/easyjson"
)

const dumpFilePermissions = 0o600

//easyjson:json
type MemStorage struct {
	Gauge         map[string]float64
	Counter       map[string]int64
	muxGauge      *sync.RWMutex
	muxCounter    *sync.RWMutex
	dumpFile      string
	sync          bool
	storeInterval time.Duration
}

var errNotFound = errors.New("not found")

const (
	gaugeKind   = "gauge"
	counterKind = "counter"
)

type GaugeListItem = struct {
	Name  string
	Value float64
}

type MetricsSlice = []struct {
	Delta *int64
	Value *float64
	ID    string
	MType string
}

func NewMemStorage(dumpPath string, restore bool, storeInterval int) (*MemStorage, func() error, error) {
	storage := MemStorage{
		Gauge:         make(map[string]float64),
		Counter:       make(map[string]int64),
		muxGauge:      &sync.RWMutex{},
		muxCounter:    &sync.RWMutex{},
		sync:          dumpPath != "" && storeInterval == 0,
		dumpFile:      dumpPath,
		storeInterval: time.Duration(storeInterval) * time.Second,
	}
	if restore {
		storage.restore()
	}
	if dumpPath != "" && storeInterval > 0 {
		go storage.periodicDump()
	}
	var closeStorage = func() error {
		storage.dump()
		return nil
	}
	return &storage, closeStorage, nil
}

func (m *MemStorage) GetGaugeList() []GaugeListItem {
	m.muxGauge.RLock()
	defer m.muxGauge.RUnlock()
	items := make([]GaugeListItem, 0, len(m.Gauge))
	for name, value := range m.Gauge {
		items = append(items, GaugeListItem{Name: name, Value: value})
	}
	return items
}

type CounterListItem = struct {
	Name  string
	Value int64
}

func (m *MemStorage) GetCounterList() []CounterListItem {
	m.muxCounter.RLock()
	defer m.muxCounter.RUnlock()
	items := make([]CounterListItem, 0, len(m.Counter))
	for name, value := range m.Counter {
		items = append(items, CounterListItem{Name: name, Value: value})
	}
	return items
}

func (m *MemStorage) GetGauge(name string) (float64, error) {
	m.muxGauge.RLock()
	defer m.muxGauge.RUnlock()
	if v, ok := m.Gauge[name]; ok {
		return v, nil
	}
	return 0, errNotFound
}

func (m *MemStorage) GetCounter(name string) (int64, error) {
	m.muxCounter.RLock()
	defer m.muxCounter.RUnlock()
	if v, ok := m.Counter[name]; ok {
		return v, nil
	}
	return 0, errNotFound
}

func (m *MemStorage) UpdateGauge(name string, value float64) {
	m.muxGauge.Lock()
	m.Gauge[name] = value
	m.muxGauge.Unlock()
	if m.sync {
		m.dump()
	}
}

func (m *MemStorage) IncrementCounter(name string, value int64) {
	m.muxCounter.Lock()
	m.Counter[name] += value
	m.muxCounter.Unlock()
	if m.sync {
		m.dump()
	}
}

func (m *MemStorage) BulkUpdate(metrics models.MetricsSlice) {
	m.muxCounter.Lock()
	m.muxGauge.Lock()
	for _, metric := range metrics {
		switch metric.MType {
		case counterKind:
			if metric.Delta == nil {
				continue
			}
			m.Counter[metric.ID] += *metric.Delta

		case gaugeKind:
			if metric.Value == nil {
				continue
			}
			m.Gauge[metric.ID] = *metric.Value
		default:
			continue
		}
	}
	m.muxGauge.Unlock()
	m.muxCounter.Unlock()
	if m.sync {
		m.dump()
	}
}

func (m *MemStorage) dump() {
	if m.dumpFile == "" {
		return
	}
	m.muxCounter.RLock()
	m.muxGauge.RLock()
	defer func() {
		m.muxCounter.RUnlock()
		m.muxGauge.RUnlock()
	}()
	data, err := easyjson.Marshal(m)
	if err != nil {
		logger.Info("Error converting storage to json.", err)
		return
	}
	err = os.WriteFile(m.dumpFile, data, dumpFilePermissions)
	if err != nil {
		logger.Info("Error writing dump.", err)
		return
	}
}

func (m *MemStorage) restore() {
	if m.dumpFile == "" {
		return
	}
	m.muxCounter.Lock()
	m.muxGauge.Lock()
	defer func() {
		m.muxCounter.Unlock()
		m.muxGauge.Unlock()
	}()
	data, err := os.ReadFile(m.dumpFile)
	if err != nil {
		logger.Info("Error reading dump.", err)
		return
	}
	err = easyjson.Unmarshal(data, m)
	if err != nil {
		logger.Info("Error unmarshalling dump.", err)
		return
	}
}

func (m *MemStorage) periodicDump() {
	for {
		time.Sleep(m.storeInterval)
		m.dump()
	}
}

func (m *MemStorage) Log() {
	m.muxGauge.RLock()
	fmt.Println(m.Gauge)
	m.muxGauge.RUnlock()
	m.muxCounter.RLock()
	fmt.Println(m.Counter)
	m.muxCounter.RUnlock()
}
