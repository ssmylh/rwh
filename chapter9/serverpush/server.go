package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

var image []byte

func init() {
	var err error
	image, err = ioutil.ReadFile("chapter9/serverpush/puzzle.jpg")
	if err != nil {
		panic(err)
	}
}

func handlerHTML(w http.ResponseWriter, r *http.Request) {
	pusher, ok := w.(http.Pusher)
	if ok {
		pusher.Push("/image", nil)
	}
	w.Header().Add("Content-Type", "text/html")
	fmt.Fprintf(w, `<html><body><img src="/image"></body></html>`)
}

func handlerImage(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "image/jpeg")
	w.Write(image)
}

// 画像と証明書等を参照しているので以下で実行
// go run chapter9/serverpush/server.go
func main() {
	http.HandleFunc("/", handlerHTML)
	http.HandleFunc("/image", handlerImage)
	fmt.Println("start http listening :18443")
	err := http.ListenAndServeTLS(":18443", "server.crt", "server.key", nil)
	fmt.Println(err)
}
