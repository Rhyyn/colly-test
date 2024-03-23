package scrapers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"

	"github.com/Rhyyn/wakfukiscraper/structs"
	"github.com/Rhyyn/wakfukiscraper/utils"
	"github.com/gocolly/colly"
)

func ScrapRessourceDetails(url string, Item *structs.Item, Recipes map[int]structs.Recipe) {
	c := GetNewCollector()
	c.OnHTML(".ak-container.ak-panel-stack.ak-glue", func(h *colly.HTMLElement) {
		Lang := utils.GetLangFromURL(h.Request.URL.String())
		// Recipe := structs.Recipe{}
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

func ScrapSingleResourceType(ScrapingParameters structs.ScrapingParameters) {
	urlSuffix := "&" + "type_1%5B%5D=" + strconv.Itoa(ScrapingParameters.SelectedId)

	fmt.Println(urlSuffix)

	var IDsList []int
	Items := make(map[int]structs.Item)
	Recipes := make(map[int]structs.Recipe)

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

		frenchURL := ScrapingParameters.ItemUrl["fr"] + itemArgName
		englishURL := ScrapingParameters.ItemUrl["en"] + "/" + itemArgName

		var Item structs.Item
		// var Recipe structs.Recipe
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
		ScrapRessourceDetails(frenchURL, &Item, Recipes)
		ScrapRessourceDetails(englishURL, &Item, Recipes)

		Items[Item.ID] = Item
		// Recipes[Recipe.RecipeId] = Recipe
		// Useless pretty print for debug
		PrettyItem, err := json.MarshalIndent(Item, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling item:", err)
			return
		}
		fmt.Println("Item after scraping:\n", string(PrettyItem))
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("ScrapItemDetails visiting:\n", r.URL)
	})

	for i := 1; i < ScrapingParameters.MaxPage; i++ {
		if i != 1 && len(IDsList) > 0 {
			AppendIDsToFile(IDsList, ScrapingParameters.SelectedType)
			AppendItemsToFile(Items, ScrapingParameters.SelectedType)
			AppendRecipesToFile(Recipes)
		}
		Items = map[int]structs.Item{}
		Recipes = map[int]structs.Recipe{}
		IDsList = []int{}
		fmt.Println(ScrapingParameters.IndexUrl)
		fmt.Println(ScrapingParameters.IndexUrl["fr"] + strconv.Itoa(i) + urlSuffix)
		c.Visit(ScrapingParameters.IndexUrl["fr"] + strconv.Itoa(i) + urlSuffix)
		fmt.Printf("setting page to %d\n", i)
	}
}
