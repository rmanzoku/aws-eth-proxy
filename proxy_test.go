package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestEthProxyHandler(t *testing.T) {

	type args struct {
		r *http.Request
	}
	tests := []struct {
		name string
		args args
	}{
		{
			name: "ok",
			args: args{
				r: httptest.NewRequest(http.MethodPost, "http://example.com/", nil),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := httptest.NewRecorder()
			EthProxyHandler(got, tt.args.r)
			t.Log(got.Body.String())
		})
	}
}
