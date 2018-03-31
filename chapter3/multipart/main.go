package main

import (
	"bytes"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
)

// 画像がこのファイルと同じディレクトリに存在するのを仮定しているので、
// プログラムはこのファイルが存在するディレクトリで実行する。
func main() {
	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	writer.WriteField("name", "Puzzle")

	fileWriter, err := writer.CreateFormFile("thumbnail", "puzzle.jpg")
	if err != nil {
		panic(err)
	}
	readFile, err := os.Open("puzzle.jpg")
	if err != nil {
		panic(err)
	}
	defer readFile.Close()
	io.Copy(fileWriter, readFile)
	writer.Close()

	resp, err := http.Post("http://localhost:18888", writer.FormDataContentType(), &buffer)
	if err != nil {
		panic(err)
	}
	log.Println("Status:", resp.Status)
}
