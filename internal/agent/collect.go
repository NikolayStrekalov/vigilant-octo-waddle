package agent

import (
	"fmt"
	"math/rand"
	"runtime"
	"time"

	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/mem"
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

func collectGopsutilStats() {
	for {
		v, err := mem.VirtualMemory()
		if err != nil {
			fmt.Println("Error reading memory stat.")
		} else {
			GopsutilStats.Store("FreeMemory", float64(v.Free))
			GopsutilStats.Store("TotalMemory", float64(v.Total))
		}
		cpus, err := cpu.Percent(0, true)
		if err != nil {
			fmt.Println("Error reading cpu stat.")
		} else {
			for i, percent := range cpus {
				GopsutilStats.Store(fmt.Sprintf("CPUutilization%d", i+1), percent)
			}
		}
		time.Sleep(time.Duration(Config.PollInterval) * time.Second)
	}
}
