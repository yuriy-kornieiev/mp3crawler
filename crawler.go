package main

import (
	"crypto/tls"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/url"
	"path"
	"regexp"
	"strconv"
	"strings"
	"time"
)

type Crawler struct {
}

var allowedExtensions = []string{"", ".html", ".htm", ".php", ".txt", ".json", ".xml", ".bml", ".cgi"}

var linkRegexp = regexp.MustCompile(`https?:\/\/(www\.)?[-a-zA-Z0-9@:%._\+~#=]{2,256}\.[a-z]{2,6}\b([-a-zA-Z0-9@:%_\+.~#?&//=]*)`)
var mp3Regexp = regexp.MustCompile(`((https?|ftp):((//)|(\\\\))+[\w\d:#@%/;$()~_?\+-=\\\.&]*\.mp3)`)

//	fmt.Printf("%q\n", mp3Regexp.FindAllString(string(html), -1))

func (c Crawler) enqueue(uri string, queue chan string) {

	url, err := url.Parse(uri)
	if err != nil {
		return
	}

	if !(Domain{}.GetDomain(*url)) {
		//log.Println("delay", "domain", uri)
		go func() { queue <- uri }() // We asynchronously enqueue what we've found
		return
	}

	for _, domain := range domainSkipList {
		if strings.Contains(url.Host, domain) {
			log.Println("skip", "domain", uri)
			return
		}
	}

	ext := path.Ext(url.Path)
	if len(ext) > 0 {
		if ok, _ := inArray(ext, allowedExtensions); !ok {
			log.Println("skip", ext, uri)
			return
		}
	}

	start := time.Now()

	transport := &http.Transport{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
	}
	client := http.Client{Transport: transport}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		log.Fatalln(err)
	}

	req.Header.Set("User-Agent", "MP3_Spider_Bot/0.1a")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	log.Println("ok", ext, resp.StatusCode, uri, time.Since(start))

	links := c.All(resp.Body)

	for _, link := range links {
		absolute := c.fixUrl(link, uri)
		if absolute != "" {
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
