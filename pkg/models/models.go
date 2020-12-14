package models


type CrawlResponse struct {
	Title string
	HtmlVersion string
	ExternalLinks int
	InternalLinks int
	InaccessibleLinks int
	Headings map[string]int
	IsLoginFormPresent bool
}
