package parser

import (
	"encoding/json"
	"math"
	"strconv"

	models "book-or-movie-api/models"

	"github.com/buger/jsonparser"
)

func ParseReccomendation(bookData string, movieData string) string {

	var recommend = new(models.Recommendation)
	recommend.Title, _ = jsonparser.GetString([]byte(bookData), "Title")
	recommend.Author, _ = jsonparser.GetString([]byte(bookData), "Author")
	recommend.BookRating, _ = jsonparser.GetString([]byte(bookData), "Rating")
	recommend.BookPlot, _ = jsonparser.GetString([]byte(bookData), "Plot")
	recommend.BookPublished, _ = jsonparser.GetString([]byte(bookData), "Year")

	recommend.MovieRating, _ = jsonparser.GetString([]byte(movieData), "Rating")
	recommend.Actors, _ = jsonparser.GetString([]byte(movieData), "Actors")
	recommend.Director, _ = jsonparser.GetString([]byte(movieData), "Director")
	recommend.MoviePlot, _ = jsonparser.GetString([]byte(movieData), "Plot")
	recommend.MovieReleased, _ = jsonparser.GetString([]byte(movieData), "ReleaseYear")
	recommend.MovieAwards, _ = jsonparser.GetString([]byte(movieData), "Awards")

	bookRatingFloat, _ := strconv.ParseFloat(recommend.BookRating, 64)
	movieRatingFloat, _ := strconv.ParseFloat(recommend.MovieRating, 64)

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
