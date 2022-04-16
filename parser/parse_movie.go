package parser

import (
	models "book-or-movie-api/models"
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
)

func ParseMovie(response []byte, title string) string {

	movieResponseData := new(models.MovieResponse)
	movieResponseData.Title = title
	movieResponseData.Director, _ = jsonparser.GetString(response, "Director")
	movieResponseData.Plot, _ = jsonparser.GetString(response, "Plot")
	movieResponseData.Actors, _ = jsonparser.GetString(response, "Actors")
	movieResponseData.Rating, _ = jsonparser.GetString(response, "imdresponseRating")
	movieResponseData.ReleaseYear, _ = jsonparser.GetString(response, "Year")
	movieResponseData.Awards, _ = jsonparser.GetString(response, "Awards")
	movieResponseData.Awards = strings.ReplaceAll(movieResponseData.Awards, "\u0026", "")

	movieData, _ := json.Marshal(movieResponseData)
	movie := string(movieData)
	return movie
}
