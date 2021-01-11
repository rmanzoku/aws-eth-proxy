package main

import (
	"log"
	"net/http"
	"os"

	proxy "github.com/rmanzoku/aws-eth-proxy"
)

var (
	port   = "9000"
	target = ""
)

func init() {
	if os.Getenv("PORT") != "" {
		port = os.Getenv("PORT")
	}
	target = os.Getenv("HTTP_ENDPOINT")
}

func run() (err error) {
	awsnode := proxy.NewClient(target, false)
	http.HandleFunc("/", awsnode.Handler)

	return http.ListenAndServe(":"+port, nil)
}

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
}
