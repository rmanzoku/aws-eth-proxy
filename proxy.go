package proxy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"time"

	"github.com/aws/aws-sdk-go/aws/defaults"
	v4 "github.com/aws/aws-sdk-go/aws/signer/v4"
)

type Client struct {
	Endpoint string
	Logger   *log.Logger
}

func NewClient(endpoint string) *Client {
	return &Client{
		Endpoint: endpoint,
		Logger:   log.New(os.Stdout, "", 0),
	}
}

func (c *Client) Handler(w http.ResponseWriter, r *http.Request) {
	err := c.Proxy(w, r)
	if err != nil {
		http.Error(w, fmt.Sprintf("...: %v", err), 500)
	}
}

type logCtx struct {
	TimeStamp int64  `json:"ts"`
	Method    string `json:"method"`
	Status    int    `json:"status"`
	Code      int64  `json:"code"`
	Message   string `json:"message"`
	Error     error  `json:"error"`
}

func (c *Client) log(lc *logCtx) {
	b, _ := json.Marshal(lc)
	c.Logger.Printf(string(b))
}

func (c *Client) Proxy(w http.ResponseWriter, r *http.Request) (err error) {
	lc := new(logCtx)
	lc.TimeStamp = time.Now().Unix()
	defer c.log(lc)

	client := new(http.Client)
	url := c.Endpoint + r.URL.Path
	req, err := http.NewRequestWithContext(r.Context(), r.Method, url, r.Body)
	if err != nil {
		lc.Error = err
		return
	}
	inBodyBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		lc.Error = err
		return
	}
	in := new(jsonRPCInput)
	err = json.Unmarshal(inBodyBytes, in)
	if err != nil {
		lc.Error = err
		return
	}
	lc.Method = in.Method

	req, err = Sign(req, inBodyBytes)
	if err != nil {
		lc.Error = err
		return
	}

	resp, err := client.Do(req)
	if err != nil {
		lc.Error = err
		return
	}
	defer resp.Body.Close()
	lc.Status = resp.StatusCode

	outBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		lc.Error = err
		return
	}

	out := new(jsonRPCOutput)
	err = json.Unmarshal(outBodyBytes, out)
	if err != nil {
		lc.Error = err
		return
	}
	lc.Code = out.Error.Code
	lc.Message = out.Error.Message

	_, err = w.Write(outBodyBytes)
	lc.Error = err
	return err
}

func Sign(r *http.Request, body []byte) (*http.Request, error) {
	config := defaults.Config()
	creds := defaults.CredChain(config, defaults.Handlers())
	signer := v4.NewSigner(creds)
	_, err := signer.Sign(r, bytes.NewReader(body), "managedblockchain", *config.Region, time.Now())
	return r, err
}

type jsonRPCInput struct {
	Jsonrpc string          `json:"jsonrpc"`
	ID      int64           `json:"id"`
	Method  string          `json:"method"`
	Params  json.RawMessage `json:"params"`
}

type jsonRPCOutput struct {
	Jsonrpc string       `json:"jsonrpc"`
	ID      int64        `json:"id"`
	Error   jsonRPCError `json:"error"`
}

type jsonRPCError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}
