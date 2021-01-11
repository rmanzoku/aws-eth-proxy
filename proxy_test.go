package proxy

import (
	"bytes"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

var getBalanceReader = bytes.NewReader([]byte(`{"jsonrpc":"2.0","id":1,"method":"eth_getBalance","params":["0xd868711BD9a2C6F1548F5f4737f71DA67d821090","latest"]}`))

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
					getBalanceReader,
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

func TestGetSignature(t *testing.T) {
	type args struct {
		r    *http.Request
		body io.ReadSeeker
	}
	tests := []struct {
		name    string
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "ok",
			args: args{
				r:    httptest.NewRequest(http.MethodPost, httpEndpoint, getBalanceReader),
				body: getBalanceReader,
			},
			want:    []byte{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetSignature(tt.args.r, tt.args.body)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetSignature() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Log(got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetSignature() = %v, want %v", got, tt.want)
			}
		})
	}
}
