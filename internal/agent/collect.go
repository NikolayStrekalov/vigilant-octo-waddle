package agent

import (
	"math/rand"
	"runtime"
	"time"
)

func collectStats() {
	for {
		statMutex.Lock()
		runtime.ReadMemStats(&RuntimeStats)
		PollCount++
		RandomValue = rand.Float64()
		statMutex.Unlock()

		time.Sleep(time.Duration(Config.PollInterval) * time.Second)
	}
}
