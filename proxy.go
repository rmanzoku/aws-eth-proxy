package proxy

import (
	"io"
	"net/http"
	"os"
)

var (
	httpEndpoint = os.Getenv("HTTP_PROXY")
)

func EthProxyHandler(w http.ResponseWriter, r *http.Request) {
	_, _ = io.WriteString(w, "Hello from handler!\n")
}
