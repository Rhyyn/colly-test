package main

import (
	cli "github.com/Rhyyn/wakfukiscraper/CLI"
)

func main() {
	cli.Execute()
	// c := colly.NewCollector()

	// Find and visit all links
	// c.OnHTML("a[href]", func(e *colly.HTMLElement) {
	// 	e.Request.Visit(e.Attr("href"))
	// })

	// c.OnRequest(func(r *colly.Request) {
	// 	fmt.Println("Visiting", r.URL)
	// })
	// c.Visit("https://wakfuki.com")

	// c.Visit("https://www.wakfu.com/fr/mmorpg/encyclopedie/objets/109-equipo/118-ropa/120-collar/s.ankama.com?text=&type_1[0]=218&level_min=0&level_max=230&page=1")
}
