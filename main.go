package main

import (
	"fmt"
	"github.com/Yakumo-zi/crawler/collect"
	"github.com/Yakumo-zi/crawler/proxy"
	"time"
)

func main() {
	proxyURLs := []string{"http://127.0.0.1:7890"}
	p, err := proxy.RoundRobinProxySwitcher(proxyURLs...)
	if err != nil {
		fmt.Println("RoundRobinProxySwitcher failed")
	}
	url := "https://google.com"
	var f collect.Fetcher = collect.BrowserFetch{
		Timeout: 3000 * time.Millisecond,
		Proxy:   p,
	}

	body, err := f.Get(url)
	if err != nil {
		fmt.Printf("read content failed:%v\\n", err)
		return
	}
	fmt.Println(string(body))
}
