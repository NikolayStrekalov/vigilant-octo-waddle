package agent

import (
	"math/rand"
	"runtime"
	"time"
)

func collectStats() {
	for {
		m.Lock()
		runtime.ReadMemStats(&RuntimeStats)
		PollCount++
		RandomValue = rand.Float64()
		m.Unlock()
		time.Sleep(time.Duration(Config.PollInterval) * time.Second)
	}
}
