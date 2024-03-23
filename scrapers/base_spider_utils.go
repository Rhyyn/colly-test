package scrapers

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/Rhyyn/wakfukiscraper/structs"
	"github.com/Rhyyn/wakfukiscraper/utils"
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

func AppendRecipesToFile(newRecipes map[int]structs.Recipe) {
	filePath := "./DATA/STATIC/ScrapedData/Recipes/recipes.json"

	var Recipes map[int]structs.Recipe

	if _, err := os.Stat(filePath); err == nil {
		file, err := os.ReadFile(filePath)
		if err != nil {
			fmt.Println("Error reading file:", err)
			return
		}

		if err := json.Unmarshal(file, &Recipes); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}
	} else if os.IsNotExist(err) {
		Recipes = map[int]structs.Recipe{}
	} else {
		fmt.Println("Error checking file status:", err)
		return
	}

	for recipeId, newRecipe := range newRecipes {
		Recipes[recipeId] = newRecipe
	}

	jsonData, err := json.MarshalIndent(Recipes, "", "  ")
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	if err := os.WriteFile(filePath, jsonData, 0o644); err != nil {
		fmt.Println("Error writing to file:", err)
		return
	}

	fmt.Println("Recipes appended successfully.")
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

func StandardizeSpaces(s string) string {
	return strings.Join(strings.Fields(s), " ")
}

func GetEnglishURL(frenchUrl string) string {
	index := strings.LastIndex(frenchUrl, "/")
	if index == -1 {
		fmt.Println("Invalid URL")
	}

	baseURL := frenchUrl[:index+1]
	resource := frenchUrl[index+1:]

	englishURL := strings.ReplaceAll(baseURL, "/fr/", "/en/") + resource
	return englishURL
}

func GetIdFromKey(key string) (int, error) {
	id, err := strconv.Atoi(key)
	if err != nil {
		return 0, err
	}
	return id, nil
}

func GetTitle(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	titleElement := htmlElement.DOM.Find("h1").First()
	title := StandardizeSpaces(titleElement.Text())
	if lang == "Fr" {
		Item.Title.Fr = title
	} else {
		Item.Title.En = title
	}
}

func GetTypeID(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	// Extracting typId from img src of type icon
	typeElement := htmlElement.DOM.Find(".ak-encyclo-detail-type.col-xs-6 span img").First()
	src, exist := typeElement.Attr("src")
	if !exist {
		fmt.Println("typeElement does not exist")
	}
	typeParts := strings.Split(src, "/")
	category := typeParts[len(typeParts)-1]
	typeId, err := strconv.Atoi(strings.Split(category, ".")[0])
	if err != nil {
		fmt.Printf("Error converting typeId")
		os.Exit(0)
	}
	Item.Params.TypeId = typeId
}

// Gets the ID from separate file
func GetStatId(statString string, ParamsStatsProperties structs.ParamsStatsProperties, isNegative bool, lang string) (int, error) {
	var id int
	var err error
	// TODO :
	// Need to handle subilmations somehow
	// ScrapItemDetails visiting:
	// https://www.wakfu.com/fr/mmorpg/encyclopedie/armures/27302-bottes-assechees
	// id of stat X% (+X Niv.) is 0
	// https://github.com/noredlace/wakfu-sublimations/blob/main/wakfu-sublimations-site/src/app/data/sublimations.json
	if isNegative {
		for key, property := range ParamsStatsProperties.AllNegativesStats {
			if lang == "Fr" {
				if strings.Contains(statString, property.Fr) {
					id, err = GetIdFromKey(key)
					return id, err
				}
			} else {
				if strings.Contains(statString, property.En) {
					id, err = GetIdFromKey(key)
					return id, err
				}
			}
		}
	} else {
		for key, property := range ParamsStatsProperties.AllPositivesStats {
			if lang == "Fr" {
				if strings.Contains(statString, property.Fr) {
					id, err = GetIdFromKey(key)
					return id, err
				}
			} else {
				if strings.Contains(statString, property.En) {
					id, err = GetIdFromKey(key)
					return id, err
				}
			}
		}
	}
	return id, nil
}

