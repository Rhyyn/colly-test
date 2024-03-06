package scrapers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/Rhyyn/wakfukiscraper/utils"
	"github.com/gocolly/colly"
)

type IndexOptions struct {
	Title     string
	Index_url string
	ID        []int
}

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

func CountItemsInPage(url string) int {
	// ex /fr/mmorpg/encyclopedie/armures?type_1[0]=119&page=33
	c := GetNewCollector()

	nOfItems := 0
	c.OnHTML(".ak-table.ak-responsivetable tbody tr", func(e *colly.HTMLElement) {
		nOfItems++
	})

	c.OnRequest(func(r *colly.Request) { fmt.Println("Visiting :", r.URL) })
	// c.OnRequest(func(r *colly.Request) { fmt.Println("body :", r.Body) })

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.Visit(BaseURL + url)
	return nOfItems
}

func UpdateMaxItemsAndPages(IndexOptions IndexOptions) {
	startingUrl := IndexOptions.Index_url
	var ID int
	// If we selected a subcategory we need to filter by its type
	if IndexOptions.ID != nil {
		ID = IndexOptions.ID[0]
		startingUrl = startingUrl + "1&" + "type_1%5B%5D=" + strconv.Itoa(IndexOptions.ID[0])
	}

	c := GetNewCollector()

	var maxItems int
	// This finds the menu for pages at the bottom of the page
	// and follows the last href and counts the number of items present in the last page
	c.OnHTML(".text-center.ak-pagination.hidden-xs .ak-pagination.pagination.ak-ajaxloader li:last-child a", func(e *colly.HTMLElement) {
		href := e.Attr("href")
		lastPageLength := CountItemsInPage(href)
		maxPages, err := strconv.Atoi(strings.Split(href, "page=")[1])
		if err != nil {
			fmt.Println(err)
		}
		numOfItemPerPage := 24
		maxItems = lastPageLength + ((maxPages - 1) * numOfItemPerPage)
		fmt.Printf("maxPages %d, maxItems %d\n", maxPages, maxItems)

		utils.EditItemsCats(utils.EditFileOptions{
			IsSubCat: true,
			ID:       ID,
			SubCat:   utils.SubCategory{Max_page: maxPages, MaxItems: maxItems},
		})
		fmt.Printf("max_items for %s is %d\n", IndexOptions.Title, maxItems)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("Visiting :", r.URL)
	})

	c.Visit(startingUrl)
}

func UpdateAllIDs() {
}
