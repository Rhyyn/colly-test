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

func parseJSON(content []byte) map[string]structs.ItemInfo {
	var allItemsInfo map[string]structs.ItemInfo
	err := json.Unmarshal(content, &allItemsInfo)
	if err != nil {
		log.Printf("Error parsing JSON in parseJSON func: %v", err)
	}
	return allItemsInfo
}

// This Gets all Item subtypes from itemTypes.json
func GetItemTypesPropertiesJSON() []structs.TypesItem {
	openedFile := OpenFile("itemTypes.json")
	defer openedFile.Close()
	fileContent := ReadFile(openedFile)

	var TypesItems []structs.TypesItem
	err := json.Unmarshal(fileContent, &TypesItems)
	if err != nil {
		fmt.Println("Error decoding JSON:", err)
		return nil
	}
	return TypesItems
}

func GetJSON(file_name string) map[string]structs.ItemInfo {
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

func GetFileContent(file_name string) []structs.ItemInfo {
	opened_file := OpenFile(file_name)
	defer opened_file.Close()
	file_content := ReadFile(opened_file)

	allItemsInfo := parseJSON(file_content)

	sortedItems := make([]structs.ItemInfo, 0, len(allItemsInfo))
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
		for _, subCategory := range category.SubCategories {
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
		for _, subCategory := range category.SubCategories {
			if len(subCategory.ID) > 0 {
				if ID == subCategory.ID[0] {
					maxPage = subCategory.MaxPage
				}
			}
		}
	}
	return maxPage
}

// Updates MaxItems and MaxPage in file
func EditItemsCats(editFileOptions structs.EditFileOptions) {
	ID := editFileOptions.ID
	allCategoriesInfo := GetJSON("updated_item_categories.json")
	for _, category := range allCategoriesInfo {
		for i, subCategory := range category.SubCategories {
			if len(subCategory.ID) > 0 {
				if ID == subCategory.ID[0] {
					category.SubCategories[i].MaxItems = editFileOptions.SubCat.MaxItems
					category.SubCategories[i].MaxPage = editFileOptions.SubCat.MaxPage
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

// This is Used to convert our string values to int and assert what format of value we're storing
// Flat means 5 HP
// Percent means 5% Armor Given
// Negative means it starts with " - ", could be - 5 HP or - 5% Armor Given
func StatPrefixToStringAndSetFormat(str string) (string, int, bool) {
	var err error
	var value int
	isNegative := false
	var trimmedString string
	if strings.HasPrefix(str, "-") {
		// if starts with " - ", set isNegative, and trim
		trimmedString = strings.TrimPrefix(str, "-")
		isNegative = true
	}

	if strings.HasSuffix(str, "%") {
		format := "percent"
		// if percent trim % sign from string and convert to int type
		value, err = strconv.Atoi(strings.TrimSuffix(trimmedString, "%"))
		if err != nil {
			fmt.Println("Error converting string to int of trimmedString inside StatPrefixTostringAndSetFormat")
			fmt.Println(err)
		}

		return format, value, isNegative
	} else {
		format := "flat"
		value, err = strconv.Atoi(str)
		if err != nil {
			fmt.Println("Error converting string to int of supposed flat str inside StatPrefixTostringAndSetFormat")
			fmt.Println(err)
		}

		return format, value, isNegative
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

func ContainsElement(statString string) bool {
	if strings.Contains(statString, "Fire") ||
		strings.Contains(statString, "Water") ||
		strings.Contains(statString, "Earth") ||
		strings.Contains(statString, "Air") {
		return true
	}
	return false
}

func FormatSingleResStrng(statString string, format string) string {
	if strings.Contains(statString, "Resistance") &&
		ContainsElement(statString) {
		prefix, suffix := strings.Split(statString, " ")[0], strings.Split(statString, " ")[1]
		newString := suffix + " " + prefix
		return newString
	}
	return statString
}

// Used to inverse some word in english
// Resistance air -> Air Resistance
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
