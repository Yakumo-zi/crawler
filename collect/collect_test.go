package collect

import (
	"fmt"
	"testing"
	"time"

	"crawler/proxy"
)

func TestBrowserFetcher(t *testing.T) {
	proxyURLs := []string{"http://127.0.0.1:7890"}
	p, err := proxy.RoundedRobinProxySwitcher(proxyURLs...)
	if err != nil {
		t.Fatal("RoundRobinSwitcher failed")
	}
	url := "https://www.google.com"
	var f Fetcher = BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}
	body, err := f.Get(url)
	if err != nil {
		t.Fatalf(`read content failed:%v`, err)
	}
	fmt.Println(string(body))
}
