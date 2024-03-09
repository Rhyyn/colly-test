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
	c := colly.NewCollector(colly.AllowedDomains("wakfu.com", "www.wakfu.com", "account.ankama.com"))

	// colly.CacheDir("./wakfu_cache"))

	c.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		Parallelism: 1,
		Delay:       3 * time.Second,
	})

	c.OnRequest(func(r *colly.Request) {
		r.Headers.Set("User-Agent", userAgent)
	})

	return c
}
