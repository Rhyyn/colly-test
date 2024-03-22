package scrapers

import (
	"fmt"
	"os"
	"strconv"

	"github.com/Rhyyn/wakfukiscraper/structs"
	"github.com/Rhyyn/wakfukiscraper/utils"
	"github.com/gocolly/colly"
)

func ScrapRessourceDetails(url string, Item *structs.Item) {
	c := GetNewCollector()
	c.OnHTML(".ak-container.ak-panel-stack.ak-glue", func(h *colly.HTMLElement) {
		Lang := utils.GetLangFromURL(h.Request.URL.String())

		// Maybe use Item for everything ?
		// Otherwise we are going to need 11 different functions
		// Maybe add omitempty field? `json:"droprates,omitempty"`
		GetTitle(h, Item, Lang)
		// fmt.Println("Got Title")
		GetTypeID(h, Item, Lang)
		// fmt.Println("Got TypeId")
		GetRarity(h, Item, Lang)
		// fmt.Println("Got Rarity")
		GetLevel(h, Item, Lang)
		// fmt.Println("Got Level")
		GetSubliStats(h, Item, Lang)
		// fmt.Println("Got Stats")
		GetDroprates(h, Item, Lang)
		// fmt.Println("Got Droprates")
		GetRecipes(h, Item, Lang)
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

func ScrapSingleResourceType(FileOptions utils.EditFileOptions) {
	urlSuffix := "&" + "type_1%5B%5D=" + strconv.Itoa(FileOptions.ID)

	var IDsList []int
	Items := make(map[int]structs.Item)

	c := GetNewCollector()

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

		frenchURL := FileOptions.SubCat.Item_url["fr"] + "/" + itemArgName
		englishURL := FileOptions.SubCat.Item_url["en"] + "/" + itemArgName

		var Item structs.Item
		// ParamsStatsProperties := structs.ParamsStatsProperties{AllPositivesStats: AllPositivesStats, AllNegativesStats: AllNegativesStats}

		Item.ID, err = utils.GetItemIDFromString(itemArgName)
		if err != nil {
			fmt.Println("error while getting the item ID", err)
			os.Exit(0)
		}

		// TODO : Refactor check for update
		// so we dont need to scrap what we already scraped

		IDsList = append(IDsList, Item.ID)
		// Scrap both FR/EN version of the item
		ScrapRessourceDetails(frenchURL, &Item)
		ScrapRessourceDetails(englishURL, &Item)

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
		fmt.Println("ScrapItemDetails visiting:\n", r.URL)
	})

	for i := 1; i < FileOptions.SubCat.MaxPage; i++ {
		if i != 1 && len(IDsList) > 0 {
			AppendIDsToFile(IDsList, FileOptions.SubCat.Title["fr"])
			AppendItemsToFile(Items, FileOptions.SubCat.Title["fr"])

		}
		IDsList = []int{}
		c.Visit(FileOptions.SubCat.Index_url["fr"] + strconv.Itoa(i) + urlSuffix)
		fmt.Printf("setting page to %d\n", i)
	}
}
