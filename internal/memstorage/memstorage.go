package memstorage

import (
	"errors"
	"fmt"
	"os"
	"sync"

	"github.com/NikolayStrekalov/vigilant-octo-waddle.git/internal/logger"
	"github.com/mailru/easyjson"
)

const dumpFilePermissions = 0o600

//easyjson:json
type MemStorage struct {
	Gauge      map[string]float64
	Counter    map[string]int64
	muxGauge   *sync.RWMutex
	muxCounter *sync.RWMutex
	dumpFile   string
	sync       bool
}

var errNotFound = errors.New("not found")

type GaugeListItem = struct {
	Name  string
	Value float64
}

func NewMemStorage(synchronous bool, dumpPath string) *MemStorage {
	return &MemStorage{
		Gauge:      make(map[string]float64),
		Counter:    make(map[string]int64),
		muxGauge:   &sync.RWMutex{},
		muxCounter: &sync.RWMutex{},
		sync:       synchronous,
		dumpFile:   dumpPath,
	}
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

type CounterListItem = struct {
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
	if m.sync {
		m.Dump()
	}
}

func (m MemStorage) IncrementCounter(name string, value int64) {
	m.muxCounter.Lock()
	m.Counter[name] += value
	m.muxCounter.Unlock()
	if m.sync {
		m.Dump()
	}
}

func (m MemStorage) Dump() {
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

func (m *MemStorage) Restore() {
	if m.dumpFile == "" {
		return
	}
	// FIXME: deadlocks with update operations; current use case is safe
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

func (m MemStorage) Log() {
	m.muxGauge.RLock()
	fmt.Println(m.Gauge)
	m.muxGauge.RUnlock()
	m.muxCounter.RLock()
	fmt.Println(m.Counter)
	m.muxCounter.RUnlock()
}
