package main

import (
	"fmt"
	"time"

	"github.com/gocolly/colly"
)

func main() {

	// var ratings = []string{}
	web_scraper := colly.NewCollector(
		colly.AllowedDomains("https://openlibrary.org", "openlibrary.org"),
	)

	// web_scraper.OnHTML(".readers-stats ", func(e *colly.HTMLElement) {
	// 	e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
	// 		ratings = append(ratings, e.DOM.Find("span").Text())
	// 	})
	// })

	var plot string = ""
	web_scraper.OnHTML(".book-description-content", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
			plot = e.Text
		})
	})

	web_scraper.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})

	web_scraper.OnRequest(func(request *colly.Request) {
		fmt.Println("Visiting...", request.URL.String())
	})

	web_scraper.Visit("https://openlibrary.org/books/OL27710722M/")

	fmt.Println(plot)
	// rating := strings.ReplaceAll(ratings[len(ratings)-1], "★", "")
	// rating_float, _ := strconv.ParseFloat(rating, 64)
	// rating_float = math.Ceil(rating_float*100) / 100
	// rating = strconv.FormatFloat(rating_float, 'E', -1, 64)
	// rating = strings.ReplaceAll(rating, "E+00", "")
	// fmt.Println(rating)

}
