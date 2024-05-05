package agent

import (
	"runtime"
)

func Start() {
	RuntimeStats = runtime.MemStats{}
	go collectStats(&RuntimeStats)
	go reportStats(&RuntimeStats)
	runtime.Goexit()
}
