package agent

import (
	"errors"
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
	statPauseTotalNs,
	statStackInuse,
	statStackSys,
	statSys,
	statTotalAlloc,
}

var RuntimeStats = runtime.MemStats{}
var PollCount int
var RandomValue float64
var errStatValueConversion = errors.New("can not convert to float64")

func getFloatStat(stat reflect.Value) (float64, error) {
	switch stat.Kind() {
	case reflect.Uint64:
		if v, ok := stat.Interface().(uint64); ok {
			return float64(v), nil
		}
	case reflect.Uint32:
		if v, ok := stat.Interface().(uint32); ok {
			return float64(v), nil
		}
	case reflect.Int:
		if v, ok := stat.Interface().(int); ok {
			return float64(v), nil
		}
	case reflect.Float64:
		if v, ok := stat.Interface().(float64); ok {
			return v, nil
		}
	default:
		return 0, errStatValueConversion
	}
	return 0, errStatValueConversion
}

func getFormatedStat(stat reflect.Value) string {
	switch stat.Kind() {
	case reflect.Uint64:
		if v, ok := stat.Interface().(uint64); ok {
			return strconv.FormatUint(v, 10)
		}
	case reflect.Uint32:
		if v, ok := stat.Interface().(uint32); ok {
			return strconv.FormatUint(uint64(v), 10)
		}
	case reflect.Int:
		if v, ok := stat.Interface().(int); ok {
			return strconv.Itoa(v)
		}
	case reflect.Float64:
		if v, ok := stat.Interface().(float64); ok {
			return strconv.FormatFloat(v, 'f', -1, 64)
		}
	default:
		return stat.String()
	}
	return stat.String()
}

var statMutex = &sync.Mutex{}
