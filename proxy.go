package proxy

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
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
	log.Print(url)
	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		return
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return err
	}
	_, err = GetSignature(req, bytes.NewReader(body))
	if err != nil {
		return err
	}
	fmt.Println(req.Header)

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	_, _ = io.Copy(w, resp.Body)
	return nil
}

func GetSignature(r *http.Request, body io.ReadSeeker) (http.Header, error) {
	serviceName := "managedblockchain"
	config := defaults.Config()
	creds := defaults.CredChain(config, defaults.Handlers())
	signer := v4.NewSigner(creds)
	_, err := signer.Sign(r, body, serviceName, *config.Region, time.Now())
	if err != nil {
		return nil, err
	}
	//	fmt.Println(r.Header)

	return r.Header, nil
}
