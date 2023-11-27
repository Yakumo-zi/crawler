package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

func main() {
	url := "https://www.thepaper.cn/"
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("fetch url error:%v\n", err)
		return
	}

	defer func(Body io.ReadCloser) {
		_ = Body.Close()
	}(resp.Body)

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%v\n", resp.StatusCode)
		return
	}
	content, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Printf("read content failed:%v\n", err)
		return
	}

	file, err := os.Create("index.html")

	if err != nil {
		fmt.Printf("create file failed:%v\n", err)
		return
	}

	defer func(file *os.File) {
		_ = file.Close()
	}(file)

	n, err := file.WriteString(string(content))
	if err != nil {
		fmt.Printf("write content to file failed:%v\n", err)
		return
	}

	fmt.Printf("write %d bytes to file %s \n", n, file.Name())
}
