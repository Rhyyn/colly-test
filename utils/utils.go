package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/Rhyyn/wakfukiscraper/structs"
)

type ScrapedItem struct {
	ID int `json:"_id"`
}

type SubCategory struct {
	Title     map[string]string `json:"title"`
	ID        []int             `json:"id"`
	Index_url map[string]string `json:"index_url"`
	Item_url  map[string]string `json:"item_url"`
	MaxPage   int               `json:"max_page"`
	MaxItems  int               `json:"max_items"`
}

type ItemInfo struct {
	Title          map[string]string `json:"title"`
	ID             []int             `json:"id"`
	Index_url      map[string]string `json:"index_url"`
	Item_url       map[string]string `json:"item_url"`
	MaxPage        int               `json:"max_page"`
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

func OpenFile(file_name string) *os.File {
	file, err := os.Open("./DATA/STATIC/" + file_name)
	if err != nil {
		log.Printf("Error opening file: %v", err)
	}
	return file
}

func ReadFile(opened_file *os.File) []byte {
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

func GetJSON(file_name string) map[string]ItemInfo {
	opened_file := OpenFile(file_name)
	defer opened_file.Close()
	file_content := ReadFile(opened_file)
	allItemsInfo := parseJSON(file_content)
	return allItemsInfo
}

// allCategoriesInfo map[string]ItemInfo
func WriteJSON(data interface{}) {
	fileContent, err := json.MarshalIndent(data, "", "    ")
	if err != nil {
		fmt.Println("Error marshaling JSON:", err)
		return
	}

	err = os.WriteFile("./DATA/STATIC/updated_item_categories.json", fileContent, 0o644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Println("JSON file updated successfully.")
}

func ReadAndPrintFile(file_name string) error {
	opened_file := OpenFile(file_name)
	defer opened_file.Close()
	file_content := ReadFile(opened_file)

	allItemsInfo := parseJSON(file_content)
	fmt.Printf("-----\n")
	for _, category := range allItemsInfo {
		fmt.Printf("Title (French): %s\nID: %v\nindex_url: %v\nitem_url: %v\n\n", category.Title["fr"], category.ID, category.Index_url["fr"], category.Item_url["fr"])
	}

	return nil
}

func GetFileContent(file_name string) []ItemInfo {
	opened_file := OpenFile(file_name)
	defer opened_file.Close()
	file_content := ReadFile(opened_file)

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

func GetMaxItems(ID int) int {
	allCategoriesInfo := GetJSON("updated_item_categories.json")
	var maxItems int
	for _, category := range allCategoriesInfo {
		for _, subCategory := range category.Sub_categories {
			if len(subCategory.ID) > 0 {
				if ID == subCategory.ID[0] {
					maxItems = subCategory.MaxItems
				}
			}
		}
	}
	return maxItems
}

func GetMaxPage(ID int) int {
	allCategoriesInfo := GetJSON("updated_item_categories.json")
	var maxPage int
	for _, category := range allCategoriesInfo {
		for _, subCategory := range category.Sub_categories {
			if len(subCategory.ID) > 0 {
				if ID == subCategory.ID[0] {
					maxPage = subCategory.MaxPage
				}
			}
		}
	}
	return maxPage
}

// Updates MaxItems and MaxPage
func EditItemsCats(editFileOptions EditFileOptions) {
	ID := editFileOptions.ID
	allCategoriesInfo := GetJSON("updated_item_categories.json")
	for _, category := range allCategoriesInfo {
		for i, subCategory := range category.Sub_categories {
			if len(subCategory.ID) > 0 {
				if ID == subCategory.ID[0] {
					category.Sub_categories[i].MaxItems = editFileOptions.SubCat.MaxItems
					category.Sub_categories[i].MaxPage = editFileOptions.SubCat.MaxPage
					break
				}
			}
		}
	}
	WriteJSON(allCategoriesInfo)
}

func ReplaceAnyNumberInString(str string, replacementStr string) string {
	re := regexp.MustCompile("[0-9]+")
	return re.ReplaceAllString(str, replacementStr)
}

func HasNumberInString(str string) (bool, int) {
	re := regexp.MustCompile("[0-9]+")
	// found := re.MatchString(str)
	matches := re.FindStringSubmatch(str)
	if len(matches) > 0 {
		num, _ := strconv.Atoi(matches[0])
		return true, num
	}
	return false, 0
	// return found
}

func StatPrefixToStringAndSetFormat(str string) (string, int) {
	if strings.HasPrefix(str, "-") {
		format := "negative"
		value, _ := strconv.Atoi(strings.TrimPrefix(str, "-"))
		return format, value
	} else if strings.HasSuffix(str, "%") {
		format := "percent"
		value, _ := strconv.Atoi(strings.TrimSuffix(str, "%"))
		return format, value
	} else {
		format := "flat"
		value, _ := strconv.Atoi(str)
		return format, value
	}
}

func FormatElementsString(str string, format string) string {
	if format == "negative" {
		switch str {
		case "- Maîtrise sur X éléments aléatoires":
			return "- Maîtrise dans X éléments"
		case "- Mastery of X random elements":
			return "- Mastery in X elements"
		case "- Résistance sur X éléments aléatoires":
			return "- Résistance dans X éléments"
		case "- Resistance to X random elements":
			return "- Resistance in X elements"
		}
	} else {
		switch str {
		case "Maîtrise sur X éléments aléatoires":
			return "Maîtrise dans X éléments"
		case "Mastery of X random elements":
			return "Mastery in X elements"
		case "Résistance sur X éléments aléatoires":
			return "Résistance dans X éléments"
		case "Resistance to X random elements":
			return "Resistance in X elements"
		}
		return str
	}
	return str
}

func FormatCritString(statString string, format string) string {
	if format == "negative" {
		switch statString {
		case "- Critical Hit":
			return "- Critical Chance (%)"
		case "- Coup critique":
			return "- Coup Critique (%)"
		}
	} else {
		switch statString {
		case "Critical Hit":
			return "Critical Chance (%)"
		case "Coup critique":
			return "Coup Critique (%)"
		}
		return statString
	}
	return statString
}

func FormatSingleResStrng(statString string, format string) string {
	if strings.Contains(statString, "Resistance") {
		prefix, suffix := strings.Split(statString, " ")[0], strings.Split(statString, " ")[1]
		newString := suffix + " " + prefix
		return newString
	}
	return statString
}

func FormatStatString(statString string, format string) string {
	newStatString := statString
	newStatString = FormatSingleResStrng(statString, format)
	if format == "negative" {
		newStatString = "- " + newStatString
	}
	newStatString = FormatElementsString(newStatString, format)
	newStatString = FormatCritString(newStatString, format)
	return newStatString
}

// This create the map of properties for stats
func HandleStatsProperties(data []byte) map[string]structs.StatProperties {
	var translations map[string]structs.StatProperties

	err := json.Unmarshal([]byte(data), &translations)
	if err != nil {
		panic(err)
	}

	return translations
}

func GetLangFromURL(str string) string {
	if strings.Contains(str, "/fr") {
		return "Fr"
	} else {
		return "En"
	}
}

func GetItemIDFromString(str string) (int, error) {
	ID, err := strconv.Atoi(strings.Split(str, "-")[0])
	if err != nil {
		fmt.Println(err)
		return 0, err
	}
	return ID, nil
}