func GetStats(htmlElement *colly.HTMLElement, Item *structs.Item, lang string, ParamsStatsProperties structs.ParamsStatsProperties) {
	statElement := htmlElement.DOM.Find(".ak-container.ak-content-list.ak-displaymode-col")
	Stats := make(map[int]structs.Stat)
	statElement.Each(func(i int, s *goquery.Selection) {
		statsDiv := s.Find("div.ak-title")

		statsDiv.Each(func(i int, stat *goquery.Selection) {
			var id int
			var idError error
			var displayString string
			var format string
			var value int
			var numElements int

			entireStatString := strings.SplitN(StandardizeSpaces(stat.Text()), " ", 2)
			// "160"
			prefixString := StandardizeSpaces(entireStatString[0])
			// "Esquive"
			suffixString := StandardizeSpaces(entireStatString[1])

			// Used to check if stat is either flat/percent/negative
			// TODO : need do make it do it can become "negative,flat" or "negative,percent"
			// or return a isNegative flag ?
			// fmt.Printf("prefixString : %s\n", prefixString)
			var isNegative bool
			format, value, isNegative = utils.StatPrefixToStringAndSetFormat(prefixString)

			// Check if has 2 values (Mastery of 3 random elements)
			// If it has set numElements and format it to (Mastery in X elements)
			if hasNumber, number := utils.HasNumberInString(suffixString); hasNumber {
				numElements = number
				// fmt.Println("numElements", numElements)
				suffixString = utils.ReplaceAnyNumberInString(suffixString, "X")
				// fmt.Println("suffixString", suffixString)
			}
			// Format the strings because Ankama's english is dogshit
			suffixString = utils.FormatStatString(suffixString, format)
			// fmt.Println("formatted suffixString", suffixString)

			// Gets the ID
			id, idError = GetStatId(suffixString, ParamsStatsProperties, isNegative, lang)
			if idError != nil {
				fmt.Printf("Error GetStatId %v", idError)
			}

			if id == 0 {
				fmt.Printf("id of stat %s is 0", suffixString)
				os.Exit(0)
			}

			displayString = suffixString

			// If scraping english vesion modify the stats to add english string
			if lang == "En" && len(Item.Stats) > 0 {
				for key, stat := range Item.Stats {
					if stat.ID == id {
						stat.Display.En = displayString
						Item.Stats[key] = stat
					}
				}
			}

			if lang == "Fr" {
				newStat := structs.Stat{
					Display: structs.Display{
						Fr: displayString,
						En: displayString,
					},
					ID:          id,
					Format:      format,
					IsNegative:  isNegative,
					Value:       value,
					NumElements: numElements,
				}

				Stats[newStat.ID] = newStat
			}
		})
	})
	if lang == "Fr" {
		Item.Stats = Stats
	}
}

func GetSubliStats(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	newSubliDetails := structs.SublimationDetails{}
	statElement := htmlElement.DOM.Find(".ak-container.ak-content-list.ak-displaymode-col")
	statElement.Each(func(i int, s *goquery.Selection) {
		statsDiv := s.Find("div.ak-title").First()

		statsDiv.Each(func(i int, stat *goquery.Selection) {
			// Extract name of the Subli
			displayString := StandardizeSpaces(statsDiv.Find("a").Text())
			if lang == "En" && len(Item.SublimationDetails.Title.Fr) > 0 {
				displayString = strings.Split(displayString, ">")[1]
				Item.SublimationDetails.Title.En = displayString
			}

			if lang == "Fr" {
				// Find and extract level increase from the string (1.Niv)
				// If contains + isEpic true
				var LevelIncValue int
				isEpic := false
				LevelIncString := strings.Split(StandardizeSpaces(statsDiv.Text()), "(")[1]
				// TODO: this does not work
				if !strings.Contains(LevelIncString, "+") {
					isEpic = true
				}
				if hasNumber, number := utils.HasNumberInString(LevelIncString); hasNumber {
					LevelIncValue = number
				}

				// Find and extract the supposed ID of the sublimation throught a leftover paramater called data-state

				cleanedDataState := strings.ReplaceAll(statsDiv.Find("a").AttrOr("data-state", ""), "\u00a0", " ")
				cleanedDataState = StandardizeSpaces(cleanedDataState)
				cleanedDataState = strings.ReplaceAll(cleanedDataState, " ", "")
				fmt.Println("Cleaned data state:", cleanedDataState)
				dataStateId, err := strconv.Atoi(cleanedDataState)
				if err != nil {
					fmt.Println("Error getting the sublimation ID through data-state")
					fmt.Println(err)
					os.Exit(0)
				}
				// TODO : fix Title lose its number Ravage II > Ravage
				newSubliDetails = structs.SublimationDetails{
					LevelInc: LevelIncValue,
					Title: structs.Title{
						Fr: displayString,
						En: displayString,
					},
					ID:     dataStateId,
					IsEpic: isEpic,
				}
			}
		})
	})
	if lang == "Fr" {
		Item.SublimationDetails = newSubliDetails
	}
}

func GetRarity(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	rarityElement := htmlElement.DOM.Find(".ak-object-rarity span span")
	rarityNumber, err := strconv.Atoi(strings.Split(rarityElement.AttrOr("class", ""), "ak-rarity-")[1])
	if err != nil {
		fmt.Println(err)
	}
	Item.Params.Rarity = rarityNumber
}

