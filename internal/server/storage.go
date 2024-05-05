package server

import (
	"errors"
	"fmt"
	"sync"
)

type StorageOperations interface {
	GetGauge(string, float64)
	GetCounter(string, int64)
	UpdateGauge(string, float64)
	IncrementCounter(string, int64)
}

type MemStorage struct {
	gauge      map[string]float64
	counter    map[string]int64
	muxGauge   *sync.RWMutex
	muxCounter *sync.RWMutex
}

var storage = MemStorage{
	gauge:      make(map[string]float64),
	counter:    make(map[string]int64),
	muxGauge:   &sync.RWMutex{},
	muxCounter: &sync.RWMutex{},
}

var errNotFound = errors.New("not found")

func (m MemStorage) getGauge(name string) (float64, error) {
	m.muxGauge.RLock()
	defer m.muxGauge.RUnlock()
	if v, ok := m.gauge[name]; ok {
		return v, nil
	}
	return 0, errNotFound
}

func (m MemStorage) getCounter(name string) (int64, error) {
	m.muxCounter.RLock()
	defer m.muxCounter.RUnlock()
	if v, ok := m.counter[name]; ok {
		return v, nil
	}
	return 0, errNotFound
}

func (m MemStorage) UpdateGauge(name string, value float64) {
	m.muxGauge.Lock()
	m.gauge[name] = value
	m.muxGauge.Unlock()
}

func (m MemStorage) IncrementCounter(name string, value int64) {
	m.muxCounter.Lock()
	m.counter[name] += value
	m.muxCounter.Unlock()
}

func (m MemStorage) Log() {
	m.muxGauge.RLock()
	fmt.Println(storage.gauge)
	m.muxGauge.RUnlock()
	m.muxCounter.RLock()
	fmt.Println(storage.counter)
	m.muxCounter.RUnlock()
}
