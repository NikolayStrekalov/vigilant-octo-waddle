package server

import (
	"reflect"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMemStorage_UpdateGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name  string
		value float64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFields fields
	}{
		// TODO: Add test cases.
		{
			name: "test RandomValue 1",
			fields: fields{
				gauge:   map[string]float64{},
				counter: map[string]int64{},
			},
			args: args{
				name:  "RandomValue",
				value: 0.31,
			},
			wantFields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{},
			},
		},
		{
			name: "test RandomValue 2",
			fields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{"PollCount": 62},
			},
			args: args{
				name:  "RandomValue2",
				value: 0.3,
			},
			wantFields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31, "RandomValue2": 0.3},
				counter: map[string]int64{"PollCount": 62},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStorage{
				gauge:      tt.fields.gauge,
				counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			m.UpdateGauge(tt.args.name, tt.args.value)
			assert.True(t, reflect.DeepEqual(m.gauge, tt.wantFields.gauge))
			assert.True(t, reflect.DeepEqual(m.counter, tt.wantFields.counter))
		})
	}
}

func TestMemStorage_IncrementCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name  string
		value int64
	}
	tests := []struct {
		name       string
		fields     fields
		args       args
		wantFields fields
	}{
		{
			name: "test PollCount 1",
			fields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{},
			},
			args: args{
				name:  "PollCount",
				value: 32,
			},
			wantFields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{"PollCount": 32},
			},
		},
		{
			name: "test PollCount 2",
			fields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{"PollCount": 30},
			},
			args: args{
				name:  "PollCount",
				value: 32,
			},
			wantFields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{"PollCount": 62},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStorage{
				gauge:      tt.fields.gauge,
				counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			m.IncrementCounter(tt.args.name, tt.args.value)
			assert.True(t, reflect.DeepEqual(m.gauge, tt.wantFields.gauge))
			assert.True(t, reflect.DeepEqual(m.counter, tt.wantFields.counter))
		})
	}
}
