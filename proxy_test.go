package proxy

import (
	"bytes"
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
				r: httptest.NewRequest(http.MethodPost, "http://example.com",
					bytes.NewBufferString(`{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0xd868711BD9a2C6F1548F5f4737f71DA67d821090","latest"]}`),
				),
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
