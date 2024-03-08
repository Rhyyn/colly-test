package cli

import (
	"fmt"
	"os"

	"github.com/Rhyyn/wakfukiscraper/scrapers"
	"github.com/Rhyyn/wakfukiscraper/utils"
)

func Execute() {
	fmt.Println("Select a command: (use numbers : 1,2..)")
	fmt.Println("1. Print item categories")
	fmt.Println("2. Scrap TYPE of item")
	fmt.Println("3. Scrap CATEGORY of item")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Println("You selected: Print item categories")
		utils.ReadAndPrintFile("item_categories.json")
	case 2:

		fmt.Println("You selected: Scrap item type")
		allCategoriesInfo := utils.GetFileContent("item_categories.json")

		// Ask for Category choice
		categoryChoice := selectCategory(allCategoriesInfo)
		selectedCategory := allCategoriesInfo[categoryChoice]
		fmt.Printf("You selected %s\n", selectedCategory.Title["fr"])


		// Ask for SubCategory choice
		typeChoice := selectSubCategory(selectedCategory)
		selectedType := selectedCategory.Sub_categories[typeChoice]
		fmt.Printf("You selected %s\n", selectedType.Title["fr"])

		// if no maxItems stored proceed without prompts else
		// Ask if we want to check for new items (max_items, max_page)
		selectedId := selectedCategory.Sub_categories[typeChoice].ID[0]
		maxItems := utils.GetMaxItems(selectedId)
		// fmt.Printf("maxItems : %d\n", maxItems)
		if maxItems != 0 {
			checkForNewItems(selectedType, selectedCategory, selectedId)
		} else {
			scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
				Title:     selectedType.Title["fr"],
				Index_url: selectedCategory.Index_url["fr"],
				ID:        selectedType.ID,
			})
			maxPage := utils.GetMaxPage(selectedId)
			scrapers.ScrapItems(selectedCategory.Index_url["fr"], maxPage, selectedId)
		}
		// defer scrapers.CrawlIndexURL(selected_category.Index_url["fr"])
		// scrapers.CrawlSingleAccesoryType(selectedType.ID[0])

	case 3:
		fmt.Println("You selected: Scrap category")
	default:
		fmt.Println("Invalid choice")
	}
}

func selectSubCategory(selectedCategory utils.ItemInfo) int {
	for index, subCategory := range selectedCategory.Sub_categories {
		fmt.Printf("%d. Type : %#v\n", index, subCategory.Title["fr"])
	}
	fmt.Print("Chose a type: (use numbers..)\n")
	var choice int
	fmt.Scanln(&choice)
	return choice
}

func selectCategory(allCategoriesInfo []utils.ItemInfo) int {
	for index, item := range allCategoriesInfo {
		fmt.Printf("Index: %v, item %#v\n", index, item.Title["fr"])
	}
	var categoryChoice int
	fmt.Print("Chose a category: (use numbers..)\n")
	fmt.Scanln(&categoryChoice)
	return categoryChoice
}

func checkForNewItems(selectedType utils.SubCategory, selectedCategory utils.ItemInfo, selectedId int) {
	var userUpdate string
	fmt.Printf("Do you want to check for new items ? (y/n)")
	fmt.Scanln(&userUpdate)
	if userUpdate != "y" && userUpdate != "n" && userUpdate != "" {
		fmt.Printf("wrong input please use (y/n), exiting..\n")
		os.Exit(0)
	} else if userUpdate == "y" {
		defer scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
			Title:     selectedType.Title["fr"],
			Index_url: selectedCategory.Index_url["fr"],
			ID:        selectedType.ID,
		})
		maxPage := utils.GetMaxPage(selectedId)
		scrapers.ScrapItems(selectedCategory.Index_url["fr"], maxPage, selectedId)
	} else {
		maxPage := utils.GetMaxPage(selectedId)
		scrapers.ScrapItems(selectedCategory.Index_url["fr"], maxPage, selectedId)
	}
}