func GetLevel(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	levelElement := htmlElement.DOM.Find(".ak-encyclo-detail-level.col-xs-6.text-right")
	levelElementFr := StandardizeSpaces(levelElement.Text())
	level, err := strconv.Atoi(strings.TrimSpace(strings.Split(levelElementFr, ": ")[1]))
	if err != nil {
		fmt.Printf("Error converting level")
		os.Exit(0)
	}
	Item.Params.Level = level
}

func GetDroprates(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	droprateContainer := htmlElement.DOM.Find(".ak-container.ak-panel:has(.ak-panel-title:contains('Peut être obtenu sur')), .ak-container.ak-panel:has(.ak-panel-title:contains('Dropped By'))")

	// TODO : what happens if return empty?
	if droprateContainer.Length() == 0 {
		fmt.Println("No elements found")
		return
	}
	Droprates := make(map[int]structs.Droprate)

	// Gets the droprates table container and forEach rows construct a droprate
	dropratesRows := htmlElement.DOM.Find(".ak-column.ak-container.col-xs-12.col-md-6")
	dropratesRows.Each(func(_ int, s *goquery.Selection) {
		var monsterArgName string
		var monsterID int
		var err error
		ak_image := s.Find(".ak-image a[href]")
		monsterHref, exists := ak_image.Attr("href")
		if exists {
			monsterArgName, err = GetItemURLArg(monsterHref)
			if err != nil {
				fmt.Printf("error getting monsterArgName from %s\n", monsterHref)
			}
			monsterID, err = utils.GetItemIDFromString(monsterArgName)
			if err != nil {
				fmt.Printf("error getting id from %s\n", monsterHref)
			}
		}

		// Gets the monsterName span
		monsterName := s.Find(".ak-content .ak-title span ").Text()

		// Gets the drop as string "0.25%", removes % and parse to float64
		dropChanceString := s.Find(".ak-aside").Text()
		dropChanceValue := strings.Split(dropChanceString, "%")[0]
		dropChance, err := strconv.ParseFloat(dropChanceValue, 64)
		if err != nil {
			fmt.Println("Error converting dropChance string")
			fmt.Println(err)
		}

		// Update the EN name only
		if lang == "En" && len(Item.Droprates) > 0 {
			for key, droprate := range Item.Droprates {
				if droprate.MonsterID == monsterID {
					droprate.MonsterName.En = monsterName
					Item.Droprates[key] = droprate
				}
			}
		}

		if lang == "Fr" {
			Droprate := structs.Droprate{
				MonsterID: monsterID,
				MonsterName: structs.Display{
					Fr: monsterName,
					En: monsterName,
				},
				DropChance: dropChance,
			}

			Droprates[Droprate.MonsterID] = Droprate
		}
	})

	// Alwasy create the Droprates inside FR scraping
	if lang == "Fr" {
		Item.Droprates = Droprates
	}
}

