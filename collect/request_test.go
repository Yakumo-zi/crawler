package collect

import (
	"bytes"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/PuerkitoBio/goquery"
)

func TestRequest(t *testing.T) {
	var worklist []*Request
	for i := 0; i <= 100; i += 29 {
		str := fmt.Sprintf("https://www.douban.com/group/szsh/discussion?start=%d", i)
		worklist = append(worklist, &Request{
			URL: str,
			ParseFunc: func(b []byte) ParseResult {
				result := ParseResult{}
				doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
				if err != nil {
					t.Fatalf("%v", err)
				}
				doc.Find("td.title a").Each(func(i int, s *goquery.Selection) {
					link, ok := s.Attr("href")
					if !ok {
						return
					}
					result.Requersrts = append(result.Requersrts, &Request{
						URL: link,
						ParseFunc: func(b []byte) ParseResult {
							r := ParseResult{}
							exists := false
							doc, err := goquery.NewDocumentFromReader(bytes.NewReader(b))
							if err != nil {
								t.Fatalf("%v", err)
							}
							doc.Find(".topic-content p").Each(func(i int, s *goquery.Selection) {
								topic := s.Text()
								if strings.Contains(topic, "阳台") && !exists {
									r.Items = append(r.Items, link)
									exists = true
								}
							})
							return r
						},
					})
				})
				return result
			},
		})
	}
	f := &BrowserFetch{
		Timeout: time.Millisecond * 3000,
	}
	for _, v := range worklist {
		v.Cookie = `bid=OhPkfNCEjFA; __utmz=30149280.1723726737.15.14.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); viewed="1007305_35556889_36573544_35233448_10555435_1148282_6859720_20436488_27016236_36683615"; _pk_id.100001.8cb4=ecf477af6c65b636.1723737783.; __yadk_uid=bUdabnrgmC8Xre0IFs5GWKPDj9jAVD2B; douban-fav-remind=1; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.215321639.1719065628.1723737785.1723787553.18; __utmc=30149280; __utmt=1; __utmb=30149280.60.5.1723790172127`
		v.UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"
		b, err := f.Get(v)
		if err != nil {
			t.Fatal(err)
		}
		fmt.Printf("%s\n", v.URL)
		result := v.ParseFunc(b)
		for _, req := range result.Requersrts {
			req.Cookie = `bid=OhPkfNCEjFA; __utmz=30149280.1723726737.15.14.utmcsr=google|utmccn=(organic)|utmcmd=organic|utmctr=(not%20provided); viewed="1007305_35556889_36573544_35233448_10555435_1148282_6859720_20436488_27016236_36683615"; _pk_id.100001.8cb4=ecf477af6c65b636.1723737783.; __yadk_uid=bUdabnrgmC8Xre0IFs5GWKPDj9jAVD2B; douban-fav-remind=1; _pk_ses.100001.8cb4=1; ap_v=0,6.0; __utma=30149280.215321639.1719065628.1723737785.1723787553.18; __utmc=30149280; __utmt=1; __utmb=30149280.60.5.1723790172127`
			req.UA = "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/127.0.0.0 Safari/537.36"
			b, err := f.Get(req)
			time.Sleep(time.Second * 1)
			if err != nil {
				t.Fatal(err)
			}
			res := req.ParseFunc(b)
			if len(res.Items) > 0 {
				fmt.Printf("%v\n", res.Items)
			}
		}

	}

}
