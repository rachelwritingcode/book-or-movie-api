package handlers

import (
	"encoding/json"
	"strings"

	models "book-or-movie-api/models"
	parser "book-or-movie-api/parser"

	"github.com/buger/jsonparser"
)

func GetMovieData(title string) string {

	responseBytes := Request(title, true)

	if strings.Contains(string(responseBytes), "Movie not found") {
		errorResp, _ := json.Marshal(models.ErrorResponse)
		errorData := string(errorResp)
		return errorData
	}

	movieData := parser.ParseMovie(responseBytes, title)
	return movieData
}

func GetBookData(title string) string {

	responseBytes := Request(title, false)
	numFound, _ := jsonparser.GetInt(responseBytes, "numFound")
	openLibraryId, _ := jsonparser.GetString(responseBytes, "docs", "[0]", "key")

	if strings.Contains(string(responseBytes), "400 Bad request") || numFound == 0 || openLibraryId == "" {
		errorResp, _ := json.Marshal(models.ErrorResponse)
		errorData := string(errorResp)
		return errorData
	}

	bookData := parser.ParseBook(responseBytes, title, openLibraryId)
	return bookData
}
