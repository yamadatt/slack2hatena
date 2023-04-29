package main

import (
	"fmt"
	"regexp"

	"github.com/gocolly/colly"
)

func main() {
	c := colly.NewCollector(
		colly.URLFilters(
			regexp.MustCompile("https://ameblo.jp/drcynthia/entry-.+html"),
			regexp.MustCompile("https://ameblo.jp/drcynthia/theme-.+html"),
		),
	)

	// Find and visit all links
	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})

	//c.Visit("https://ameblo.jp/drcynthia")
	c.Visit("https://ameblo.jp/drcynthia/theme-10007146990.html")

}
