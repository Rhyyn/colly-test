package scrapers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/Rhyyn/wakfukiscraper/structs"
	"github.com/Rhyyn/wakfukiscraper/utils"
	"github.com/gocolly/colly"
)

func ScrapItemDetails(url string, Item *structs.Item, Recipes map[int]structs.Recipe, ParamsStatsProperties structs.ParamsStatsProperties) {
	c := GetNewCollector()
	c.OnHTML(".ak-container.ak-panel-stack.ak-glue", func(h *colly.HTMLElement) {
		Lang := utils.GetLangFromURL(h.Request.URL.String())
		GetTitle(h, Item, Lang)
		// fmt.Println("Got Title")
		GetTypeID(h, Item, Lang)
		// fmt.Println("Got TypeId")
		GetRarity(h, Item, Lang)
		// fmt.Println("Got Rarity")
		GetLevel(h, Item, Lang)
		// fmt.Println("Got Level")
		GetStats(h, Item, Lang, ParamsStatsProperties)
		// fmt.Println("Got Stats")
		GetDroprates(h, Item, Lang)
		// fmt.Println("Got Droprates")
		GetRecipes(h, Item, Recipes, Lang)
		// fmt.Println("Got Recipes")
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("ScrapItemDetails visiting:\n", r.URL)
	})
	c.Visit(url)
}

func GetItemURLArg(url string) (string, error) {
	// /fr/mmorpg/encyclopedie/armures/29086-bottes-vaal-enthia
	parts := strings.Split(url, "/")
	if len(parts) >= 5 {
		return parts[5], nil
	} else {
		return "", errors.New("URL format error: invalid number of parts")
	}
}

func ScrapItems(ScrapingParameters structs.ScrapingParameters) {
	urlSuffix := "&" + "type_1%5B%5D=" + strconv.Itoa(ScrapingParameters.SelectedId)
	AllPositivesStats := utils.HandleStatsProperties(utils.ReadFile(utils.OpenFile("all_positives_stats.json")))
	AllNegativesStats := utils.HandleStatsProperties(utils.ReadFile(utils.OpenFile("all_negatives_stats.json")))

	var IDsList []int
	Items := make(map[int]structs.Item)
	Recipes := make(map[int]structs.Recipe)

	fmt.Printf("ScrapItems called for id %d with maxPage %d\n", ScrapingParameters.SelectedId, ScrapingParameters.MaxPage)
	c := GetNewCollector()

	// ON EVERY TR IN THE TABLE
	c.OnHTML(".ak-table.ak-responsivetable tbody tr", func(h *colly.HTMLElement) {
		// extract each item href from each td
		href, exists := h.DOM.Find("td").Eq(1).Find("a[href]").Attr("href")
		if !exists {
			fmt.Printf("NO TD FOUND FOR %s\n", href)
		}
		itemArgName, err := GetItemURLArg(href)
		if err != nil {
			fmt.Println(err)
			fmt.Printf("error getting item url arg")
			os.Exit(0)
		}

		frenchURL := ScrapingParameters.ItemUrl["fr"] + "/" + itemArgName
		englishURL := ScrapingParameters.ItemUrl["en"] + "/" + itemArgName

		var Item structs.Item
		ParamsStatsProperties := structs.ParamsStatsProperties{AllPositivesStats: AllPositivesStats, AllNegativesStats: AllNegativesStats}

		Item.ID, err = utils.GetItemIDFromString(itemArgName)
		if err != nil {
			fmt.Println("error while getting the item ID", err)
			os.Exit(0)
		}
		// TODO : Refactor check for update
		// so we dont need to scrap what we already scraped

		IDsList = append(IDsList, Item.ID)
		// Scrap both FR/EN version of the item
		ScrapItemDetails(frenchURL, &Item, Recipes, ParamsStatsProperties)
		ScrapItemDetails(englishURL, &Item, Recipes, ParamsStatsProperties)

		Items[Item.ID] = Item
		// Useless pretty print for debug
		// PrettyItem, err := json.MarshalIndent(Item, "", "    ")
		// if err != nil {
		// 	fmt.Println("Error marshaling item:", err)
		// 	return
		// }
		// fmt.Println("Item after scraping:\n", string(PrettyItem))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("ScrapItems visiting:\n", r.URL)
	})

	if !ScrapingParameters.SingleItemMode {
		for i := 1; i < ScrapingParameters.MaxPage; i++ {
			if i != 1 && len(IDsList) > 0 {
				AppendIDsToFile(IDsList, ScrapingParameters.SelectedType)
				AppendItemsToFile(Items, ScrapingParameters.SelectedType)

			}
			IDsList = []int{}
			Items = make(map[int]structs.Item)
			c.Visit(ScrapingParameters.IndexUrl["fr"] + strconv.Itoa(i) + urlSuffix)
			fmt.Printf("setting page to %d\n", i)
		}
	} else {
		var Item structs.Item
		var Recipes map[int]structs.Recipe

		ParamsStatsProperties := structs.ParamsStatsProperties{AllPositivesStats: AllPositivesStats, AllNegativesStats: AllNegativesStats}
		ScrapItemDetails(ScrapingParameters.SingleItemURL["Fr"], &Item, Recipes, ParamsStatsProperties)
		ScrapItemDetails(ScrapingParameters.SingleItemURL["En"], &Item, Recipes, ParamsStatsProperties)

		PrettyItem, err := json.MarshalIndent(Item, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling item:", err)
			return
		}
		fmt.Println("Item after scraping:\n", string(PrettyItem))
	}
}
