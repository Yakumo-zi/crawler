package proxy

import (
	"errors"
	"net/http"
	"net/url"
	"sync/atomic"
)

type ProxyFunc func(*http.Request) (*url.URL, error)

func RoundedRobinProxySwitcher(ProxyURLs ...string) (ProxyFunc, error) {
	if len(ProxyURLs) < 1 {
		return nil, errors.New("proxy URL list is empty")
	}

	urls := make([]*url.URL, len(ProxyURLs))
	for i, u := range ProxyURLs {
		parsed, err := url.Parse(u)
		if err != nil {
			return nil, err
		}
		urls[i] = parsed
	}
	return (&roundRobinSwitcher{urls, 0}).GetProxy, nil
}

type roundRobinSwitcher struct {
	proxyURLs []*url.URL
	index     uint32
}

func (r *roundRobinSwitcher) GetProxy(req *http.Request) (*url.URL, error) {
	index := atomic.AddUint32(&r.index, 1) - 1
	u := r.proxyURLs[index%uint32(len(r.proxyURLs))]
	return u, nil
}
