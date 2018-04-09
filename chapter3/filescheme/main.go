package main

import (
	"log"
	"net/http"
	"net/http/httputil"
)

// このプログラムはこのファイルが存在する場所で実行する必要がある。
func main() {
	transport := &http.Transport{}
	transport.RegisterProtocol("file", http.NewFileTransport(http.Dir(".")))
	//transport.RegisterProtocol("file", http.NewFileTransport(http.Dir("/")))

	client := http.Client{
		Transport: transport,
	}
	resp, err := client.Get("file://./main.go")
	//resp, err := client.Get("file:///Users/ssmylh/works/go/rwh/chapter3/filescheme/main.go")
	if err != nil {
		panic(err)
	}

	dump, err := httputil.DumpResponse(resp, true)
	if err != nil {
		panic(err)
	}
	log.Println(string(dump))
}
