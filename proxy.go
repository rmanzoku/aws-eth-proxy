package proxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/defaults"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
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
	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	req, err = Sign(req, body)
	if err != nil {
		return err
	}

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(w, resp.Body)
	return nil
}

func Sign(r *http.Request, body []byte) (*http.Request, error) {
	config := defaults.Config()
	creds := defaults.CredChain(config, defaults.Handlers())
	signer := v4.NewSigner(creds)
	_, err := signer.Sign(r, bytes.NewReader(body), "managedblockchain", *config.Region, time.Now())
	return r, err
}
