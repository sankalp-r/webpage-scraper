package handlers

import (
	"fmt"
	"testing"
	"net/http"
    "net/http/httptest"
    "github.com/stretchr/testify/assert"
)




func TestCrawl(t *testing.T) {
	
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			myMockPage := `<!DOCTYPE html><html lang=en-US class=no-js><head><title>What does the Otter say?</title></head><body>
					<h1>Head1</h1>
					<p><a href="`+r.URL.String()+`">Visit W3Schools.com!</a></p>
					</body>
					</html>`
			w.Header().Set("Content-Type", "text/html")

			fmt.Fprintln(w, myMockPage)
	}))
	defer ts.Close()

	req, err := http.NewRequest("GET", "/crawl", nil)
	q := req.URL.Query()
	q.Add("url", ts.URL)
	req.URL.RawQuery = q.Encode()

	if err != nil {
		t.Fatal(err)
	}


	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Crawl)
	handler.ServeHTTP(rr, req)
	
	expected := `{"Title":"What does the Otter say?","HtmlVersion":"HTML5","ExternalLinks":0,"InternalLinks":1,"InaccessibleLinks":0,"Headings":{"h1":1,"h2":0,"h3":0,"h4":0,"h5":0,"h6":0},"IsLoginFormPresent":false}`
	assert.True(t, rr.Body.String() == expected)
}
