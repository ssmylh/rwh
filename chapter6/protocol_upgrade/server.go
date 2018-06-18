package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"time"
)

func handlerUpgrade(w http.ResponseWriter, r *http.Request) {
	if r.Header.Get("Connection") != "Upgrade" || r.Header.Get("Upgrade") != "MyProtocol" {
		w.WriteHeader(400)
		return
	}

	fmt.Println("Upgrade to MyProtocol")

	hj, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Webserver doesn't support hijacking", http.StatusInternalServerError)
		return
	}
	conn, bufrw, err := hj.Hijack()
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	response := http.Response{
		StatusCode: 101,
		Header:     make(http.Header),
	}
	response.Header.Set("Upgrade", "MyProtocol")
	response.Header.Set("Connection", "Upgrade")
	response.Write(conn)

	for i := 1; i <= 10; i++ {
		fmt.Fprintf(bufrw, "%d\n", i)
		fmt.Println("->", i)
		bufrw.Flush()
		recv, err := bufrw.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		fmt.Printf("<- %s", string(recv))
		time.Sleep(500 * time.Millisecond)
	}
}

func main() {
	var httpServerr http.Server
	http.HandleFunc("/upgrade", handlerUpgrade)

	log.Println("start http listening: 18888")
	httpServerr.Addr = ":18888"
	log.Println(httpServerr.ListenAndServe())
}
