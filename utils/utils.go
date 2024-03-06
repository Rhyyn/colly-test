package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"sort"
)

type SubCategory struct {
	Title     map[string]string `json:"title"`
	ID        []int             `json:"id"`
	Index_url map[string]string `json:"index_url"`
	Item_url  map[string]string `json:"item_url"`
	Max_page  int               `json:"max_page"`
	MaxItems  int               `json:"max_items"`
}

type ItemInfo struct {
	Title          map[string]string `json:"title"`
	ID             []int             `json:"id"`
	Index_url      map[string]string `json:"index_url"`
	Item_url       map[string]string `json:"item_url"`
	Max_page       int               `json:"max_page"`
	Sub_categories []SubCategory     `json:"sub_categories"`
}

type EditFileOptions struct {
	IsSubCat bool
	ID       int
	SubCat   SubCategory
}

type FileResult struct {
	File *os.File
}

func openFile(file_name string) *os.File {
	file, err := os.Open("./DATA/STATIC/" + file_name)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}
	return file
}

func readFile(opened_file *os.File) []byte {
	content, err := io.ReadAll(opened_file)
	if err != nil {
		log.Printf("Error reading file: %v", err)
	}
	return content
}

func parseJSON(content []byte) map[string]ItemInfo {
	var allItemsInfo map[string]ItemInfo
	err := json.Unmarshal(content, &allItemsInfo)
	if err != nil {
		log.Printf("Error parsing JSON: %v", err)
	}
	return allItemsInfo
}

func WriteJSON(allCategoriesInfo map[string]ItemInfo) {
	// fmt.Println(allCategoriesInfo)
	updatedJSON, err := json.Marshal(allCategoriesInfo)
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("item_categories.json", updatedJSON, 0o644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("JSON file updated successfully.")
}

func ReadAndPrintFile(file_name string) error {
	opened_file := openFile(file_name)
	defer opened_file.Close()
	file_content := readFile(opened_file)

	allItemsInfo := parseJSON(file_content)
	fmt.Printf("-----\n")
	for _, category := range allItemsInfo {
		fmt.Printf("Title (French): %s\nID: %v\nindex_url: %v\nitem_url: %v\n\n", category.Title["fr"], category.ID, category.Index_url["fr"], category.Item_url["fr"])
	}

	return nil
}

func GetFileContent(file_name string) []ItemInfo {
	opened_file := openFile(file_name)
	defer opened_file.Close()
	file_content := readFile(opened_file)
	allItemsInfo := parseJSON(file_content)

	sortedItems := make([]ItemInfo, 0, len(allItemsInfo))
	for _, item := range allItemsInfo {
		sortedItems = append(sortedItems, item)
	}
	sort.Slice(sortedItems, func(i, j int) bool {
		return sortedItems[i].ID[0] < sortedItems[j].ID[0]
	})

	return sortedItems
}

func EditItemsCats(editFileOptions EditFileOptions) {
	fmt.Println(editFileOptions.SubCat.Max_page)
	ID := editFileOptions.ID
	isSubCat := editFileOptions.IsSubCat
	fmt.Printf("ID : %d , isSubCat %t\n", ID, isSubCat)
	opened_file := openFile("item_categories.json")
	defer opened_file.Close()
	file_content := readFile(opened_file)
	allCategoriesInfo := parseJSON(file_content)
	for _, category := range allCategoriesInfo {
		for i, subCategory := range category.Sub_categories {
			if len(subCategory.ID) > 0 {
				fmt.Println(subCategory.ID[0])
				if ID == subCategory.ID[0] {
					category.Sub_categories[i].MaxItems = editFileOptions.SubCat.MaxItems
					category.Sub_categories[i].Max_page = editFileOptions.SubCat.Max_page
					break
				}
			}
		}
	}
	WriteJSON(allCategoriesInfo)
}
