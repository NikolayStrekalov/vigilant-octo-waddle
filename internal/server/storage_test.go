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

func TestMemStorage_getGauge(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    float64
		wantErr bool
	}{
		{
			name: "Get gauge",
			fields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{"PollCount": 30},
			},
			args: args{
				name: "RandomValue",
			},
			want:    0.31,
			wantErr: false,
		},
		{
			name: "Miss gauge",
			fields: fields{
				gauge:   map[string]float64{"RandomValue": 0.31},
				counter: map[string]int64{"PollCount": 30},
			},
			args: args{
				name: "Random",
			},
			want:    0,
			wantErr: true,
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
			got, err := m.getGauge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.getGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.getGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_getCounter(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "Get counter",
			fields: fields{
				gauge:   map[string]float64{"Value": 0.31},
				counter: map[string]int64{"Poll": 30},
			},
			args: args{
				name: "Poll",
			},
			want:    30,
			wantErr: false,
		},
		{
			name: "Miss counter",
			fields: fields{
				gauge:   map[string]float64{"Value": 0.31},
				counter: map[string]int64{"Poll": 30},
			},
			args: args{
				name: "Count",
			},
			want:    0,
			wantErr: true,
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
			got, err := m.getCounter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.getCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.getCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}
