package agent

import (
	"runtime"
)

var exitCodeWrongArgs = 2
var exitCodeMisconfigured = 3

func Start() {
	setupConfig()
	go collectStats()
	go collectGopsutilStats()
	go reportStats()
	runtime.Goexit()
}
