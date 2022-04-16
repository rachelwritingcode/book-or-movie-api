package parser

import (
	models "book-or-movie-api/models"
	"encoding/json"
	"strconv"

	"github.com/buger/jsonparser"
)

func ParseBook(response []byte, title string, openLibraryId string) string {

	bookResponseData := new(models.BookResponse)

	var rating string
	var plot string
	rating, plot = WebScrape(openLibraryId)
	data_title, _ := jsonparser.GetString(response, "docs", "[0]", "title")
	yearInt, _ := jsonparser.GetInt(response, "docs", "[0]", "first_publish_year")
	yearString := strconv.FormatInt(int64(yearInt), 10)
	author, _ := jsonparser.GetString(response, "docs", "[0]", "author_name")

	bookResponseData.Title = data_title
	bookResponseData.Author = author
	bookResponseData.Year = yearString
	bookResponseData.Plot = plot
	bookResponseData.Rating = rating

	bookData, _ := json.Marshal(bookResponseData)
	book := string(bookData)
	return book
}
