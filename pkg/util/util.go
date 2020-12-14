package util

import (
	"github.com/gocolly/colly"
	"golang.org/x/net/html"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
)

var HtmlVersionMap = map[string]string{
	"html" : "HTML5",
	"HTML PUBLIC \"-//W3C//DTD HTML 4.01//EN\" \"http://www.w3.org/TR/html4/strict.dtd\"":"HTML 4.01 Strict",
	"HTML PUBLIC \"-//W3C//DTD HTML 4.01 Transitional//EN\" \"http://www.w3.org/TR/html4/loose.dtd\"":"HTML 4.01 Transitional",
	"HTML PUBLIC \"-//W3C//DTD HTML 4.01 Frameset//EN\" \"http://www.w3.org/TR/html4/frameset.dtd\"" : "HTML 4.01 Frameset",
	"html PUBLIC \"-//W3C//DTD XHTML 1.0 Strict//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-strict.dtd\"" : "XHTML 1.0 Strict",
	"html PUBLIC \"-//W3C//DTD XHTML 1.0 Transitional//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-transitional.dtd\"":"XHTML 1.0 Transitional",
	"html PUBLIC \"-//W3C//DTD XHTML 1.0 Frameset//EN\" \"http://www.w3.org/TR/xhtml1/DTD/xhtml1-frameset.dtd\"":"XHTML 1.0 Frameset",
	"html PUBLIC \"-//W3C//DTD XHTML 1.1//EN\" \"http://www.w3.org/TR/xhtml11/DTD/xhtml11.dtd\"":"XHTML 1.1",
}

func GetHref(t html.Token) (ok bool, href string) {
	for _, a := range t.Attr {
		if a.Key == "href"{
			if strings.Index( a.Val,"javascript:")==0{
				ok = false
				return
			}
			href = a.Val
			ok = true
		}
	}
	return
}

func IsLoginForm(url string, ch chan  bool){

		isPassword := false
		c := colly.NewCollector(
			colly.MaxDepth(1),
			colly.UserAgent("abc"),
		)

		c.OnHTML("input", func(e *colly.HTMLElement) {
			link := e.Attr("type")
			if link=="password"{
				isPassword = true
			}
		})

		c.Visit(url)
		ch <- isPassword

}

func Fetch2(url string,ch chan<- string, wg *sync.WaitGroup,client *http.Client){
	defer wg.Done()
	log.Println("Visiting: "+url)

	request, err := http.NewRequest("HEAD", url, nil)
	if err != nil {
		ch <- strconv.Itoa(-1) // send to channel ch
		return
	}
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")
	request.Header.Set("Dnt", "1")
	request.Header.Set("Upgrade-Insecure-Requests", "1")

	resp, err := client.Do(request)

	if err != nil {
		ch <- strconv.Itoa(-1)
		return
	}
	defer resp.Body.Close()
	ch <- strconv.Itoa(200)
	log.Println("Visited: "+url)
}