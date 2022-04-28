package parser

import (
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const URL = "http://openlibrary.org"
const domain = "openlibrary.org"

func WebScrape(open_library_id string) (string, string) {

	var Rating = ""
	var Plot string = ""
	var ratings = []string{}

	web_scraper := colly.NewCollector(
		colly.AllowedDomains(URL, domain),
	)
	web_scraper.OnHTML(".readers-stats ", func(e *colly.HTMLElement) {
		e.ForEach("li", func(_ int, elem *colly.HTMLElement) {
			ratings = append(ratings, e.DOM.Find("span").Text())
		})
	})
	web_scraper.OnHTML(".book-description-content", func(e *colly.HTMLElement) {
		e.ForEach("p", func(_ int, elem *colly.HTMLElement) {
			Plot = e.Text
		})
	})
	web_scraper.Limit(&colly.LimitRule{
		DomainGlob:  "*",
		RandomDelay: 1 * time.Second,
	})
	web_scraper.OnRequest(func(request *colly.Request) {
	})
	web_scraper.Visit(URL + open_library_id)

	Rating = strings.ReplaceAll(ratings[len(ratings)-1], "★", "")
	ratingStr := strings.Split(Rating, ".")
	first_digits := ratingStr[0]
	end_digits := strings.Split(ratingStr[1], "")
	ending_decimals := end_digits[0] + end_digits[1]
	Rating = first_digits + "." + ending_decimals

	Plot = strings.ReplaceAll(Plot, "\n", "")
	Plot = strings.Trim(Plot, " ")
	Plot = strings.ReplaceAll(Plot, "\\", "")

	return Rating, Plot
}
