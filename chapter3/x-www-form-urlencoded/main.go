package main

import (
	"log"
	"net/http"
	"net/url"
)

func main() {
	values := url.Values{
		"test": {"value 1 & ="},
	}

	// RFC1866のURLエンコードで送信される(スペースが`+`になる)。
	resp, err := http.PostForm("http://localhost:18888", values)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}
