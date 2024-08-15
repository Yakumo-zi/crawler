package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"net/http"

	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

func main() {

	body, err := Fetch("https://www.thepaper.cn/channel_25950")
	if err != nil {
		panic(err)
	}
	doc, err := goquery.NewDocumentFromReader(bytes.NewReader(body))
	if err != nil {
		panic(err)
	}
	doc.Find(".ant-card-body a").Each(func(i int, s *goquery.Selection) {
		title := s.Find("h2").Text()
		link, ok := s.Attr("href")
		if !ok {
			return
		}
		if len(title) == 0 {
			return
		}
		fmt.Printf("No.%d\n", i)
		fmt.Printf("\tTitle:%s\n", title)
		fmt.Printf("\tLink:%s\n", link)
	})
}

func Fetch(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%v", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewEncoder())
	return io.ReadAll(utf8Reader)

}

func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		fmt.Printf("fetch error:%v", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