func GetRecipes(htmlElement *colly.HTMLElement, Item *structs.Item, Recipes map[int]structs.Recipe, lang string) {
	jobs := map[string]map[string]string{
		"75": {"fr": "Pêcheur", "en": "Fisherman"},
		"71": {"fr": "Forestier", "en": "Lumberjack"},
		"72": {"fr": "Herboriste", "en": "Herbalist"},
		"64": {"fr": "Paysan", "en": "Farmer"},
		"73": {"fr": "Mineur", "en": "Miner"},
		"74": {"fr": "Trappeur", "en": "Trapper"},
		"77": {"fr": "Armurier", "en": "Armorer"},
		"78": {"fr": "Bijoutier", "en": "Jeweler"},
		"40": {"fr": "Boulanger", "en": "Baker"},
		"76": {"fr": "Cuisinier", "en": "Chef"},
		"81": {"fr": "Ébéniste", "en": "Handyman"},
		"83": {"fr": "Maitre d'armes", "en": "Weapons Master"},
		"80": {"fr": "Maroquinier", "en": "Leather Dealer"},
		"79": {"fr": "Tailleur", "en": "Tailor"},
	}

	recipesContainer := htmlElement.DOM.Find(".ak-container.ak-panel.ak-crafts:has(.ak-panel-title:contains('Recette')), .ak-container.ak-panel.ak-crafts:has(.ak-panel-title:contains('Recipe'))")

	recipesContainer.Each(func(i int, rc *goquery.Selection) {
		// fmt.Println(recipesContainer.Text())
		// Check if container is a proper reciper or a "used in X recipes"
		if rc.Find(".ak-panel-content .ak-container.ak-panel").Length() > 0 {
			// Recipes := make(map[int]structs.Recipe)

			// statElement.Each(func(i int, s *goquery.Selection) {
			recipesContainer.Each(func(i int, rc *goquery.Selection) {
				// statsDiv := s.Find("div.ak-title")
				recipesDivs := rc.Find("div.ak-panel-content")

				// Each recipe
				recipesDivs.Each(func(i int, rd *goquery.Selection) {
					Recipe := structs.Recipe{}
					var recipeId int
					var jobId int
					jobStringLevel := htmlElement.DOM.Find(".ak-panel-intro").Text()
					jobString := StandardizeSpaces(strings.Split(jobStringLevel, "-")[0])
					jobLevelString := StandardizeSpaces(strings.Split(jobStringLevel, "-")[1])
					jobLevel, err := strconv.Atoi(strings.Split(jobLevelString, " ")[1])
					if err != nil {
						fmt.Println("Error converting joblevel")
						fmt.Println(err)
					}

					JobName := structs.Display{}
					if lang == "Fr" {
						for key := range jobs {
							if jobs[key]["fr"] == jobString {
								jobId, err = strconv.Atoi(StandardizeSpaces(key))
								if err != nil {
									fmt.Println("Error getting jobId")
									jobId = 0
								}
								JobName.Fr = jobs[key]["fr"]
								JobName.En = jobs[key]["en"]
							}
						}
					}

					Ingredients := make(map[int]structs.Ingredient)
					ingredientsRows := rd.Find(".ak-list-element")
					// Each individual ingredients
					ingredientsRows.Each(func(i int, ir *goquery.Selection) {
						Ingredient := structs.Ingredient{}
						// Quantity
						quantityString := StandardizeSpaces(ir.Find(".ak-front").Text())
						quantityValue, err := strconv.Atoi(StandardizeSpaces(strings.Split(quantityString, "x")[0]))
						if err != nil {
							fmt.Printf("Error getting quant value of %d\n", recipeId)
							fmt.Println(err)
						}

						var ingredientArgName string
						var ingId int
						ingNameDiv := ir.Find(".ak-title")
						ingTypeName := StandardizeSpaces(ingNameDiv.Find(".ak-text").Text())
						ingredientHref, exists := ingNameDiv.Find("a").Attr("href")
						if exists {
							ingredientArgName, err = GetItemURLArg(ingredientHref)
							if err != nil {
								fmt.Printf("error getting ingredientArgName from %s\n", ingredientHref)
							}
							ingId, err = utils.GetItemIDFromString(ingredientArgName)
							if err != nil {
								fmt.Printf("error getting id from %s\n", ingredientHref)
							}
						} else {
							fmt.Println("Ingredient href does not exists")
						}

						// ingName
						ingName := StandardizeSpaces(ingNameDiv.Find(".ak-linker").Text())

						// Compare ingTypeName with title.fr inside itemTypes.json
						if lang == "Fr" {
							ItemTypesProperties := utils.GetItemTypesPropertiesJSON()
							for _, t := range ItemTypesProperties {
								if t.Title.Fr == ingTypeName {
									Ingredient.TypeID = t.Definition.ID
								}
							}
							Ingredient.Quantity = quantityValue
							Ingredient.IngName.Fr = ingName
							Ingredient.ID = ingId
							Ingredients[ingId] = Ingredient
						}
						// fmt.Println("jobId: ", jobId)
						// fmt.Println("i: ", i)
						// fmt.Println("Item.Params.TypeId: ", Item.Params.TypeId)
						// fmt.Println("Item.ID : ", Item.ID)
						recipeId, err = strconv.Atoi(fmt.Sprintf("%d%d%d%d", jobId, i, Item.Params.TypeId, Item.ID))
						if err != nil {
							fmt.Println("error while concatenating recipeId")
							fmt.Println(err)
							os.Exit(0)
						}
						fmt.Println("recipeId:", recipeId)

						if lang == "En" && len(Recipes) > 0 {
							for _, recipe := range Recipes {
								for ingKey, ing := range Recipes[recipeId].Ingredients {
									if ing.ID == ingId {
										ing.IngName.En = ingName

										recipe.Ingredients[ingKey] = ing

										// fmt.Printf("Updated English name for ingredient with ID %d to %s\n", ingId, ing.IngName.En)

										break
									}
								}
							}
						}
					})

					Recipe.JobID = jobId
					Recipe.JobLevel = jobLevel
					Recipe.JobName = JobName
					Recipe.RecipeId = recipeId
					Recipe.Ingredients = Ingredients

					fmt.Println("Recipe", Recipe)
					fmt.Println("recipeId", recipeId)

					Recipes[recipeId] = Recipe
					fmt.Println("Recipes --- ", Recipes)
				})
			})

			// if lang == "Fr" {
			// 	Item.Recipes = Recipes
			// }
		}
	})

	// TODO : what happens if return empty?
	if recipesContainer.Length() == 0 {
		fmt.Println("No recipes found")
		return
	}
}
