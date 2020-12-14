# webpage-scraper

## Simple web scraper app

### How to build and run
 * Navigate to `webpage-scraper/cmd/crawler/main` folder of project run : go run main.go
 * After local-server starts, enter `http://localhost:10000/` in the browser to open webapp
 * Enter the Url you want to scrape in the input box like: `https://www.google.com/` (Please enter full Url with http or https)
 * Click Go button to fetch the result.
 * Alternatively, you can also make an API call: `http://localhost:10000/crawl?url=https://www.google.com/`
 * Sample API reponse looks like : 
 	```json
   {
   "Title":"W3Schools Online Web Tutorials",
   "HtmlVersion":"HTML5",
   "ExternalLinks":8,
   "InternalLinks":150,
   "InaccessibleLinks":0,
   "Headings":{
      "h1":6,
      "h2":12,
      "h3":31,
      "h4":18,
      "h5":0,
      "h6":0
   },
   "IsLoginFormPresent":true
  }
```
 * Note: app has a timeout of 10s. If the response doesnt come back within 10s, you will get error.

### Thought process
 * For scraping the webpage, golang libraries have been used.
 * To find the inaccessible links, "HEAD" call is made to each link.
 * We might have many links in the webpage, so goroutines are used for making concurrent calls.
 * Also timeout configurations have been used for both server and httpclient to prevent resource throttling.
