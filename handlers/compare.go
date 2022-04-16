package handlers

import (
	"encoding/json"
	"strings"

	models "book-or-movie-api/models"
	parser "book-or-movie-api/parser"

	"github.com/buger/jsonparser"
)

func CompareBookAndMovieData(movieData string, bookData string, title string) string {

	recommend := &models.Recommendation{}
	bookStatus, _ := jsonparser.GetString([]byte(bookData), "Message")
	movieStatus, _ := jsonparser.GetString([]byte(movieData), "Message")

	if strings.Contains(bookStatus, models.NoTitle) &&
		strings.Contains(movieStatus, models.NoTitle) {

		errorResp, _ := json.Marshal(models.ErrorResponse)
		errorData := string(errorResp)
		return errorData
	} else if strings.Contains(movieStatus, models.NoTitle) {

		recommend.Recommendation = models.NoMovieFound
		recommend.Title, _ = jsonparser.GetString([]byte(bookData), "Title")
		recommend.Author, _ = jsonparser.GetString([]byte(bookData), "Author")
		recommend.BookRating, _ = jsonparser.GetString([]byte(bookData), "Rating")
		recommend.BookPlot, _ = jsonparser.GetString([]byte(bookData), "Plot")
		recommend.BookPublished, _ = jsonparser.GetString([]byte(bookData), "Year")
	} else if strings.Contains(bookStatus, models.NoTitle) {

		recommend.Recommendation = models.NoBookFound
		recommend.Title, _ = jsonparser.GetString([]byte(movieData), "Title")
		recommend.MovieRating, _ = jsonparser.GetString([]byte(movieData), "Rating")
		recommend.Actors, _ = jsonparser.GetString([]byte(movieData), "Actors")
		recommend.Director, _ = jsonparser.GetString([]byte(movieData), "Director")
		recommend.MoviePlot, _ = jsonparser.GetString([]byte(movieData), "Plot")
		recommend.MovieReleased, _ = jsonparser.GetString([]byte(movieData), "ReleaseYear")
		recommend.MovieAwards, _ = jsonparser.GetString([]byte(movieData), "Awards")
	} else {
		return parser.ParseReccomendation(bookData, movieData)
	}
	recommendationData, _ := json.Marshal(recommend)
	reccommendResponse := string(recommendationData)
	return reccommendResponse
}
