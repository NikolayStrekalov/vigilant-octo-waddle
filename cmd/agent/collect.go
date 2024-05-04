package main

import (
	"math/rand"
	"runtime"
	"time"
)

var pollInterval = 2

func collectStats(stats *runtime.MemStats) {
	for {
		m.Lock()
		runtime.ReadMemStats(stats)
		PollCount++
		RandomValue = rand.Float64()
		m.Unlock()
		time.Sleep(time.Duration(pollInterval) * time.Second)
	}
}
