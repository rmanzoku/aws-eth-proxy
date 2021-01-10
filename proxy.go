package proxy

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
)

var (
	httpEndpoint = os.Getenv("HTTP_ENDPOINT")
)

func EthProxyHandler(w http.ResponseWriter, r *http.Request) {
	err := Proxy(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("...: %v", err), 500)
	}

}

func Proxy(w http.ResponseWriter, r *http.Request) (err error) {
	client := new(http.Client)
	url := httpEndpoint + r.URL.Path
	log.Print(url)
	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(w, resp.Body)
	return nil
}
