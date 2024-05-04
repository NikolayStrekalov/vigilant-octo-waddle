package main

import (
	"reflect"
	"runtime"
	"strconv"
	"sync"
)

type StatKind string

const (
	gaugeKind   StatKind = "gauge"
	counterKind StatKind = "counter"
)

type StatName string

const (
	statAlloc         StatName = "Alloc"
	statBuckHashSys   StatName = "BuckHashSys"
	statFrees         StatName = "Frees"
	statGCCPUFraction StatName = "GCCPUFraction"
	statGCSys         StatName = "GCSys"
	statHeapAlloc     StatName = "HeapAlloc"
	statHeapIdle      StatName = "HeapIdle"
	statHeapInuse     StatName = "HeapInuse"
	statHeapObjects   StatName = "HeapObjects"
	statHeapReleased  StatName = "HeapReleased"
	statHeapSys       StatName = "HeapSys"
	statLastGC        StatName = "LastGC"
	statLookups       StatName = "Lookups"
	statMCacheInuse   StatName = "MCacheInuse"
	statMCacheSys     StatName = "MCacheSys"
	statMSpanInuse    StatName = "MSpanInuse"
	statMSpanSys      StatName = "MSpanSys"
	statMallocs       StatName = "Mallocs"
	statNextGC        StatName = "NextGC"
	statNumForcedGC   StatName = "NumForcedGC"
	statNumGC         StatName = "NumGC"
	statOtherSys      StatName = "OtherSys"
	statPauseTotalNs  StatName = "PauseTotalNs"
	statStackInuse    StatName = "StackInuse"
	statStackSys      StatName = "StackSys"
	statSys           StatName = "Sys"
	statTotalAlloc    StatName = "TotalAlloc"
	statRandomValue   StatName = "RandomValue"
	statPollCount     StatName = "PollCount"
)

var runtimeStatList = [...]StatName{
	statAlloc,
	statBuckHashSys,
	statFrees,
	statGCCPUFraction,
	statGCSys,
	statHeapAlloc,
	statHeapIdle,
	statHeapInuse,
	statHeapObjects,
	statHeapReleased,
	statHeapSys,
	statLastGC,
	statLookups,
	statMCacheInuse,
	statMCacheSys,
	statMSpanInuse,
	statMSpanSys,
	statMallocs,
	statNextGC,
	statNumForcedGC,
	statNumGC,
	statOtherSys,
	statStackInuse,
	statStackSys,
	statSys,
	statTotalAlloc,
}

var RuntimeStats runtime.MemStats
var PollCount int
var RandomValue float64

func getFormatedStat(stat reflect.Value) string {
	switch stat.Kind() {
	case reflect.Uint64:
		return strconv.FormatUint(stat.Interface().(uint64), 10)
	case reflect.Uint32:
		return strconv.FormatUint(uint64(stat.Interface().(uint32)), 10)
	case reflect.Int:
		return strconv.Itoa(stat.Interface().(int))
	case reflect.Float64:
		return strconv.FormatFloat(stat.Interface().(float64), 'f', -1, 64)
	}
	return stat.String()
}

var m = &sync.Mutex{}
