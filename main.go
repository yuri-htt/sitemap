package main

import (
	"flag"
	"fmt"
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
	flag.Parse()

	fmt.Println(*urlFlag)
	resp, err := http.Get(*urlFlag)
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
		fmt.Println("確認2", baseUrl.Scheme) // https
		fmt.Println("確認3", baseUrl.Host) // gophercises.com
	*/
	base := baseUrl.String()
	// fmt.Println("Request URL:", reqUrl.String())
	// fmt.Println("Base URL:", base)

	links, _ := link.Parse(resp.Body)
	var hrefs []string
	for _, l := range links {
		switch {
		case strings.HasPrefix(l.Href, "/"):
			hrefs = append(hrefs, base+l.Href)
		case strings.HasPrefix(l.Href, "http"):
			hrefs = append(hrefs, l.Href)
		}
	}
	for _, href := range hrefs {
		fmt.Println(href)
	}
}
