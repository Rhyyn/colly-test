package scrapers

import (
	"fmt"

	// utils "github.com/Rhyyn/wakfukiscraper/utils"
	"github.com/gocolly/colly"
)

func StartCollySraper() {
	c := colly.NewCollector()

	c.OnHTML("a[href]", func(e *colly.HTMLElement) {
		e.Request.Visit(e.Attr("href"))
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting", r.URL)
	})
	c.Visit("https://wakfuki.com")
}

func CrawlSingleAccesoryType(type_id int) {
	fmt.Printf("Crawling %d...\n", type_id)
}

func CrawlAllAccessories() error {
	return nil
}
