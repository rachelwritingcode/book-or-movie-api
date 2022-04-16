package parser

import (
	models "book-or-movie-api/models"
	"encoding/json"
	"strings"

	"github.com/buger/jsonparser"
)

func ParseMovie(response []byte, title string) string {

	movieResponseData := new(models.MovieResponse)
	director, _ := jsonparser.GetString(response, "Director")
	plot, _ := jsonparser.GetString(response, "Plot")
	actors, _ := jsonparser.GetString(response, "Actors")
	rating, _ := jsonparser.GetString(response, "imdresponseRating")
	year, _ := jsonparser.GetString(response, "Year")
	awards, _ := jsonparser.GetString(response, "Awards")
	awards = strings.ReplaceAll(awards, "\u0026", "")

	movieResponseData.Title = title
	movieResponseData.Director = director
	movieResponseData.Plot = plot
	movieResponseData.Actors = actors
	movieResponseData.Rating = rating
	movieResponseData.ReleaseYear = year
	movieResponseData.Awards = awards

	movieData, _ := json.Marshal(movieResponseData)
	movie := string(movieData)
	return movie
}
