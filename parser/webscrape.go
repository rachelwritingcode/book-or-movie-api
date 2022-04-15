package parser

import (
	"math"
	"strconv"
	"strings"
	"time"

	"github.com/gocolly/colly"
)

const openLibraryAPI = "http://openlibrary.org/search.json?title="
const openLibraryURL = "http://openlibrary.org"

func WebScrape(open_library_id string) (string, string) {

	var Rating = ""
	var Plot string = ""
	var ratings = []string{}

	web_scraper := colly.NewCollector(
		colly.AllowedDomains(openLibraryURL, "openlibrary.org"),
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

	web_scraper.Visit(openLibraryURL + open_library_id)

	Rating = strings.ReplaceAll(ratings[len(ratings)-1], "★", "")
	rating_float, _ := strconv.ParseFloat(Rating, 64)
	rating_float = (math.Ceil(rating_float*100) / 100) / 5 * 10
	Rating = strconv.FormatFloat(rating_float, 'E', -1, 64)
	Rating = strings.ReplaceAll(Rating, "E+00", "")

	Plot = strings.ReplaceAll(Plot, "\n", "")
	Plot = strings.Trim(Plot, " ")
	Plot = strings.ReplaceAll(Plot, "\\", "")

	return Rating, Plot
}
