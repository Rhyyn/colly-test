package scrapers

import (
	"encoding/json"
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

var StatPropertiesMap = make(map[string]structs.StatProperties)

// This create the map of properties for stats
func handleStatsProperties(data []byte) map[string]structs.StatProperties {
	var translations map[string]structs.StatProperties

	err := json.Unmarshal([]byte(data), &translations)
	if err != nil {
		panic(err)
	}

	return translations
}

func ScrapItems(urlPrefix string, maxPage int, selectedId int) {
	// construct the suffix
	urlSuffix := "&" + "type_1%5B%5D=" + strconv.Itoa(selectedId)

	fmt.Printf("ScrapItems called for id %d with maxPage %d\n", selectedId, maxPage)
	c := GetNewCollector()

	c.OnHTML(".ak-table.ak-responsivetable tbody tr", func(h *colly.HTMLElement) {
		// finds all the tds, only follow the href inside the 2nd
		tds := h.DOM.Find("td")
		secondTd := tds.Eq(1)
		link := secondTd.Find("a[href]")
		href, _ := link.Attr("href")
		frenchUrl := h.Request.AbsoluteURL(href)
		scrapeItemPage(frenchUrl)
	})

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("ScrapItems visiting:\n", r.URL)
	})

	// visit each page until = maxPage
	for i := 1; i < maxPage; i++ {
		c.Visit(urlPrefix + strconv.Itoa(i) + urlSuffix)
		fmt.Printf("setting page to %d\n", i)
	}
}

func scrapeItemPage(url string) {
	c := GetNewCollector()
	defer c.Wait()
	all_positives_stats := handleStatsProperties(utils.ReadFile(utils.OpenFile("all_positives_stats.json")))
	all_negatives_stats := handleStatsProperties(utils.ReadFile(utils.OpenFile("all_negatives_stats.json")))

	// This gets english href
	// c.OnHTML(".ak-idbar-box.ak-box-lang .ak-flag-en", func(h *colly.HTMLElement) {
	// 	enHref := h.Attr("href")
	// 	fmt.Println(enHref)
	// })

	var FrItem structs.Item
	// channel for FrItem pop
	done := make(chan bool)
	// item page -- /fr/
	c.OnHTML(".ak-container.ak-panel-stack.ak-glue ", func(e *colly.HTMLElement) {
		// Title
		titleElement := e.DOM.Find("h1").First()
		titleFr := StandardizeSpaces(titleElement.Text())

		// typeId
		typeElement := e.DOM.Find(".ak-encyclo-detail-type.col-xs-6 span img").First()
		src, _ := typeElement.Attr("src")
		typeParts := strings.Split(src, "/")
		category := typeParts[len(typeParts)-1]
		typeId, err := strconv.Atoi(strings.Split(category, ".")[0])
		if err != nil {
			fmt.Printf("Error converting typeId")
			os.Exit(0)
		}

		// level
		levelElement := e.DOM.Find(".ak-encyclo-detail-level.col-xs-6.text-right")
		levelElementFr := StandardizeSpaces(levelElement.Text())
		level, err := strconv.Atoi(strings.TrimSpace(strings.Split(levelElementFr, ": ")[1]))
		if err != nil {
			fmt.Printf("Error converting level")
			os.Exit(0)
		}

		// Rarity
		rarityElement := e.DOM.Find(".ak-object-rarity span span")
		rarityNumber, _ := strconv.Atoi(strings.Split(rarityElement.AttrOr("class", ""), "ak-rarity-")[1])

		// Stats
		statElement := e.DOM.Find(".ak-container.ak-content-list.ak-displaymode-col")
		Stats := make(map[int]structs.Stat)
		statElement.Each(func(i int, s *goquery.Selection) {
			statsDiv := s.Find("div.ak-title")

			statsDiv.Each(func(i int, stat *goquery.Selection) {
				var id int
				var displayFr string
				var displayEn string
				var format string
				var value int
				var numElements int

				entireStatString := strings.SplitN(StandardizeSpaces(stat.Text()), " ", 2)
				// "160"
				prefixString := StandardizeSpaces(entireStatString[0])
				// "Esquive"
				suffixString := StandardizeSpaces(entireStatString[1])

				format, value = utils.StatPrefixToStringAndSetFormat(prefixString)

				// Check if has 2 values (Mastery of 3 random elements)
				// If it has set numElements and format it to (Mastery in X elements)
				if hasNumber, number := utils.HasNumberInString(suffixString); hasNumber {
					numElements = number
					suffixString = utils.FormatElementsString(utils.ReplaceAnyNumberInString(suffixString, "X"))
				}

				displayFr = suffixString
				displayEn = suffixString
				if strings.HasPrefix(suffixString, "-") {
					fmt.Println("NEGATIVE SUFFIX")
					for key, property := range all_negatives_stats {
						if strings.Contains(suffixString, property.Fr) {
							id, _ = strconv.Atoi(key)
						}
					}
				} else {
					for key, property := range all_positives_stats {
						if strings.Contains(suffixString, property.Fr) {
							id, _ = strconv.Atoi(key)
						}
					}
				}
				newStat := structs.Stat{
					Display: structs.Display{
						Fr: displayFr,
						En: displayEn,
					},
					ID:          id,
					Format:      format,
					Value:       value,
					NumElements: numElements,
				}

				Stats[newStat.ID] = newStat
				// fmt.Println(newStat.ID)
				// fmt.Println(Stats)
				// fmt.Printf("prefixValue : %d | suffixString %s\n", prefixValue, suffixString)
			})

			// need to fetch IDs from separate file (manual)
		})

		// Droprates

		// Recipe
		fmt.Printf("titleFr %s\n", titleFr)
		fmt.Printf("typeId %d\n", typeId)
		fmt.Printf("level : %d\n", level)
		fmt.Printf("rarityNumber %d\n", rarityNumber)
		fmt.Printf("Stats : %v\n", Stats)

		FrItem.Title.Fr = titleFr
		FrItem.Params.TypeId = typeId
		FrItem.Params.Level = level
		FrItem.Params.Rarity = rarityNumber
		FrItem.Stats = Stats

		fmt.Println("------VV FrItem INSIDE c.OnHTML VV-----")
		fmt.Println(FrItem)
		fmt.Println("end of c.OnHTML")
		done <- true
	})

	go func() {
		<-done
		// useless pretty print
		PrettyFrItem, err := json.MarshalIndent(FrItem, "", "    ")
		if err != nil {
			fmt.Println("Error marshaling item:", err)
			return
		}
		fmt.Println("FrItem after scraping:\n", string(PrettyFrItem))
	}()

	c.OnError(func(r *colly.Response, err error) {
		fmt.Println("Request URL:", r.Request.URL, "failed with response:", r, "\nError:", err)
	})

	c.OnRequest(func(r *colly.Request) {
		fmt.Println("scrapeItemPage visiting :\n", r.URL)
	})

	c.Visit(url)
}
