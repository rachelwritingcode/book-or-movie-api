package parser

import (
	"encoding/json"
	"math"
	"strconv"

	models "book-or-movie-api/models"

	"github.com/buger/jsonparser"
)

//TODO fix bug
func ParseReccomendation(bookData string, movieData string) string {

	var recommend = new(models.Recommendation)
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

	// Movie Data
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

	// Compare Movie and Book Ratings
	bookRatingFloat, _ := strconv.ParseFloat(bookRatingString, 64)
	movieRatingFloat, _ := strconv.ParseFloat(movieRatingString, 64)

	if math.Max(bookRatingFloat, movieRatingFloat) == bookRatingFloat {
		recommend.Recommendation = models.BookFirst
	} else if math.Max(bookRatingFloat, movieRatingFloat) == movieRatingFloat {
		recommend.Recommendation = models.MovieFirst
	} else {
		recommend.Recommendation = models.BothFirst
	}
	recommendationData, _ := json.Marshal(recommend)

	reccommendResponse := string(recommendationData)
	return reccommendResponse

}
