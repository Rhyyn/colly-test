package cli

import (
	"fmt"

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
		var category_choice int
		fmt.Println("You selected: Scrap item type")
		allCategoriesInfo := utils.GetFileContent("item_categories.json")

		for index, item := range allCategoriesInfo {
			fmt.Printf("Index: %v, item %#v\n", index, item.Title["fr"])
		}

		fmt.Print("Chose a category: (use numbers..)\n")
		fmt.Scanln(&category_choice)

		selected_category := allCategoriesInfo[category_choice]
		fmt.Printf("You selected %s\n", selected_category.Title["fr"])
		for index, subCategory := range selected_category.Sub_categories {
			fmt.Printf("%d. Type : %#v\n", index, subCategory.Title["fr"])
		}
		fmt.Print("Chose a type: (use numbers..)\n")

		var type_choice int
		fmt.Scanln(&type_choice)
		selected_type := selected_category.Sub_categories[type_choice]
		fmt.Printf("You selected %s\n", selected_type.Title["fr"])

		fmt.Printf("Begin crawling %s...\n", selected_type.Title["fr"])
		defer scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
			Title:     selected_type.Title["fr"],
			Index_url: selected_category.Index_url["fr"],
			ID:        selected_type.ID,
		})
		// defer scrapers.CrawlIndexURL(selected_category.Index_url["fr"])
		// scrapers.CrawlSingleAccesoryType(selectedType.ID[0])

	case 3:
		fmt.Println("You selected: Scrap category")
	default:
		fmt.Println("Invalid choice")
	}
}
