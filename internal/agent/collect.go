package agent

import (
	"math/rand"
	"runtime"
	"time"
)

func collectStats(stats *runtime.MemStats) {
	for {
		m.Lock()
		runtime.ReadMemStats(stats)
		PollCount++
		RandomValue = rand.Float64()
		m.Unlock()
		time.Sleep(time.Duration(Config.PollInterval) * time.Second)
	}
}
