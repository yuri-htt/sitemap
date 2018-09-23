package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	// "github.com/gophercises/link"
	"github.com/yuri-swift/link"
)

/*
	1, Get the webpage 5:41
	2. parse all the links on the page 15:02
	3. build proper urls with our links
	4, filter out any links w/ a diff domain 11:42
	5. Find all pages (BFS) 20:34
	6. print out XML 13:55
*/

func main() {
	urlFlag := flag.String("url", "https://gophercises.com/", "the url that you waant to build a sitemap for")
	maxDepth := flag.Int("depth", 2, "the maximum number of links deep to traverse")
	flag.Parse()

	pages := bfs(*urlFlag, *maxDepth)
	for _, page := range pages {
		fmt.Println(page)
	}
}

func bfs(urlStr string, maxDepth int) []string {
	// variable := make(map[key_type]value_type)
	// variable := map[key_type] value_type {}
	seen := make(map[string]struct{})
	var q map[string]struct{}
	nq := map[string]struct{}{
		urlStr: struct{}{},
	}
	for i := 0; i < maxDepth; i++ {
		fmt.Println("■■■")
		fmt.Println(q)
		fmt.Println(nq)
		q, nq = nq, make(map[string]struct{})
		fmt.Println(q)
		fmt.Println(nq)
		for url, _ := range q {
			if _, ok := seen[url]; ok {
				continue
			}
			seen[url] = struct{}{}
			for _, link := range get(url) {
				nq[link] = struct{}{}
			}
		}
	}
	ret := make([]string, 0, len(seen))
	for url, _ := range seen {
		ret = append(ret, url)
	}
	return ret
}

func get(urlStr string) []string {
	resp, err := http.Get(urlStr)
	if err != nil {
		panic(err)
	}
	// 遅延実行: エラーが発生した場合でもBodyCloseを行う
	defer resp.Body.Close()

	reqUrl := resp.Request.URL
	/*
		// アドレスでアクセスしたのに構造体に値を代入できるの？
		test := &url.URL{}
		fmt.Println("確認1", test.String()) // 空
	*/
	baseUrl := &url.URL{
		Scheme: reqUrl.Scheme,
		Host:   reqUrl.Host,
	}
	/*
		fmt.Println("スキーマ", baseUrl.Scheme) // https
		fmt.Println("ホスト", baseUrl.Host) // gophercises.com
	*/
	base := baseUrl.String()
	// fmt.Println("Request URL:", reqUrl.String())
	// fmt.Println("Base URL:", base)
	return filter(hrefs(resp.Body, base), withPrefix(base))
}

func hrefs(r io.Reader, base string) []string {
	links, _ := link.Parse(r)
	var ret []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			ret = append(ret, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			ret = append(ret, l.Href)
		}
	}
	return ret
}

// https://gophercises.com/ではじまるリンクでフィルタリング
func filter(links []string, keepFn func(string) bool) []string {
	var ret []string
	for _, link := range links {
		if keepFn(link) {
			ret = append(ret, link)
		}
	}

	return ret
}

func withPrefix(pfx string) func(string) bool {
	return func(link string) bool {
		return strings.HasPrefix(link, pfx)
	}
}
