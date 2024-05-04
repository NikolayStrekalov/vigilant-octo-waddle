package main

import (
	"runtime"
)

func main() {
	RuntimeStats = runtime.MemStats{}
	go collectStats(&RuntimeStats)
	go reportStats(&RuntimeStats)
	runtime.Goexit()
}
