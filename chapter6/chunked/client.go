package main

import (
	"bufio"
	"io"
	"log"
	"net"
	"net/http"
	"strconv"
	"time"
)

func main() {
	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}
	conn, err := dialer.Dial("tcp", "localhost:18888")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	request, err := http.NewRequest("GET", "http://localhost:18888/chunked", nil)
	err = request.Write(conn)
	if err != nil {
		panic(err)
	}

	reader := bufio.NewReader(conn)
	resp, err := http.ReadResponse(reader, request)
	if err != nil {
		panic(err)
	}
	if resp.TransferEncoding[0] != "chunked" {
		panic("wrong transfer encoding")
	}

	for {
		read, err := reader.ReadBytes('\n')
		if err == io.EOF {
			break
		}
		// 末尾のCRLFを除いたサイズを取得する( `server.go` の `\n` は除いていない)
		size, err := strconv.ParseInt(string(read[:len(read)-2]), 16, 64)
		if err != nil {
			panic(err)
		}
		if size == 0 {
			break
		}

		line := make([]byte, int(size))
		reader.Read(line)
		reader.Discard(2) // 末尾のCRLF分を読み捨てる
		log.Println(" ", string(line))
	}
}
