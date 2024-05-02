package main

import (
	"fmt"
	"math/rand"
	"reflect"
	"runtime"
	"time"
)

func gatherStats(stats *runtime.MemStats) {
	for {
		m.Lock()
		runtime.ReadMemStats(stats)
		PollCount += 1
		RandomValue = rand.Float64()
		m.Unlock()
		time.Sleep(2 * time.Second)
	}
}

func reportStats(stats *runtime.MemStats) {
	for {
		m.Lock()
		runtime.ReadMemStats(stats)
		r := reflect.ValueOf(*stats)
		for _, statName := range statList {
			f := reflect.Indirect(r).FieldByName(string(statName))
			fmt.Printf("%s: %s\n", statName, getFormatedStat(f))
		}
		fmt.Printf("%s: %s\n", "RandomValue", getFormatedStat(reflect.ValueOf(RandomValue)))
		fmt.Printf("%s: %s\n", "PollCount", getFormatedStat(reflect.ValueOf(PollCount)))
		fmt.Println()
		PollCount = 0
		m.Unlock()
		time.Sleep(10 * time.Second)
	}
}

func main() {
	RuntimeStats = runtime.MemStats{}
	go gatherStats(&RuntimeStats)
	go reportStats(&RuntimeStats)
	runtime.Goexit()
}
