package agent

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_getFormatedStat(t *testing.T) {
	type args struct {
		stat reflect.Value
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test Float",
			args: args{stat: reflect.ValueOf(float64(5.83))},
			want: "5.83",
		},
		{
			name: "Test Uint64",
			args: args{stat: reflect.ValueOf(uint64(882))},
			want: "882",
		},
		{
			name: "Test Uint32",
			args: args{stat: reflect.ValueOf(uint32(583))},
			want: "583",
		},
		{
			name: "Test Int",
			args: args{stat: reflect.ValueOf(int(583))},
			want: "583",
		},
		{
			name: "Test Bool",
			args: args{stat: reflect.ValueOf(true)},
			want: "<bool Value>",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, getFormatedStat(tt.args.stat), tt.want)
		})
	}
}
