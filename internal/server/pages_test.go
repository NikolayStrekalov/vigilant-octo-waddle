package server

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var example = `<!DOCTYPE html>
<html lang="en">
  <head>
    <meta charset="utf-8">
    <title>Metrics list</title>
  </head>
  <body>
	<ul>
	<li>gauge poi 4.292772423</li>
	<li>gauge lsd rew 3e&#43;11</li>
	<li>counter Rw 43</li>
	<li>counter qwe -4</li>
	</ul>
  </body>
</html>
`

func Test_renderIndexPage(t *testing.T) {
	type args struct {
		counters []CounterListItem
		gauges   []GaugeListItem
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		// TODO: Add test cases.
		{
			name: "Simple list",
			args: args{
				counters: []CounterListItem{
					{
						Name:  "Rw",
						Value: 43,
					},
					{
						Name:  "qwe",
						Value: -4,
					},
				},
				gauges: []GaugeListItem{
					{
						Name:  "poi",
						Value: 4.292772423,
					},
					{
						Name:  "lsd rew",
						Value: 300000000000,
					},
				},
			},
			want: example,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlBuf, err := renderIndexPage(tt.args.counters, tt.args.gauges)
			assert.Nil(t, err)
			assert.Equal(t, tt.want, htmlBuf.String())
		})
	}
}
