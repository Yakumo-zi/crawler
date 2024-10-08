package collect

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"time"

	"crawler/proxy"

	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/transform"
)

type Fetcher interface {
	Get(req *Request) ([]byte, error)
}
type BaseFetch struct {
}

func (BaseFetch) Get(req *Request) ([]byte, error) {
	resp, err := http.Get(req.URL)
	if err != nil {
		return nil, err
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

type BrowserFetch struct {
	Proxy   proxy.ProxyFunc
	Timeout time.Duration
}

func (b BrowserFetch) Get(request *Request) ([]byte, error) {
	client := &http.Client{
		Timeout: b.Timeout,
	}
	if b.Proxy != nil {
		transport := http.DefaultTransport.(*http.Transport)
		transport.Proxy = b.Proxy
		client.Transport = transport
	}
	req, err := http.NewRequest("GET", request.URL, nil)
	if err != nil {
		return nil, err
	}
	if len(request.UA) > 0 {
		req.Header.Set("User-Agent", request.UA)
	}
	if len(request.Cookie) > 0 {
		req.Header.Set("Cookie", request.Cookie)
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error status code:%v\n", resp.StatusCode)
	}

	bodyReader := bufio.NewReader(resp.Body)
	e := DeterminEncoding(bodyReader)
	utf8Reader := transform.NewReader(bodyReader, e.NewEncoder())
	return io.ReadAll(utf8Reader)
}

func DeterminEncoding(r *bufio.Reader) encoding.Encoding {
	bytes, err := r.Peek(1024)
	if err != nil {
		fmt.Printf("fetch error:%v\n", err)
		return unicode.UTF8
	}
	e, _, _ := charset.DetermineEncoding(bytes, "")
	return e
}
