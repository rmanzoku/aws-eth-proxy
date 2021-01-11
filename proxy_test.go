package proxy

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

var getBalance = []byte(`{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0xd868711BD9a2C6F1548F5f4737f71DA67d821090","latest"]}`)

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
				r: httptest.NewRequest(http.MethodPost, "http://example.com", bytes.NewReader(getBalance)),
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

func TestSign(t *testing.T) {
	type args struct {
		r    *http.Request
		body []byte
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "sign",
			args: args{r: httptest.NewRequest(http.MethodPost, httpEndpoint, bytes.NewReader(getBalance)),
				body: getBalance,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Sign(tt.args.r, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("Sign() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got.Header.Get("Authorization"))
			// if !reflect.DeepEqual(got, tt.want) {
			// 	t.Errorf("Sign() = %v, want %v", got, tt.want)
			// }
		})
	}
}
