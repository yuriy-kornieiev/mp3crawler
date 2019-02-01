package main

import (
	"crypto/tls"
	"fmt"
	"golang.org/x/net/html"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Crawler struct {
}

func (c Crawler) Start(startPage string) {

	//re1 := regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
	//
	//html, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("%q\n", re1.FindAllString(string(html), -1))

}

func (c Crawler) enqueue(uri string, queue chan string) {
	fmt.Println("fetching", uri)
	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	resp, err := client.Get(uri)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	links := c.All(resp.Body)

	for _, link := range links {
		absolute := c.fixUrl(link, uri)
		if uri != "" {
			go func() { queue <- absolute }() // We asynchronously enqueue what we've found
		}
	}
}

func (c Crawler) fixUrl(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return ""
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return ""
	}
	uri = baseUrl.ResolveReference(uri)
	return uri.String()
}

func (c Crawler) All(httpBody io.Reader) []string {
	links := []string{}
	col := []string{}
	page := html.NewTokenizer(httpBody)
	for {
		tokenType := page.Next()
		if tokenType == html.ErrorToken {
			return links
		}
		token := page.Token()
		if tokenType == html.StartTagToken && token.DataAtom.String() == "a" {
			for _, attr := range token.Attr {
				if attr.Key == "href" {
					tl := c.trimHash(attr.Val)
					col = append(col, tl)
					c.resolv(&links, col)
				}
			}
		}
	}
}

// trimHash slices a hash # from the link
func (c Crawler) trimHash(l string) string {
	if strings.Contains(l, "#") {
		var index int
		for n, str := range l {
			if strconv.QuoteRune(str) == "'#'" {
				index = n
				break
			}
		}
		return l[:index]
	}
	return l
}

// resolv adds links to the link slice and insures that there is no repetition
// in our collection.
func (c Crawler) resolv(sl *[]string, ml []string) {
	for _, str := range ml {
		if c.check(*sl, str) == false {
			*sl = append(*sl, str)
		}
	}
}

// check looks to see if a url exits in the slice.
func (c Crawler) check(sl []string, s string) bool {
	var check bool
	for _, str := range sl {
		if str == s {
			check = true
			break
		}
	}
	return check
}
