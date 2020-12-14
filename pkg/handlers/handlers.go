package handlers

import (
	"crypto/tls"
	"encoding/json"
	"fmt"
	"golang.org/x/net/html"
	"h24/pkg/models"
	"h24/pkg/util"
	"html/template"
	"log"
	"net"
	"net/http"
	url2 "net/url"
	"strings"
	"sync"
	"time"
	//"io/ioutil"
)



var client *http.Client

func init(){
	tr := &http.Transport{
		Proxy: http.ProxyFromEnvironment,
		DialContext: (&net.Dialer{
			Timeout:   30 * time.Second,
			KeepAlive: 30 * time.Second,
		}).DialContext,
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true,
		},
		MaxIdleConnsPerHost: 1024,
		MaxIdleConns: 1024,
		TLSHandshakeTimeout: 10 * time.Second,
		//ResponseHeaderTimeout: 5 * time.Second,
		IdleConnTimeout: 90 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,

	}
	client = &http.Client{
		//Timeout: 10 * time.Second,
		Transport: tr,
	}
}

func Crawl(w http.ResponseWriter, r *http.Request){
	err:=r.ParseForm()
	if err!=nil{
		log.Println(err)
		w.WriteHeader(500)
		fmt.Fprintf(w,err.Error())
		return
	}
	queryUrl := r.Form["url"][0]

	request, err := http.NewRequest("GET", queryUrl, nil)
	request.Header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_4) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/83.0.4103.97 Safari/537.36")
	request.Header.Set("Dnt", "1")
	request.Header.Set("Upgrade-Insecure-Requests", "1")

	resposne, err := client.Do(request)
	if err!=nil{
		log.Println(err)
		w.WriteHeader(500)
		fmt.Fprintf(w,err.Error())
		return
	}

	externalLinkSet := make(map[string]bool)
	internalLinkSet := make(map[string]bool)
	headingCount := map[string]int{
		"h1":  0,
		"h2":   0,
		"h3": 0,
		"h4": 0,
		"h5": 0,
		"h6":0,
	}
	ch := make(chan string)
	var wg sync.WaitGroup
	parsedQueryUrl,err:=url2.Parse(queryUrl)

	if err!=nil{
		w.WriteHeader(500)
		fmt.Fprintf(w,err.Error())
		return
	}
	parsedQueryUrl.RawQuery = ""
	
	defer resposne.Body.Close()
	z := html.NewTokenizer(resposne.Body)
	crawlResponse := models.CrawlResponse{IsLoginFormPresent:false}

	for {
		tt := z.Next()
		testt := z.Token()
		//fmt.Println(testt.Data)
		switch {
		case tt == html.ErrorToken:

			inaccessible := 0
			isLogin := make(chan bool)

			go util.IsLoginForm(queryUrl,isLogin)

			go func() {
				wg.Wait()
				close(ch)
			}()

			for val := range ch{
				if "200" != val && val != "-1"{
					inaccessible++
				}
			}

			crawlResponse.ExternalLinks = len(externalLinkSet)
			crawlResponse.InternalLinks = len(internalLinkSet)
			crawlResponse.InaccessibleLinks = inaccessible
			crawlResponse.Headings = headingCount
			crawlResponse.IsLoginFormPresent = <-isLogin
			jsonReponse,err := json.Marshal(crawlResponse)

			if err!=nil{
				w.WriteHeader(500)
				fmt.Fprintf(w,err.Error())
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.Write(jsonReponse)
			fmt.Println("Returned")
			return

		case tt == html.StartTagToken:
			isAnchor := testt.Data == "a"

			if !isAnchor {
				if testt.Data == "title"{
					tt = z.Next()
					if tt == html.TextToken{
						testt = z.Token()

						crawlResponse.Title = testt.Data
					}
				} else {
							if _,exist := headingCount[strings.ToLower(testt.Data)]; exist{
								headingCount[strings.ToLower(testt.Data)]++
							}

				}
			} else{
				ok, url := util.GetHref(testt)
				//fmt.Println(url)
				if !ok {
					continue
				}

				parsedUrl,err:=url2.Parse(url)
				if err!=nil{
					log.Println(err.Error()+":"+url)
				}

				if parsedUrl.IsAbs(){
					temp := parsedUrl.String()
					if parsedQueryUrl.Host!=parsedUrl.Host{

						if _,exist := externalLinkSet[temp]; !exist{
							externalLinkSet[temp]=true
							wg.Add(1)
							go util.Fetch2(temp,ch,&wg,client)
						}

					} else{

						if _,exist := internalLinkSet[temp]; !exist{
							internalLinkSet[temp]=true
							wg.Add(1)
							go util.Fetch2(temp,ch, &wg,client)
						}
					}
				} else{
					if url[0]!='/'{
						url=string('/')+url
					}
					if parsedQueryUrl.Path == "" || parsedQueryUrl.Path[len(parsedQueryUrl.Path)-1]!='/'{
						temp := parsedQueryUrl.String()+url
						if _, exist := internalLinkSet[temp]; !exist{
							internalLinkSet[temp] = true
							wg.Add(1)
							go util.Fetch2(temp,ch, &wg,client)
						}

					} else{
						temp := parsedQueryUrl.String()+url[1:]
						if _, exist := internalLinkSet[temp]; !exist{
							internalLinkSet[temp] = true
							wg.Add(1)
							go util.Fetch2(temp,ch, &wg,client)
						}
					}
				}
			}
		case tt == html.DoctypeToken:
			if _,exist := util.HtmlVersionMap[testt.Data]; exist{
				crawlResponse.HtmlVersion = util.HtmlVersionMap[testt.Data]
			} else{
				crawlResponse.HtmlVersion = "Unknown"

			}


		}
	}

}

func HomePage(w http.ResponseWriter, r *http.Request) {
	t, err := template.ParseFiles("index.html")
	if err != nil { // if there is an error
		log.Print("template parsing error: ", err)
	}
	err = t.Execute(w,nil)
	if err != nil {
		log.Print("template executing error: ", err)
	}
}
