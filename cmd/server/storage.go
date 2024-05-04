package main

import (
	"fmt"
	"sync"
)

type StorageOperations interface {
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
