package scrapers

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	"github.com/Rhyyn/wakfukiscraper/structs"
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

func contains(slice []int, val int) bool {
	for _, item := range slice {
		if item == val {
			return true
		}
	}
	return false
}

func AppendIDsToFile(newIds []int, selectedType string) {
	filePath := "./DATA/STATIC/ScrapedIds/" + selectedType + ".json"

	var data []int

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		if err := json.Unmarshal(file, &data); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	} else if os.IsNotExist(err) {
		data = []int{}
	} else {
		fmt.Println("Error checking file status:", err)
		return
	}

	for _, id := range newIds {
		if !contains(data, id) {
			data = append(data, id)
		}
	}

	jsonData, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	if err := os.WriteFile(filePath, jsonData, 0o644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Data appended successfully.")
}

func AppendItemsToFile(newItems map[int]structs.Item, selectedType string) {
	filePath := "./DATA/STATIC/ScrapedData/" + selectedType + ".json"

	var Items map[int]structs.Item

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		if err := json.Unmarshal(file, &Items); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	} else if os.IsNotExist(err) {
		Items = map[int]structs.Item{}
	} else {
		fmt.Println("Error checking file status:", err)
		return
	}

	for id, newItem := range newItems {
		Items[id] = newItem
	}

	jsonData, err := json.MarshalIndent(Items, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	if err := os.WriteFile(filePath, jsonData, 0o644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Items appended successfully.")
}
