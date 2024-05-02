package main

import (
	"reflect"
	"runtime"
	"strconv"
	"sync"
)

type StatName string

const (
	statAlloc         StatName = "Alloc"
	statBuckHashSys   StatName = "BuckHashSys"
	statFrees         StatName = "Frees"
	statGCCPUFraction StatName = "GCCPUFraction"
)

var statList = [...]StatName{
	statAlloc,
	statBuckHashSys,
	statFrees,
	statGCCPUFraction,
}

var RuntimeStats runtime.MemStats
var PollCount int
var RandomValue float64

func getFormatedStat(stat reflect.Value) string {
	switch stat.Kind() {
	case reflect.Uint64:
		return strconv.FormatUint(stat.Interface().(uint64), 10)
	case reflect.Int:
		return strconv.Itoa(stat.Interface().(int))
	case reflect.Float64:
		return strconv.FormatFloat(stat.Interface().(float64), 'f', -1, 64)
	}
	return ""
}

var m sync.Mutex
