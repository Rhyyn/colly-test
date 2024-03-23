package cli

import (
	"fmt"
	"os"

	"github.com/Rhyyn/wakfukiscraper/scrapers"
	"github.com/Rhyyn/wakfukiscraper/structs"
	"github.com/Rhyyn/wakfukiscraper/utils"
)

func Execute() {
	fmt.Println("Select a command: (use numbers : 1,2..)")
	fmt.Println("1. Print item categories")
	fmt.Println("2. Scrap TYPE of item")
	fmt.Println("3. Scrap CATEGORY of item")
	fmt.Println("4. Scrap url of item")

	var choice int
	fmt.Print("Enter your choice: ")
	fmt.Scanln(&choice)

	switch choice {
	case 1:
		fmt.Println("You selected: Print item categories")
		utils.ReadAndPrintFile("updated_item_categories.json")
	case 2:

		fmt.Println("You selected: Scrap item type")
		allCategoriesInfo := utils.GetFileContent("updated_item_categories.json")

		// Ask for Category choice
		categoryChoice := selectCategory(allCategoriesInfo)
		selectedCategory := allCategoriesInfo[categoryChoice]
		fmt.Printf("You selected %s\n", selectedCategory.Title["fr"])

		// Ask for SubCategory choice
		typeChoice := selectSubCategory(selectedCategory)
		selectedType := selectedCategory.SubCategories[typeChoice]
		fmt.Printf("You selected %s\n", selectedType.Title["fr"])

		// This if Else is temporary, this needs better logic
		// If SubCategory == ressources , ask for sub of sub category
		if selectedType.Title["fr"] == "Améliorations" {
			selectedSubType := selectTypeInsideSubCategory(selectedType)
			// ID / ParentID / TITLES

			EditFileOptions := scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
				Title:     selectedSubType.Title.Fr,
				Index_url: selectedCategory.Index_url["fr"],
				ID:        append([]int{}, selectedSubType.Definition.ID),
			})

			ScrapingParameters := structs.ScrapingParameters{
				IndexUrl:       selectedCategory.Index_url,
				ItemUrl:        selectedCategory.Item_url,
				MaxPage:        EditFileOptions.SubCat.MaxPage,
				SelectedId:     selectedSubType.Definition.ID,
				SingleItemMode: false,
				SelectedType:   selectedType.Title["fr"],
			}
			// call scraper here with params from EditFileOptions
			scrapers.ScrapSingleResourceType(ScrapingParameters)
			fmt.Println(EditFileOptions.SubCat)
		} else {

			// if no maxItems stored proceed without prompts else
			// Ask if we want to check for new items (max_items, max_page)
			selectedId := selectedCategory.SubCategories[typeChoice].ID[0]
			maxItems := utils.GetMaxItems(selectedId)
			// fmt.Printf("maxItems : %d\n", maxItems)
			if maxItems != 0 {
				checkForNewItems(selectedType, selectedCategory, selectedId)
			} else {
				// TODO : DRY, need to extract
				// Call a new separate Collector to get info on the max number of items in that category
				EditFileOptions := scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
					Title:     selectedType.Title["fr"],
					Index_url: selectedCategory.Index_url["fr"],
					ID:        selectedType.ID,
				})
				// Update the json file with it
				defer utils.EditItemsCats(EditFileOptions)
				// scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
				// 	Title:     selectedType.Title["fr"],
				// 	Index_url: selectedCategory.Index_url["fr"],
				// 	ID:        selectedType.ID,
				// })
				maxPage := utils.GetMaxPage(selectedId)
				ScrapingParameters := structs.ScrapingParameters{
					IndexUrl:       selectedCategory.Index_url,
					ItemUrl:        selectedCategory.Item_url,
					MaxPage:        maxPage,
					SelectedId:     selectedId,
					SingleItemMode: false,
					SelectedType:   selectedType.Title["fr"],
				}
				scrapers.ScrapItems(ScrapingParameters)
				// scrapers.ScrapRedirect(selectedCategory.Index_url["fr"], maxPage, selectedId)
			}
			// defer scrapers.CrawlIndexURL(selected_category.Index_url["fr"])
			// scrapers.CrawlSingleAccesoryType(selectedType.ID[0])

		}

	case 3:
		fmt.Println("You selected: Scrap category")
	case 4:
		fmt.Println("You selected: Scrap url of item")
		fmt.Println("Please Input the FRENCH url of the item: ")
		var frURL, enURL string
		urls := make(map[string]string)

		fmt.Scanln(&frURL)
		fmt.Println("Please Input the ENGLISH url of the item: ")
		fmt.Scanln(&enURL)

		urls["Fr"] = frURL
		urls["En"] = enURL

		ScrapingParameters := structs.ScrapingParameters{
			SingleItemURL:  urls,
			SingleItemMode: true,
		}
		scrapers.ScrapItems(ScrapingParameters)
	default:
		fmt.Println("Invalid choice")
	}
}

