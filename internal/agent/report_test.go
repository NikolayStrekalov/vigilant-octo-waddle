package agent

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_sendStat(t *testing.T) {
	type args struct {
		kind  StatKind
		name  StatName
		value string
	}
	tests := []struct {
		name     string
		args     args
		wantPath string
	}{
		{
			name: "Test Post path 1",
			args: args{
				kind:  counterKind,
				name:  statPollCount,
				value: "314",
			},
			wantPath: "/counter/PollCount/314",
		},
		{
			name: "Test Post path 2",
			args: args{
				kind:  gaugeKind,
				name:  statBuckHashSys,
				value: "3.1415",
			},
			wantPath: "/gauge/BuckHashSys/3.1415",
		},
	}
	var requestPath string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Hello, client")
		requestPath = r.URL.Path
	}))
	defer ts.Close()
	serverBase = ts.URL
	for _, tt := range tests {
		sendStat(tt.args.kind, tt.args.name, tt.args.value)
		assert.Equal(t, requestPath, tt.wantPath)
	}
}
