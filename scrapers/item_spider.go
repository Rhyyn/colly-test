package scrapers

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/Rhyyn/wakfukiscraper/structs"
	"github.com/Rhyyn/wakfukiscraper/utils"
	"github.com/gocolly/colly"
)

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

func GetStatId(statString string, ParamsStatsProperties structs.ParamsStatsProperties, format string, lang string) (int, error) {
	var id int
	var err error
	if format == "negative" {
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
			format, value = utils.StatPrefixToStringAndSetFormat(prefixString)

			// Check if has 2 values (Mastery of 3 random elements)
			// If it has set numElements and format it to (Mastery in X elements)
			if hasNumber, number := utils.HasNumberInString(suffixString); hasNumber {
				numElements = number
				suffixString = utils.ReplaceAnyNumberInString(suffixString, "X")
			}
			// Format the strings because Ankama's english is dogshit
			suffixString = utils.FormatStatString(suffixString, format)

			// Gets the ID
			id, idError = GetStatId(suffixString, ParamsStatsProperties, format, lang)
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

func GetTypeID(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
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

func GetTitle(htmlElement *colly.HTMLElement, Item *structs.Item, lang string) {
	titleElement := htmlElement.DOM.Find("h1").First()
	title := StandardizeSpaces(titleElement.Text())
	if lang == "Fr" {
		Item.Title.Fr = title
	} else {
		Item.Title.En = title
	}
}

func ScrapItemDetails(url string, Item *structs.Item, ParamsStatsProperties structs.ParamsStatsProperties) {
	c := GetNewCollector()
	c.OnHTML(".ak-container.ak-panel-stack.ak-glue", func(h *colly.HTMLElement) {
		Lang := utils.GetLangFromURL(h.Request.URL.String())
		GetTitle(h, Item, Lang)
		GetTypeID(h, Item, Lang)
		GetRarity(h, Item, Lang)
		GetLevel(h, Item, Lang)
		GetStats(h, Item, Lang, ParamsStatsProperties)
		GetDroprates(h, Item, Lang)
		// GetRecipes ( need to check if multiple recipes exists)
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

func ScrapItems(indexURL map[string]string, itemURL map[string]string, maxPage int, selectedId int) {
	urlSuffix := "&" + "type_1%5B%5D=" + strconv.Itoa(selectedId)
	AllPositivesStats := utils.HandleStatsProperties(utils.ReadFile(utils.OpenFile("all_positives_stats.json")))
	AllNegativesStats := utils.HandleStatsProperties(utils.ReadFile(utils.OpenFile("all_negatives_stats.json")))

	fmt.Printf("ScrapItems called for id %d with maxPage %d\n", selectedId, maxPage)
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

		frenchURL := itemURL["fr"] + "/" + itemArgName
		englishURL := itemURL["en"] + "/" + itemArgName

		var Item structs.Item
		ParamsStatsProperties := structs.ParamsStatsProperties{AllPositivesStats: AllPositivesStats, AllNegativesStats: AllNegativesStats}

		Item.ID, _ = utils.GetItemIDFromString(itemArgName)
		// Scrap both FR/EN version of the item
		ScrapItemDetails(frenchURL, &Item, ParamsStatsProperties)
		ScrapItemDetails(englishURL, &Item, ParamsStatsProperties)

		// TODO: Add to separate map of item and write to file

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
		fmt.Println("ScrapItems visiting:\n", r.URL)
	})

	// visit each page until = maxPage
	for i := 1; i < maxPage; i++ {
		c.Visit(indexURL["fr"] + strconv.Itoa(i) + urlSuffix)
		fmt.Printf("setting page to %d\n", i)
	}
}