// ask and select a sub category
func selectSubCategory(selectedCategory structs.ItemInfo) int {
	for index, subCategory := range selectedCategory.SubCategories {
		fmt.Printf("%d. Type : %#v\n", index, subCategory.Title["fr"])
	}

	fmt.Print("Chose a type: (use numbers..)\n")
	var choice int
	fmt.Scanln(&choice)

	if choice < 0 || choice >= len(selectedCategory.SubCategories) {
		fmt.Println("Invalid choice.")
		os.Exit(0)
	}

	return choice
}

// ask and select a type of object inside a sub category
func selectTypeInsideSubCategory(selectedType structs.SubCategory) structs.TypesItem {
	ItemTypesProperties := utils.GetItemTypesPropertiesJSON()

	var selectedObjects []structs.TypesItem
	for _, t := range ItemTypesProperties {
		if t.Definition.ParentID == selectedType.ID[0] {
			selectedObjects = append(selectedObjects, t)
		}
	}

	for index, st := range selectedObjects {
		fmt.Printf("%d. Type : %#v\n", index, st.Title.Fr)
	}

	fmt.Print("Chose a type: (use numbers..)\n")
	var choice int
	fmt.Scanln(&choice)

	if choice < 0 || choice >= len(selectedObjects) {
		fmt.Println("Invalid choice.")
		os.Exit(0)
	}

	return selectedObjects[choice]
}

func selectCategory(allCategoriesInfo []structs.ItemInfo) int {
	for index, item := range allCategoriesInfo {
		fmt.Printf("Index: %v, item %#v\n", index, item.Title["fr"])
	}
	var categoryChoice int
	fmt.Print("Chose a category: (use numbers..)\n")
	fmt.Scanln(&categoryChoice)

	if categoryChoice < 0 || categoryChoice >= len(allCategoriesInfo) {
		fmt.Println("Invalid choice.")
		os.Exit(0)
	}
	return categoryChoice
}

func checkForNewItems(selectedType structs.SubCategory, selectedCategory structs.ItemInfo, selectedId int) {
	var userUpdate string
	fmt.Printf("Do you want to check for new items ? (y/n)")
	fmt.Scanln(&userUpdate)
	if userUpdate != "y" && userUpdate != "n" && userUpdate != "" {
		fmt.Printf("wrong input please use (y/n), exiting..\n")
		os.Exit(0)
	} else if userUpdate == "y" {
		// Call a new separate Collector to get info on the max number of items in that category
		EditFileOptions := scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
			Title:     selectedType.Title["fr"],
			Index_url: selectedCategory.Index_url["fr"],
			ID:        selectedType.ID,
		})
		// Update the json file with it
		defer utils.EditItemsCats(EditFileOptions)
		// defer scrapers.UpdateMaxItemsAndPages(scrapers.IndexOptions{
		// 	Title:     selectedType.Title["fr"],
		// 	Index_url: selectedCategory.Index_url["fr"],
		// 	ID:        selectedType.ID,
		// })
		maxPage := utils.GetMaxPage(selectedId)
		ScrapingParameters := structs.ScrapingParameters{
			IndexUrl:       selectedCategory.Index_url,
			ItemUrl:        selectedCategory.Item_url,
			MaxPage:        maxPage,
			SelectedId:     selectedId,
			SingleItemMode: false,
			SelectedType:   selectedType.Title["fr"],
		}
		scrapers.ScrapItems(ScrapingParameters)
		// scrapers.ScrapRedirect(selectedCategory.Index_url["fr"], maxPage, selectedId)
	} else {
		maxPage := utils.GetMaxPage(selectedId)
		ScrapingParameters := structs.ScrapingParameters{
			IndexUrl:       selectedCategory.Index_url,
			ItemUrl:        selectedCategory.Item_url,
			MaxPage:        maxPage,
			SelectedId:     selectedId,
			SingleItemMode: false,
			SelectedType:   selectedType.Title["fr"],
		}
		scrapers.ScrapItems(ScrapingParameters)
		// scrapers.ScrapRedirect(selectedCategory.Index_url["fr"], maxPage, selectedId)
	}
}
