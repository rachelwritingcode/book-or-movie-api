package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	models "book-or-movie-api/models"
	parser "book-or-movie-api/parser"

	"github.com/buger/jsonparser"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
)

const openLibraryAPI = "http://openlibrary.org/search.json?title="
const openLibraryURL = "http://openlibrary.org"

func GetMovie(c echo.Context) error {
	title := c.QueryParam("title")
	movieData := GetMovieData(title)
	return c.String(http.StatusOK, movieData)
}

func GetMovieData(title string) string {

	// Load Environment Variables
	err := godotenv.Load("local.env")
	if err != nil {
		log.Fatalln(err)
	}

	// Environment Variables
	omdb_api := os.Getenv("OMDB_API")
	omdb_key := os.Getenv("OMDB_KEY")

	request_url := omdb_api + "=" + omdb_key + "&t=" + title
	req, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	if strings.Contains(string(b), "Movie not found") {
		errorResp, _ := json.Marshal(models.ErrorResponse)
		errorData := string(errorResp)
		return errorData
	}

	director, _ := jsonparser.GetString(b, "Director")
	plot, _ := jsonparser.GetString(b, "Plot")
	actors, _ := jsonparser.GetString(b, "Actors")
	rating, _ := jsonparser.GetString(b, "imdbRating")
	year, _ := jsonparser.GetString(b, "Year")
	awards, _ := jsonparser.GetString(b, "Awards")
	awards = strings.ReplaceAll(awards, "\u0026", "")

	movieResponseData := new(models.MovieResponse)
	movieResponseData = &models.MovieResponse{
		title,
		rating,
		actors,
		director,
		plot,
		year,
		awards}

	movieData, _ := json.Marshal(movieResponseData)
	return string(movieData)
}

func GetBookData(title string) string {

	request_url := openLibraryAPI + title
	req, err := http.NewRequest("GET", request_url, nil)
	if err != nil {
		log.Fatalln(err)
	}
	req.Header.Set("Accept", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalln(err)
	}

	defer resp.Body.Close()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	numFound, _ := jsonparser.GetInt(b, "numFound")

	if strings.Contains(string(b), "400 Bad request") || numFound == 0 {
		errorResp, _ := json.Marshal(models.ErrorResponse)
		errorData := string(errorResp)
		return errorData
	}

	open_library_id, _ := jsonparser.GetString(b, "docs", "[0]", "key")
	var rating string
	var plot string

	rating, plot = parser.WebScrape(open_library_id)

	data_title, _ := jsonparser.GetString(b, "docs", "[0]", "title")
	yearInt, _ := jsonparser.GetInt(b, "docs", "[0]", "first_publish_year")
	yearString := strconv.FormatInt(int64(yearInt), 10)

	author, _ := jsonparser.GetString(b, "docs", "[0]", "author_name")

	bookResponse := &models.BookResponse{
		data_title,
		author,
		rating,
		plot,
		yearString,
	}
	bookResponseData, _ := json.Marshal(bookResponse)
	bookData := string(bookResponseData)
	return bookData
}

// Add %20 to the request url
func cleanTitleQuery(title string) string {

	query_parameter := ""
	title_slice := strings.Split(title, " ")

	for _, s := range title_slice {
		query_parameter += s + "%20"
	}

	return query_parameter
}

func GetBook(c echo.Context) error {

	originalTitle := c.QueryParam("title")
	title := cleanTitleQuery(originalTitle)
	bookData := GetBookData(title)
	return c.String(http.StatusOK, bookData)
}

func GetReccomend(c echo.Context) error {

	originalTitle := c.QueryParam("title")
	title := cleanTitleQuery(originalTitle)

	bookData := GetBookData(title)
	movieData := GetMovieData(title)

	recommendation := compareBookAndMovieData(movieData, bookData, title)
	status, _ := jsonparser.GetString([]byte(recommendation), "Message")

	if strings.Contains(status, models.NoBookOrMovieFound) {
		return c.String(http.StatusNotFound, recommendation)
	}
	return c.String(http.StatusOK, recommendation)
}

//TODO cleaning up this function
func compareBookAndMovieData(movieData string, bookData string, title string) string {

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
		parser.ParseReccomendation(bookData, movieData)
	}

	recommendationData, _ := json.Marshal(recommend)
	reccommendResponse := string(recommendationData)
	return reccommendResponse
}
