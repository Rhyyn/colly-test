package scrapers

import (
	"time"

	"github.com/gocolly/colly"
)

var (
	BaseURL   = "https://www.wakfu.com"
	userAgent = "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/109.0.0.0 Safari/537.36"
)

func GetNewCollector() *colly.Collector {
	c := colly.NewCollector()

	c.Limit(&colly.LimitRule{
		Delay: 2 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

	return c
}
