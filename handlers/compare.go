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

	// TODO Handle use case when ratings don't exist for book or movie
	if strings.Contains(bookStatus, "Unable to find this title") &&
		strings.Contains(movieStatus, "Unable to find this title") {
		errorResp, _ := json.Marshal(models.ErrorResponse)
		errorData := string(errorResp)
		return errorData
	} else if strings.Contains(movieStatus, "Unable to find this title") { // Only book data can be retrieved

		recommend.Recommendation = models.NoMovieFound
		titleString, _ := jsonparser.GetString([]byte(bookData), "Title")
		recommend.Title = titleString
		authorString, _ := jsonparser.GetString([]byte(bookData), "Author")
		recommend.Author = authorString
		bookRatingString, _ := jsonparser.GetString([]byte(bookData), "Rating")
		recommend.BookRating = bookRatingString
		bookPlotString, _ := jsonparser.GetString([]byte(bookData), "Plot")
		recommend.BookPlot = bookPlotString
		publishingYear, _ := jsonparser.GetString([]byte(bookData), "Year")
		recommend.BookPublished = publishingYear

	} else if strings.Contains(bookStatus, "Unable to find this title") { // Only movie data can be retrieved

		recommend.Recommendation = models.NoBookFound
		titleString, _ := jsonparser.GetString([]byte(movieData), "Title")
		recommend.Title = titleString
		movieRatingString, _ := jsonparser.GetString([]byte(movieData), "Rating")
		recommend.MovieRating = movieRatingString
		actors, _ := jsonparser.GetString([]byte(movieData), "Actors")
		recommend.Actors = actors
		director, _ := jsonparser.GetString([]byte(movieData), "Director")
		recommend.Director = director
		moviePlotString, _ := jsonparser.GetString([]byte(movieData), "Plot")
		recommend.MoviePlot = moviePlotString
		movieYear, _ := jsonparser.GetString([]byte(movieData), "ReleaseYear")
		recommend.MovieReleased = movieYear
		movieAwards, _ := jsonparser.GetString([]byte(movieData), "Awards")
		recommend.MovieAwards = movieAwards

	} else { // Both book and movie data is retrieved
		return parser.ParseReccomendation(bookData, movieData)
	}

	recommendationData, _ := json.Marshal(recommend)
	reccommendResponse := string(recommendationData)
	return reccommendResponse
}
