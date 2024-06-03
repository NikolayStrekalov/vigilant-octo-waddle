package memstorage

import (
	"os"
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
				Gauge:      tt.fields.gauge,
				Counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
				sync:       false,
			}
			m.UpdateGauge(tt.args.name, tt.args.value)
			assert.True(t, reflect.DeepEqual(m.Gauge, tt.wantFields.gauge))
			assert.True(t, reflect.DeepEqual(m.Counter, tt.wantFields.counter))
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
				Gauge:      tt.fields.gauge,
				Counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			m.IncrementCounter(tt.args.name, tt.args.value)
			assert.True(t, reflect.DeepEqual(m.Gauge, tt.wantFields.gauge))
			assert.True(t, reflect.DeepEqual(m.Counter, tt.wantFields.counter))
		})
	}
}

func TestMemStorage_GetGauge(t *testing.T) {
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
				Gauge:      tt.fields.gauge,
				Counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			got, err := m.GetGauge(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetGauge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetGauge() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetCounter(t *testing.T) {
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
				Gauge:      tt.fields.gauge,
				Counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			got, err := m.GetCounter(tt.args.name)
			if (err != nil) != tt.wantErr {
				t.Errorf("MemStorage.GetCounter() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("MemStorage.GetCounter() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMemStorage_GetGaugeList(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []GaugeListItem
	}{
		{
			name: "Gauge List",
			fields: fields{
				gauge:   map[string]float64{"Random": 0.31, "Value": 0.003},
				counter: map[string]int64{"Count": 62, "Poll": -2},
			},
			want: []GaugeListItem{
				{
					Name:  "Random",
					Value: 0.31,
				},
				{
					Name:  "Value",
					Value: 0.003,
				},
			},
		}, {
			name: "Gauge Empty List",
			fields: fields{
				gauge:   map[string]float64{},
				counter: map[string]int64{"Count": 62, "Poll": -2},
			},
			want: []GaugeListItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStorage{
				Gauge:      tt.fields.gauge,
				Counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			got := m.GetGaugeList()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestMemStorage_GetCounterList(t *testing.T) {
	type fields struct {
		gauge   map[string]float64
		counter map[string]int64
	}
	tests := []struct {
		name   string
		fields fields
		want   []CounterListItem
	}{
		{
			name: "Counter List",
			fields: fields{
				gauge:   map[string]float64{"Random": 0.31, "Value": 0.003},
				counter: map[string]int64{"Count": 62, "Poll": -2},
			},
			want: []CounterListItem{
				{
					Name:  "Count",
					Value: 62,
				},
				{
					Name:  "Poll",
					Value: -2,
				},
			},
		}, {
			name: "Counter Empty List",
			fields: fields{
				gauge:   map[string]float64{"Random": 0.31, "Value": 0.003},
				counter: map[string]int64{},
			},
			want: []CounterListItem{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := MemStorage{
				Gauge:      tt.fields.gauge,
				Counter:    tt.fields.counter,
				muxGauge:   &sync.RWMutex{},
				muxCounter: &sync.RWMutex{},
			}
			got := m.GetCounterList()
			assert.ElementsMatch(t, tt.want, got)
		})
	}
}

func TestMemStorageSync(t *testing.T) {
	f, err := os.CreateTemp("", "tmpfile-") // in Go version older than 1.17 you can use ioutil.TempFile
	if err != nil {
		t.Errorf("create temp file error: %v", err)
		return
	}
	defer func() {
		_ = f.Close()
	}()
	defer func() {
		_ = os.Remove(f.Name())
	}()
	storage := NewMemStorage(true, f.Name())
	storage.IncrementCounter("some", 10)
	storage.UpdateGauge("any", 3.1415)
	data, err := os.ReadFile(f.Name())
	if err != nil {
		t.Errorf("reading temp file error: %v", err)
		return
	}
	assert.Equal(t, `{"Gauge":{"any":3.1415},"Counter":{"some":10}}`, string(data))
}
